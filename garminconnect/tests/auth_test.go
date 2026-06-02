package garminconnect_test

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

// TestLogin_FetchesProfile verifies that Login() populates DisplayName via the
// social profile endpoint when a valid token is already present.
func TestLogin_FetchesProfile(t *testing.T) {
	cassettePath := "testdata/cassettes/login_profile"
	if _, err := os.Stat(cassettePath + ".yaml"); os.IsNotExist(err) {
		t.Skip("cassette login_profile not found; record it with record_cassettes.sh")
	}

	r, err := recorder.New(cassettePath,
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

	require.NoError(t, c.Login("", ""))
	assert.Equal(t, "testuser", c.DisplayName())
}

// TestLogin_SSO exercises the full SSO login flow end-to-end:
//
//	POST sso.garmin.com → POST diauth.garmin.com → GET socialProfile
func TestLogin_SSO(t *testing.T) {
	c, stop := newAuthVCRClient(t, "login_sso")
	defer stop()

	email := os.Getenv("GARMIN_EMAIL")
	if email == "" {
		email = "test@example.com"
	}
	password := os.Getenv("GARMIN_PASSWORD")
	if password == "" {
		password = "test"
	}

	require.NoError(t, c.Login(email, password))
	assert.Equal(t, "testuser", c.DisplayName())
	assert.NotEmpty(t, c.Token())
}

// TestLogin_MFARequired verifies that Login returns ErrMFARequired when the
// SSO server demands MFA and no prompt callback is configured.
func TestLogin_MFARequired(t *testing.T) {
	transport := fixedTransport(http.StatusOK, `{
		"responseStatus": {"type": "MFA_REQUIRED"},
		"customerMfaInfo": {"mfaLastMethodUsed": "email"}
	}`)
	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: transport}),
	)
	require.ErrorIs(t, c.Login("test@example.com", "test"), gc.ErrMFARequired)
}

// TestRefreshToken verifies the 401 → refresh → retry flow: when an API call
// returns 401, the client exchanges the refresh token for a new access token
// and transparently retries the request.
func TestRefreshToken(t *testing.T) {
	calls := 0
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		header := http.Header{"Content-Type": []string{"application/json"}}

		switch req.URL.Host {
		case "diauth.garmin.com":
			body := `{"access_token":"new_token","refresh_token":"new_refresh","expires_in":3600}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     header,
			}, nil
		default:
			if calls == 1 {
				return &http.Response{
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(strings.NewReader(`{}`)),
					Header:     make(http.Header),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"displayName":"testuser"}`)),
				Header:     header,
			}, nil
		}
	})

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: transport}),
		gc.WithToken("expired_token"),
		gc.WithRefreshToken("test_refresh_token"),
	)

	require.NoError(t, c.Login("", ""))
	assert.Equal(t, "testuser", c.DisplayName())
	assert.Equal(t, "new_token", c.Token())
}
