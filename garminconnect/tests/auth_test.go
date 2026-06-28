package garminconnect_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

// jsonResp builds a canned JSON HTTP response for the SSO/auth stubs below.
func jsonResp(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": {"application/json"}},
	}
}

// TestLogin_FetchesProfile verifies that Login() populates DisplayName via the
// social profile endpoint when a valid token is already present.
func TestLogin_FetchesProfile(t *testing.T) {
	r, err := recorder.New("testdata/cassettes/login_profile",
		recorder.WithMode(recorder.ModeReplayOnly),
		recorder.WithMatcher(func(req *http.Request, i cassette.Request) bool {
			cu, _ := url.Parse(i.URL)
			return req.Method == i.Method && normaliseURL(req.URL) == normaliseURL(cu)
		}),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, r.Stop()) }()

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: r}),
		gc.WithToken("test"),
	)

	require.NoError(t, c.Login(t.Context(), "", ""))
	assert.Equal(t, "testuser", c.DisplayName())
}

// TestLogout clears the in-memory token and removes the cached token file.
func TestLogout(t *testing.T) {
	tokenFile := filepath.Join(t.TempDir(), "token.json")
	require.NoError(t, os.WriteFile(tokenFile, []byte(`{"access_token":"x"}`), 0600))

	c := gc.NewClient(tokenFile, gc.WithToken("x"), gc.WithDisplayName("testuser"))
	require.NoError(t, c.Logout())

	assert.Empty(t, c.Token())
	assert.Empty(t, c.DisplayName())
	assert.NoFileExists(t, tokenFile)
}

// TestLogout_NoTokenFile is a no-op when no token file is configured.
func TestLogout_NoTokenFile(t *testing.T) {
	c := gc.NewClient("", gc.WithToken("x"))
	require.NoError(t, c.Logout())
	assert.Empty(t, c.Token())
}

// TestLogin_SSO drives the full mobile-SSO flow (login -> ticket -> token ->
// profile) against stubbed endpoints, and verifies the cached token file is
// written owner-only even when it pre-exists world-readable (GHSA-wjhr-76vg-2hvc).
func TestLogin_SSO_WritesOwnerOnlyToken(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX file modes only")
	}
	tokenFile := filepath.Join(t.TempDir(), "token.json")
	require.NoError(t, os.WriteFile(tokenFile, []byte("{}"), 0644)) // pre-existing loose perms

	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Host + r.URL.Path {
		case "sso.garmin.com/mobile/api/login":
			return jsonResp(200, `{"serviceTicketId":"ST-1","serviceURL":"https://x"}`), nil
		case "diauth.garmin.com/di-oauth2-service/oauth/token":
			return jsonResp(200, `{"access_token":"a","refresh_token":"r","expires_in":3600}`), nil
		case "connectapi.garmin.com/userprofile-service/socialProfile":
			return jsonResp(200, `{"displayName":"testuser"}`), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", r.URL)
	})

	c := gc.NewClient(tokenFile, gc.WithHTTPClient(&http.Client{Transport: rt}))
	require.NoError(t, c.Login(t.Context(), "e@x.com", "pw"))
	assert.Equal(t, "testuser", c.DisplayName())

	info, err := os.Stat(tokenFile)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

// TestLogin_MFARequired returns ErrMFARequired when the SSO flow needs MFA and
// no prompt is configured.
func TestLogin_MFARequired(t *testing.T) {
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Host+r.URL.Path == "sso.garmin.com/mobile/api/login" {
			return jsonResp(200, `{"responseStatus":{"type":"MFA_REQUIRED"},"customerMfaInfo":{"mfaLastMethodUsed":"email"}}`), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", r.URL)
	})

	c := gc.NewClient("", gc.WithHTTPClient(&http.Client{Transport: rt}))
	err := c.Login(t.Context(), "e@x.com", "pw")
	require.ErrorIs(t, err, gc.ErrMFARequired)
}

// TestLogin_MFASuccess completes the MFA flow: the prompt supplies a code, the
// verify call returns a ticket, and login finishes.
func TestLogin_MFASuccess(t *testing.T) {
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Host + r.URL.Path {
		case "sso.garmin.com/mobile/api/login":
			return jsonResp(200, `{"responseStatus":{"type":"MFA_REQUIRED"},"customerMfaInfo":{"mfaLastMethodUsed":"email"}}`), nil
		case "sso.garmin.com/mobile/api/mfa/verifyCode":
			return jsonResp(200, `{"serviceTicketId":"ST-1","serviceURL":"https://x"}`), nil
		case "diauth.garmin.com/di-oauth2-service/oauth/token":
			return jsonResp(200, `{"access_token":"a","refresh_token":"r","expires_in":3600}`), nil
		case "connectapi.garmin.com/userprofile-service/socialProfile":
			return jsonResp(200, `{"displayName":"mfauser"}`), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", r.URL)
	})

	promptCalled := false
	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: rt}),
		gc.WithMFAPrompt(func() (string, error) { promptCalled = true; return "123456", nil }),
	)
	require.NoError(t, c.Login(t.Context(), "e@x.com", "pw"))
	assert.True(t, promptCalled)
	assert.Equal(t, "mfauser", c.DisplayName())
}

// TestLogin_RefreshesOn401 verifies that a 401 triggers a token refresh and the
// request is retried with the new token.
func TestLogin_RefreshesOn401(t *testing.T) {
	profileCalls := 0
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Host + r.URL.Path {
		case "connectapi.garmin.com/userprofile-service/socialProfile":
			profileCalls++
			if profileCalls == 1 {
				return jsonResp(401, ``), nil
			}
			return jsonResp(200, `{"displayName":"refreshed"}`), nil
		case "diauth.garmin.com/di-oauth2-service/oauth/token":
			return jsonResp(200, `{"access_token":"new","refresh_token":"r2","expires_in":3600}`), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", r.URL)
	})

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: rt}),
		gc.WithToken("old"),
		gc.WithRefreshToken("r1"),
	)
	require.NoError(t, c.Login(t.Context(), "", ""))
	assert.Equal(t, "refreshed", c.DisplayName())
	assert.Equal(t, "new", c.Token())
	assert.Equal(t, 2, profileCalls)
}

// TestLogin_RefreshFailsUnauthorized returns ErrUnauthorized when the 401 retry
// cannot refresh the token.
func TestLogin_RefreshFailsUnauthorized(t *testing.T) {
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Host + r.URL.Path {
		case "connectapi.garmin.com/userprofile-service/socialProfile":
			return jsonResp(401, ``), nil
		case "diauth.garmin.com/di-oauth2-service/oauth/token":
			return jsonResp(400, `{}`), nil
		}
		return nil, fmt.Errorf("unexpected request: %s", r.URL)
	})

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: rt}),
		gc.WithToken("old"),
		gc.WithRefreshToken("r1"),
	)
	err := c.Login(t.Context(), "", "")
	require.ErrorIs(t, err, gc.ErrUnauthorized)
}
