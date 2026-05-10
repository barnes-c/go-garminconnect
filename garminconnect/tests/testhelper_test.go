package garminconnect_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
	"gopkg.in/dnaeon/go-vcr.v2/cassette"
	"gopkg.in/dnaeon/go-vcr.v2/recorder"
)

// newVCRClient returns a Client wired to the named cassette. The returned stop
// function must be deferred by the caller.
func newVCRClient(t *testing.T, cassetteName string) (*gc.Client, func()) {
	t.Helper()
	r, err := recorder.New("testdata/cassettes/" + cassetteName)
	if err != nil {
		t.Fatalf("recorder.New: %v", err)
	}
	r.SetMatcher(func(req *http.Request, i cassette.Request) bool {
		cu, err := url.Parse(i.URL)
		if err != nil {
			return false
		}
		return req.Method == i.Method && normaliseURL(req.URL) == normaliseURL(cu)
	})
	c := gc.NewClient("",
		gc.WithHTTPClient(&http.Client{Transport: r}),
		gc.WithToken("test"),
		gc.WithDisplayName("testuser"),
	)
	return c, func() {
		if err := r.Stop(); err != nil {
			t.Errorf("recorder.Stop: %v", err)
		}
	}
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
	copy := *u
	copy.RawQuery = copy.Query().Encode()
	return copy.String()
}
