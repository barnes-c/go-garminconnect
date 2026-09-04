package garminconnect

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// The methods in this file are the escape hatch for endpoints this library
// does not wrap yet. They apply the same authentication, token refresh and
// error translation as the typed methods, but leave the request path and the
// response shape entirely to the caller.
//
// path is either relative to the Garmin Connect API base URL
// ("/wellness-service/wellness/dailyStress/2026-01-01") or an absolute URL, in
// which case it is used as-is — some Garmin payloads embed links to other
// hosts.

func (c *Client) resolve(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return c.baseURL + path
}

// Get performs an authenticated GET against path and JSON-decodes the response
// body into out. Pass a *json.RawMessage or map[string]any for an unmodelled
// endpoint, or nil to discard the body.
func (c *Client) Get(ctx context.Context, path string, params url.Values, out any) error {
	return c.getURL(ctx, c.resolve(path), params, out)
}

// GetBytes performs an authenticated GET against path and returns the raw
// response body, for endpoints that do not return JSON (FIT, GPX, TCX, CSV).
func (c *Client) GetBytes(ctx context.Context, path string, params url.Values) ([]byte, error) {
	return c.getBytesURL(ctx, c.resolve(path), params)
}

// Post performs an authenticated POST against path. body is JSON-encoded when
// non-nil; the response is JSON-decoded into out when non-nil.
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	return c.doRequest(ctx, http.MethodPost, c.resolve(path), body, out)
}

// Put performs an authenticated PUT against path. body is JSON-encoded when
// non-nil; the response is JSON-decoded into out when non-nil.
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	return c.doRequest(ctx, http.MethodPut, c.resolve(path), body, out)
}

// Delete performs an authenticated DELETE against path.
func (c *Client) Delete(ctx context.Context, path string) error {
	return c.doRequest(ctx, http.MethodDelete, c.resolve(path), nil, nil)
}

// Upload posts data to path as a multipart/form-data file named filename,
// JSON-decoding the response into out when non-nil.
func (c *Client) Upload(ctx context.Context, path string, data []byte, filename string, out any) error {
	return c.uploadURL(ctx, c.resolve(path), data, filename, out)
}
