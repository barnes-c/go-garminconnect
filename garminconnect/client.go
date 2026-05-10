// Copyright Christopher Barnes
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
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
}

// NewClient returns a Client that caches tokens at tokenFile.
func NewClient(tokenFile string) *Client {
	return &Client{
		http:      newUTLSClient(),
		tokenFile: tokenFile,
	}
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
	u := connectAPI + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, u, nil)
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
		return &APIError{StatusCode: resp.StatusCode, Path: path}
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func date(t time.Time) string { return t.Format("2006-01-02") }
