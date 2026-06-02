package garminconnect_test

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
)

// newVCRClient returns a Client wired to the named cassette for replay.
// To record a new cassette set either:
//   - GARMIN_TOKEN + GARMIN_DISPLAY_NAME  (preferred: no extra SSO call)
//   - GARMIN_EMAIL + GARMIN_PASSWORD      (triggers one SSO login per test)
//
// Real credentials are scrubbed from saved cassettes so they can be committed:
//   - Authorization header → "Bearer test"
//   - real display name in URLs → "testuser"
func newVCRClient(t *testing.T, cassetteName string) (*gc.Client, func()) {
	t.Helper()

	cassettePath := "testdata/cassettes/" + cassetteName
	token := os.Getenv("GARMIN_TOKEN")
	displayName := os.Getenv("GARMIN_DISPLAY_NAME")
	email := os.Getenv("GARMIN_EMAIL")
	password := os.Getenv("GARMIN_PASSWORD")

	needsRecording := false
	if _, err := os.Stat(cassettePath + ".yaml"); os.IsNotExist(err) {
		needsRecording = true
	}

	if needsRecording && token == "" && (email == "" || password == "") {
		t.Skipf("cassette %q not found; set GARMIN_TOKEN+GARMIN_DISPLAY_NAME or GARMIN_EMAIL+GARMIN_PASSWORD to record", cassetteName)
	}

	// Obtain a live token only when recording. Prefer the pre-fetched token
	// (no extra SSO call); fall back to email/password login if needed.
	var liveToken, liveDisplayName string
	if needsRecording {
		if token != "" {
			liveToken = token
			liveDisplayName = displayName
		} else if email != "" && password != "" {
			authClient := gc.NewClient("")
			if err := authClient.Login(email, password); err != nil {
				t.Fatalf("garmin login: %v", err)
			}
			liveToken = authClient.Token()
			liveDisplayName = authClient.DisplayName()
		}
	}

	mode := recorder.ModeReplayWithNewEpisodes
	if !needsRecording {
		mode = recorder.ModeReplayOnly
	}

	opts := []recorder.Option{
		recorder.WithMode(mode),
		recorder.WithMatcher(func(req *http.Request, i cassette.Request) bool {
			cu, err := url.Parse(i.URL)
			if err != nil {
				return false
			}
			return req.Method == i.Method && normaliseURL(req.URL) == normaliseURL(cu)
		}),
	}

	if liveDisplayName != "" {
		opts = append(opts, recorder.WithHook(func(i *cassette.Interaction) error {
			i.Request.Headers.Set("Authorization", "Bearer test")
			i.Request.URL = strings.ReplaceAll(i.Request.URL, url.PathEscape(liveDisplayName), "testuser")
			i.Request.URL = strings.ReplaceAll(i.Request.URL, liveDisplayName, "testuser")
			return nil
		}, recorder.BeforeSaveHook))
	}

	r, err := recorder.New(cassettePath, opts...)
	if err != nil {
		t.Fatalf("recorder.New: %v", err)
	}

	// Use real credentials when recording, synthetic ones when replaying.
	clientToken := "test"
	clientName := "testuser"
	if liveToken != "" {
		clientToken = liveToken
	}
	if liveDisplayName != "" {
		clientName = liveDisplayName
	}

	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: r}),
		gc.WithToken(clientToken),
		gc.WithDisplayName(clientName),
	)
	return c, func() {
		if err := r.Stop(); err != nil {
			t.Errorf("recorder.Stop: %v", err)
		}
	}
}

// newServerClient returns a Client pointed at srv with a pre-loaded token.
// The server is closed automatically via t.Cleanup.
func newServerClient(t *testing.T, srv *httptest.Server) *gc.Client {
	t.Helper()
	t.Cleanup(srv.Close)
	return gc.NewClient("",
		gc.WithBaseURL(srv.URL),
		gc.WithToken("test"),
		gc.WithDisplayName("testuser"),
	)
}

// fixedTransport returns a RoundTripper that always responds with the given
// status code and optional body, without making any real network calls.
func fixedTransport(code int, body string) http.RoundTripper {
	return roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: code,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func normaliseURL(u *url.URL) string {
	cp := *u
	cp.RawQuery = cp.Query().Encode()
	return cp.String()
}

// newAuthVCRClient returns a VCR-backed Client for auth flow tests.
// Unlike newVCRClient, no token is pre-loaded so the full auth flow runs through
// the cassette. Recording requires GARMIN_EMAIL and GARMIN_PASSWORD in env;
// if the account has MFA, the user is prompted on stdin for the code.
func newAuthVCRClient(t *testing.T, cassetteName string) (*gc.Client, func()) {
	t.Helper()

	cassettePath := "testdata/cassettes/" + cassetteName

	needsRecording := false
	if _, err := os.Stat(cassettePath + ".yaml"); os.IsNotExist(err) {
		needsRecording = true
	}

	if needsRecording && (os.Getenv("GARMIN_EMAIL") == "" || os.Getenv("GARMIN_PASSWORD") == "") {
		t.Skipf("cassette %q not found; set GARMIN_EMAIL and GARMIN_PASSWORD to record", cassetteName)
	}

	mode := recorder.ModeReplayWithNewEpisodes
	if !needsRecording {
		mode = recorder.ModeReplayOnly
	}

	opts := []recorder.Option{
		recorder.WithMode(mode),
		recorder.WithMatcher(func(req *http.Request, i cassette.Request) bool {
			cu, err := url.Parse(i.URL)
			if err != nil {
				return false
			}
			return req.Method == i.Method && normaliseURL(req.URL) == normaliseURL(cu)
		}),
	}

	// When recording, scrub Authorization headers before the cassette is saved.
	// Request and response bodies still need sanitize_cassettes.py to clean
	// passwords, tokens, and tickets.
	if needsRecording {
		opts = append(opts, recorder.WithHook(func(i *cassette.Interaction) error {
			i.Request.Headers.Del("Authorization")
			return nil
		}, recorder.BeforeSaveHook))
	}

	r, err := recorder.New(cassettePath, opts...)
	if err != nil {
		t.Fatalf("recorder.New: %v", err)
	}

	clientOpts := []gc.Option{gc.WithHTTPClient(&http.Client{Transport: r})}
	if needsRecording {
		clientOpts = append(clientOpts, gc.WithMFAPrompt(func() (string, error) {
			fmt.Fprint(os.Stderr, "Enter Garmin MFA code: ")
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return "", err
			}
			return strings.TrimSpace(line), nil
		}))
	} else {
		clientOpts = append(clientOpts, gc.WithMFAPrompt(func() (string, error) {
			return "123456", nil
		}))
	}

	c := gc.NewClient("", clientOpts...)
	return c, func() {
		if err := r.Stop(); err != nil {
			t.Errorf("recorder.Stop: %v", err)
		}
	}
}

// skipAPIError calls t.Skip when err signals a non-2xx response captured in the
// cassette (account doesn't have access to this endpoint).
func skipAPIError(t *testing.T, err error) {
	t.Helper()
	if errors.Is(err, gc.ErrUnauthorized) {
		t.Skipf("cassette captured 401 from API")
	}
	var ae *gc.APIError
	if errors.As(err, &ae) {
		t.Skipf("cassette captured HTTP %d from API", ae.StatusCode)
	}
}
