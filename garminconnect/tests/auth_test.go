package garminconnect_test

import (
	"net/http"
	"net/url"
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
