// Package garminconnect is a Go client for the Garmin Connect API.
//
// Create a [Client] with [NewClient], authenticate with [Client.Login], then
// call any of the data methods. See the README for a capability overview.
package garminconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	utls "github.com/refraction-networking/utls"
)

const connectAPI = "https://connectapi.garmin.com"

// Client is an authenticated Garmin Connect API client. It is safe for
// concurrent use by multiple goroutines.
type Client struct {
	http        *http.Client
	mu          sync.Mutex // guards token
	refreshMu   sync.Mutex // serializes 401-triggered token refreshes
	token       *diToken
	tokenFile   string
	displayName string
	baseURL     string // defaults to connectAPI; overridable in tests
	mfaPrompt   func() (string, error)
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient replaces the default utls-backed HTTP client. Useful for
// injecting a custom transport, proxy, or test recorder.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithToken pre-loads an access token, skipping the SSO login flow.
// Unless a refresh token is also configured (WithRefreshToken), no refresh or
// SSO will be performed. If the token is a JWT its expiry is read from the
// "exp" claim; otherwise it is assumed valid for 24 hours.
func WithToken(accessToken string) Option {
	return func(c *Client) {
		expiresAt := time.Now().Add(24 * time.Hour)
		if exp, ok := jwtExpiry(accessToken); ok {
			expiresAt = exp
		}
		if c.token == nil {
			c.token = &diToken{}
		}
		c.token.AccessToken = accessToken
		c.token.ExpiresAt = expiresAt
	}
}

// WithDisplayName sets the Garmin Connect display name, skipping the profile
// fetch that Login normally performs.
func WithDisplayName(name string) Option {
	return func(c *Client) { c.displayName = name }
}

// WithRefreshToken attaches a refresh token to the client so that expired or
// revoked access tokens are automatically exchanged for a new one.
func WithRefreshToken(refreshToken string) Option {
	return func(c *Client) {
		if c.token == nil {
			c.token = &diToken{}
		}
		c.token.RefreshToken = refreshToken
		if c.token.ClientID == "" {
			c.token.ClientID = diClientIDs[0]
		}
	}
}

// WithTokenJSON pre-loads a full token from its JSON form — the same format
// written to the token file and returned by [Client.TokenJSON]. It lets
// callers persist tokens in an external store (e.g. a secret backend)
// instead of a file. Invalid or empty input is ignored, in which case Login
// falls back to the SSO flow.
func WithTokenJSON(data []byte) Option {
	return func(c *Client) {
		var tok diToken
		if err := json.Unmarshal(data, &tok); err != nil || (tok.AccessToken == "" && tok.RefreshToken == "") {
			return
		}
		if tok.ClientID == "" {
			tok.ClientID = diClientIDs[0]
		}
		c.token = &tok
	}
}

// WithMFAPrompt sets a callback that is invoked when Garmin's SSO requires
// multi-factor authentication. The callback should return the MFA code
// (e.g. read from stdin, a channel, or an HTTP handler). If no prompt is
// configured and MFA is required, Login returns ErrMFARequired.
func WithMFAPrompt(fn func() (string, error)) Option {
	return func(c *Client) { c.mfaPrompt = fn }
}

// WithBaseURL overrides the default API base URL. Primarily useful in tests
// to point the client at an httptest.Server.
func WithBaseURL(u string) Option {
	return func(c *Client) { c.baseURL = u }
}

// NewClient returns a Client that caches tokens at tokenFile.
func NewClient(tokenFile string, opts ...Option) *Client {
	c := &Client{
		http:      newUTLSClient(),
		tokenFile: tokenFile,
		baseURL:   connectAPI,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func newUTLSClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, _, _ := net.SplitHostPort(addr)
				conn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
				if err != nil {
					return nil, err
				}
				uconn := utls.UClient(conn, &utls.Config{ServerName: host}, utls.HelloAndroid_11_OkHttp)
				return uconn, uconn.Handshake()
			},
		},
	}
}

// DisplayName returns the authenticated user's Garmin Connect display name.
func (c *Client) DisplayName() string { return c.displayName }

// Token returns the current access token, or empty string if not authenticated.
func (c *Client) Token() string {
	if t := c.currentToken(); t != nil {
		return t.AccessToken
	}
	return ""
}

// TokenJSON returns the current token serialized in the token file format,
// for persistence outside the client (e.g. a secret store). Pass the result
// to [WithTokenJSON] to resume the session later. Returns ErrUnauthorized if
// the client has no token.
func (c *Client) TokenJSON() ([]byte, error) {
	tok := c.currentToken()
	if tok == nil {
		return nil, ErrUnauthorized
	}
	return json.Marshal(tok)
}

func (c *Client) currentToken() *diToken {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.token
}

func (c *Client) setToken(t *diToken) {
	c.mu.Lock()
	c.token = t
	c.mu.Unlock()
}

func (c *Client) fetchProfile(ctx context.Context) error {
	var profile struct {
		DisplayName string `json:"displayName"`
	}
	if err := c.get(ctx, "/userprofile-service/socialProfile", nil, &profile); err != nil {
		return fmt.Errorf("fetch profile: %w", err)
	}
	c.displayName = profile.DisplayName
	return nil
}

// withRefresh calls fn with the current access token, and if the first
// attempt returns a 401, tries to refresh the token and retries fn once.
// Returns (nil, ErrUnauthorized) if the client has no token, the refresh
// fails, or there is no refresh token available.
func (c *Client) withRefresh(ctx context.Context, fn func(accessToken string) (*http.Response, error)) (*http.Response, error) {
	tok := c.currentToken()
	if tok == nil {
		return nil, ErrUnauthorized
	}
	resp, err := fn(tok.AccessToken)
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}
	resp.Body.Close()
	if c.refreshAfter401(ctx, tok.AccessToken) != nil {
		return nil, ErrUnauthorized
	}
	return fn(c.Token())
}

// refreshAfter401 exchanges the refresh token after a request using usedToken
// got a 401. Concurrent callers are serialized; whoever arrives after a
// successful refresh sees a changed access token and skips the exchange.
func (c *Client) refreshAfter401(ctx context.Context, usedToken string) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()
	tok := c.currentToken()
	if tok == nil || tok.RefreshToken == "" {
		return ErrUnauthorized
	}
	if tok.AccessToken != usedToken {
		return nil
	}
	return c.refreshToken(ctx, tok)
}

// get performs an authenticated GET against the Garmin Connect API and
// JSON-decodes the response body into out.
func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	return c.getURL(ctx, c.baseURL+path, params, out)
}

func (c *Client) getURL(ctx context.Context, rawURL string, params url.Values, out any) error {
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	resp, err := c.withRefresh(ctx, func(token string) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		return c.http.Do(req)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimit
	default:
		return &APIError{StatusCode: resp.StatusCode, Path: rawURL}
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// doRequest executes a non-GET HTTP request against the Garmin Connect API.
func (c *Client) doRequest(ctx context.Context, method, rawURL string, body any, out any) error {
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	resp, err := c.withRefresh(ctx, func(token string) (*http.Response, error) {
		var br io.Reader
		if data != nil {
			br = bytes.NewReader(data)
		}
		req, err := http.NewRequestWithContext(ctx, method, rawURL, br)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		if data != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return c.http.Do(req)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimit
	default:
		return &APIError{StatusCode: resp.StatusCode, Path: rawURL}
	}
	if out == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) post(ctx context.Context, path string, body any, out any) error {
	return c.doRequest(ctx, http.MethodPost, c.baseURL+path, body, out)
}

func (c *Client) put(ctx context.Context, path string, body any, out any) error {
	return c.doRequest(ctx, http.MethodPut, c.baseURL+path, body, out)
}

func (c *Client) del(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, c.baseURL+path, nil, nil)
}

// getBytes performs an authenticated GET and returns the raw response body.
// Used for binary downloads (FIT, GPX, TCX, etc.).
func (c *Client) getBytes(ctx context.Context, path string, params url.Values) ([]byte, error) {
	rawURL := c.baseURL + path
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	resp, err := c.withRefresh(ctx, func(token string) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return c.http.Do(req)
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusTooManyRequests:
		return nil, ErrRateLimit
	default:
		return nil, &APIError{StatusCode: resp.StatusCode, Path: rawURL}
	}
	return io.ReadAll(resp.Body)
}

// upload sends a file as multipart/form-data to the given path and
// JSON-decodes the response into out (may be nil).
func (c *Client) upload(ctx context.Context, path string, data []byte, filename string, out any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return err
	}
	if _, err := fw.Write(data); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}

	bodyBytes := buf.Bytes()
	contentType := w.FormDataContentType()
	rawURL := c.baseURL + path

	resp, err := c.withRefresh(ctx, func(token string) (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", contentType)
		return c.http.Do(req)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusTooManyRequests:
		return ErrRateLimit
	default:
		return &APIError{StatusCode: resp.StatusCode, Path: rawURL}
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func date(t time.Time) string { return t.Format("2006-01-02") }

var errNoDisplayName = errors.New("display name not set: call Login or use WithDisplayName")

// displayNamePath returns the display name escaped for use as a URL path
// segment, or an error when the client has none (e.g. WithToken without
// WithDisplayName and no Login).
func (c *Client) displayNamePath() (string, error) {
	if c.displayName == "" {
		return "", errNoDisplayName
	}
	return url.PathEscape(c.displayName), nil
}
