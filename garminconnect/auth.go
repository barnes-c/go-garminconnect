package garminconnect

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ssoLoginURL     = "https://sso.garmin.com/mobile/api/login"
	ssoMFAVerifyURL = "https://sso.garmin.com/mobile/api/mfa/verifyCode"
	diAuthURL       = "https://diauth.garmin.com/di-oauth2-service/oauth/token"
	ssoClientID     = "GCM_IOS_DARK"
	ssoServiceURL   = "https://mobile.integration.garmin.com/gcm/ios"
	ssoUserAgent    = "GCM-iOS-5.7.2.1 (com.garmin.connect.mobile.sso)"
)

var diClientIDs = []string{
	"GARMIN_CONNECT_MOBILE_ANDROID_DI_2025Q2",
	"GARMIN_CONNECT_MOBILE_ANDROID_DI_2024Q4",
	"GARMIN_CONNECT_MOBILE_ANDROID_DI",
}

type diToken struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	ClientID         string    `json:"client_id"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at,omitzero"`
}

func (t *diToken) valid() bool {
	return t != nil && t.AccessToken != "" && time.Now().Before(t.ExpiresAt)
}

// jwtExpiry reads the "exp" claim from a JWT access token without verifying its
// signature. It returns false if the token is not a JWT or has no exp claim.
func jwtExpiry(token string) (time.Time, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}, false
	}
	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Exp == 0 {
		return time.Time{}, false
	}
	return time.Unix(claims.Exp, 0), true
}

// Login ensures the client has a valid token, then fetches the user profile.
// It loads from disk, refreshes if needed, or performs a full SSO login as a
// last resort.
func (c *Client) Login(ctx context.Context, username, password string) error {
	if err := c.ensureToken(ctx, username, password); err != nil {
		return err
	}
	return c.fetchProfile(ctx)
}

// Logout clears the in-memory token and removes the cached token file on disk,
// if one is configured. The next Login then runs a full SSO flow instead of
// resuming the cached token.
func (c *Client) Logout() error {
	c.setToken(nil)
	c.displayName = ""
	if c.tokenFile == "" {
		return nil
	}
	if err := os.Remove(c.tokenFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (c *Client) ensureToken(ctx context.Context, username, password string) error {
	if c.currentToken().valid() {
		return nil
	}
	var refreshErr error
	// An expired preset token (WithTokenJSON/WithRefreshToken) can still be
	// exchanged without a full SSO round trip.
	if tok := c.currentToken(); tok != nil && tok.RefreshToken != "" {
		if refreshErr = c.refreshToken(ctx, tok); refreshErr == nil {
			return nil
		}
	}
	if tok, err := c.loadToken(); err == nil {
		if tok.valid() {
			c.setToken(tok)
			return nil
		}
		if tok.RefreshToken != "" {
			if refreshErr = c.refreshToken(ctx, tok); refreshErr == nil {
				return nil
			}
		}
	}
	// A full SSO login needs credentials (and interactive MFA); when it fails
	// on a headless refresh, surface why the cached token could not be reused
	// instead of only the misleading SSO error.
	if err := c.ssoLogin(ctx, username, password); err != nil {
		if refreshErr != nil {
			return fmt.Errorf("%w (cached token refresh failed: %v)", err, refreshErr)
		}
		return err
	}
	return nil
}

func (c *Client) loadToken() (*diToken, error) {
	if c.tokenFile == "" {
		return nil, fmt.Errorf("no token file configured")
	}
	data, err := os.ReadFile(c.tokenFile)
	if err != nil {
		return nil, err
	}
	var tok diToken
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func (c *Client) saveToken(tok *diToken) error {
	if c.tokenFile == "" {
		return nil
	}
	data, err := json.Marshal(tok) //nolint:gosec // intentionally marshaling OAuth token to disk cache
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(c.tokenFile), ".garmin-token-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, c.tokenFile)
}

type ssoResponse struct {
	ServiceTicketID string             `json:"serviceTicketId"`
	ServiceURL      string             `json:"serviceURL"`
	ResponseStatus  ssoResponseStatus  `json:"responseStatus"`
	CustomerMfaInfo ssoCustomerMfaInfo `json:"customerMfaInfo"`
}

type ssoResponseStatus struct {
	Type string `json:"type"`
}

type ssoCustomerMfaInfo struct {
	MfaLastMethodUsed string `json:"mfaLastMethodUsed"`
}

func (c *Client) ssoQueryParams() string {
	return fmt.Sprintf("?clientId=%s&locale=en-US&service=%s",
		ssoClientID, url.QueryEscape(ssoServiceURL))
}

func (c *Client) ssoLogin(ctx context.Context, username, password string) error {
	body, _ := json.Marshal(map[string]any{
		"username":     username,
		"password":     password,
		"rememberMe":   true,
		"captchaToken": "",
	})

	loginURL := ssoLoginURL + c.ssoQueryParams()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", ssoUserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sso login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return authError("sso login", resp)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("sso login read body: %w", err)
	}

	var ssoResp ssoResponse
	if err := json.Unmarshal(rawBody, &ssoResp); err != nil {
		return fmt.Errorf("sso login decode: %w (body: %s)", err, rawBody)
	}

	if ssoResp.ResponseStatus.Type == "MFA_REQUIRED" {
		return c.handleMFA(ctx, ssoResp.CustomerMfaInfo.MfaLastMethodUsed)
	}

	// Garmin returns 200 with this status when its bot protection wants a
	// CAPTCHA solved in a browser; without this it surfaces as "no ticket".
	if ssoResp.ResponseStatus.Type == "CAPTCHA_REQUIRED" {
		return fmt.Errorf("sso login: CAPTCHA required — sign in at " +
			"https://connect.garmin.com in a browser to clear the challenge")
	}

	if ssoResp.ServiceTicketID == "" {
		return fmt.Errorf("sso login: no ticket in response (body: %s)", rawBody)
	}

	return c.exchangeTicket(ctx, ssoResp.ServiceTicketID, ssoResp.ServiceURL)
}

func (c *Client) handleMFA(ctx context.Context, mfaMethod string) error {
	if c.mfaPrompt == nil {
		return ErrMFARequired
	}
	if mfaMethod == "" {
		mfaMethod = "email"
	}

	code, err := c.mfaPrompt()
	if err != nil {
		return fmt.Errorf("mfa prompt: %w", err)
	}

	body, _ := json.Marshal(map[string]any{
		"mfaMethod":           mfaMethod,
		"mfaVerificationCode": code,
		"rememberMyBrowser":   true,
		"reconsentList":       []any{},
		"mfaSetup":            false,
	})

	verifyURL := ssoMFAVerifyURL + c.ssoQueryParams()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mfa verify request: %w", err)
	}
	req.Header.Set("User-Agent", ssoUserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mfa verify: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return authError("mfa verify", resp)
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mfa verify read body: %w", err)
	}

	var mfaResp ssoResponse
	if err := json.Unmarshal(rawBody, &mfaResp); err != nil {
		return fmt.Errorf("mfa verify decode: %w (body: %s)", err, rawBody)
	}

	if mfaResp.ServiceTicketID == "" {
		return fmt.Errorf("mfa verify: no ticket in response (body: %s)", rawBody)
	}

	return c.exchangeTicket(ctx, mfaResp.ServiceTicketID, mfaResp.ServiceURL)
}

func (c *Client) exchangeTicket(ctx context.Context, ticket, serviceURL string) error {
	var errs []error
	for _, clientID := range diClientIDs {
		tok, err := c.doTokenRequest(ctx, url.Values{
			"client_id":      {clientID},
			"service_ticket": {ticket},
			"grant_type":     {"https://connectapi.garmin.com/di-oauth2-service/oauth/grant/service_ticket"},
			"service_url":    {serviceURL},
		}, clientID)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", clientID, err))
			continue
		}
		tok.ClientID = clientID
		c.setToken(tok)
		return c.saveToken(tok)
	}
	return fmt.Errorf("di token exchange failed for all client IDs: %w", errors.Join(errs...))
}

func (c *Client) refreshToken(ctx context.Context, old *diToken) error {
	tok, err := c.doTokenRequest(ctx, url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {old.ClientID},
		"refresh_token": {old.RefreshToken},
	}, old.ClientID)
	if err != nil {
		return err
	}
	tok.ClientID = old.ClientID
	if tok.RefreshToken == "" {
		tok.RefreshToken = old.RefreshToken
	}
	// Garmin does not always return refresh_token_expires_in on a refresh, even
	// when it rotates the refresh token. Carrying the old expiry forward keeps a
	// (conservative) deadline on disk instead of dropping the field entirely.
	if tok.RefreshExpiresAt.IsZero() {
		tok.RefreshExpiresAt = old.RefreshExpiresAt
	}
	c.setToken(tok)
	return c.saveToken(tok)
}

func (c *Client) doTokenRequest(ctx context.Context, params url.Values, clientID string) (*diToken, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":"))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, diAuthURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, authError("di token request", resp)
	}

	var raw struct {
		AccessToken           string `json:"access_token"`
		RefreshToken          string `json:"refresh_token"`
		ExpiresIn             int    `json:"expires_in"`
		RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("di token decode: %w", err)
	}

	// Expire a minute early so a token is refreshed before the server rejects
	// it, but never let the buffer push the expiry into the past.
	expiry := time.Duration(raw.ExpiresIn) * time.Second
	if expiry > time.Minute {
		expiry -= time.Minute
	}
	tok := &diToken{
		AccessToken:  raw.AccessToken,
		RefreshToken: raw.RefreshToken,
		ExpiresAt:    time.Now().Add(expiry),
	}
	if raw.RefreshTokenExpiresIn > 0 {
		tok.RefreshExpiresAt = time.Now().Add(time.Duration(raw.RefreshTokenExpiresIn) * time.Second)
	}
	return tok, nil
}

// maxAuthErrorBody caps how much of a failed auth response is kept in the
// error message.
const maxAuthErrorBody = 512

// authError converts a non-200 response from an auth endpoint into an error
// that says why. Garmin signals rate limiting and bot challenges only through
// the status, the Retry-After header, and the body, so all three are preserved
// — a 429 additionally wraps ErrRateLimit so callers can test for it.
func authError(stage string, resp *http.Response) error {
	snippet := bodySnippet(resp.Body)

	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		detail := "rate limited by Garmin"
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			detail += ", Retry-After: " + ra
		}
		return fmt.Errorf("%s: %w (%s)%s", stage, ErrRateLimit, detail, snippet)
	case isBotChallenge(resp, snippet):
		return fmt.Errorf("%s: blocked by bot protection (status %d); "+
			"the TLS fingerprint or client identity was rejected%s", stage, resp.StatusCode, snippet)
	default:
		return fmt.Errorf("%s: status %d%s", stage, resp.StatusCode, snippet)
	}
}

// bodySnippet returns a bounded, single-line rendering of an error response
// body, or "" when there is nothing to show.
func bodySnippet(r io.Reader) string {
	data, err := io.ReadAll(io.LimitReader(r, maxAuthErrorBody))
	if err != nil || len(data) == 0 {
		return ""
	}
	s := strings.Join(strings.Fields(string(data)), " ")
	if s == "" {
		return ""
	}
	return " (body: " + s + ")"
}

// isBotChallenge reports whether a response looks like a Cloudflare or WAF
// interstitial rather than a Garmin API error.
func isBotChallenge(resp *http.Response, snippet string) bool {
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusServiceUnavailable {
		return false
	}
	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return true
	}
	lower := strings.ToLower(snippet)
	for _, marker := range []string{"cloudflare", "cf-ray", "just a moment", "attention required"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}
