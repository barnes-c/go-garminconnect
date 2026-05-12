package garminconnect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	utls "github.com/refraction-networking/utls"
)

const connectAPI = "https://connectapi.garmin.com"

// Client is an authenticated Garmin Connect API client.
type Client struct {
	http        *http.Client
	token       *diToken
	tokenFile   string
	displayName string
	baseURL     string // defaults to connectAPI; overridable in tests
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient replaces the default utls-backed HTTP client. Useful for
// injecting a custom transport, proxy, or test recorder.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.http = hc }
}

// WithToken pre-loads an access token, skipping the SSO login flow.
// The token is assumed valid; no refresh or SSO will be performed.
func WithToken(accessToken string) Option {
	return func(c *Client) {
		c.token = &diToken{AccessToken: accessToken, ExpiresAt: time.Now().Add(24 * time.Hour)}
	}
}

// WithDisplayName sets the Garmin Connect display name, skipping the profile
// fetch that Login normally performs.
func WithDisplayName(name string) Option {
	return func(c *Client) { c.displayName = name }
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
	return &http.Client{
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
	if c.token == nil {
		return ""
	}
	return c.token.AccessToken
}

func (c *Client) fetchProfile() error {
	var profile struct {
		DisplayName string `json:"displayName"`
	}
	if err := c.get("/userprofile-service/socialProfile", nil, &profile); err != nil {
		return fmt.Errorf("fetch profile: %w", err)
	}
	c.displayName = profile.DisplayName
	return nil
}

// get performs an authenticated GET against the Garmin Connect API and
// JSON-decodes the response body into out.
func (c *Client) get(path string, params url.Values, out any) error {
	return c.getURL(c.baseURL+path, params, out)
}

func (c *Client) getURL(rawURL string, params url.Values, out any) error {
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
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
func (c *Client) doRequest(method, rawURL string, body any, out any) error {
	var br io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		br = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, rawURL, br)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
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

func (c *Client) post(path string, body any, out any) error {
	return c.doRequest(http.MethodPost, c.baseURL+path, body, out)
}

func (c *Client) put(path string, body any, out any) error {
	return c.doRequest(http.MethodPut, c.baseURL+path, body, out)
}

func (c *Client) del(path string) error {
	return c.doRequest(http.MethodDelete, c.baseURL+path, nil, nil)
}

// getBytes performs an authenticated GET and returns the raw response body.
// Used for binary downloads (FIT, GPX, TCX, etc.).
func (c *Client) getBytes(path string, params url.Values) ([]byte, error) {
	rawURL := c.baseURL + path
	if len(params) > 0 {
		rawURL += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := c.http.Do(req)
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
func (c *Client) upload(path string, data []byte, filename string, out any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return err
	}
	if _, err := fw.Write(data); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, c.baseURL+path, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
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
		return &APIError{StatusCode: resp.StatusCode, Path: c.baseURL + path}
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func date(t time.Time) string { return t.Format("2006-01-02") }
