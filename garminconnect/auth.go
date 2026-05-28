package garminconnect

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ClientID     string    `json:"client_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (t *diToken) valid() bool {
	return t != nil && t.AccessToken != "" && time.Now().Before(t.ExpiresAt)
}

// Login ensures the client has a valid token, then fetches the user profile.
// It loads from disk, refreshes if needed, or performs a full SSO login as a
// last resort.
func (c *Client) Login(username, password string) error {
	if err := c.ensureToken(username, password); err != nil {
		return err
	}
	return c.fetchProfile()
}

func (c *Client) ensureToken(username, password string) error {
	if c.token.valid() {
		return nil
	}
	if tok, err := c.loadToken(); err == nil {
		if tok.valid() {
			c.token = tok
			return nil
		}
		if tok.RefreshToken != "" {
			if err := c.refreshToken(tok); err == nil {
				return nil
			}
		}
	}
	return c.ssoLogin(username, password)
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
	return os.WriteFile(c.tokenFile, data, 0600)
}

type ssoResponse struct {
	ServiceTicketID string `json:"serviceTicketId"`
	ServiceURL      string `json:"serviceURL"`
	ResponseStatus  struct {
		Type string `json:"type"`
	} `json:"responseStatus"`
	CustomerMfaInfo struct {
		MfaLastMethodUsed string `json:"mfaLastMethodUsed"`
	} `json:"customerMfaInfo"`
}

func (c *Client) ssoQueryParams() string {
	return fmt.Sprintf("?clientId=%s&locale=en-US&service=%s",
		ssoClientID, url.QueryEscape(ssoServiceURL))
}

func (c *Client) ssoLogin(username, password string) error {
	body, _ := json.Marshal(map[string]any{
		"username":     username,
		"password":     password,
		"rememberMe":   true,
		"captchaToken": "",
	})

	loginURL := ssoLoginURL + c.ssoQueryParams()

	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewReader(body))
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
		return fmt.Errorf("sso login: status %d", resp.StatusCode)
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
		return c.handleMFA(ssoResp.CustomerMfaInfo.MfaLastMethodUsed)
	}

	if ssoResp.ServiceTicketID == "" {
		return fmt.Errorf("sso login: no ticket in response (body: %s)", rawBody)
	}

	return c.exchangeTicket(ssoResp.ServiceTicketID, ssoResp.ServiceURL)
}

func (c *Client) handleMFA(mfaMethod string) error {
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
	req, err := http.NewRequest(http.MethodPost, verifyURL, bytes.NewReader(body))
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
		return fmt.Errorf("mfa verify: status %d", resp.StatusCode)
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

	return c.exchangeTicket(mfaResp.ServiceTicketID, mfaResp.ServiceURL)
}

func (c *Client) exchangeTicket(ticket, serviceURL string) error {
	for _, clientID := range diClientIDs {
		tok, err := c.doTokenRequest(url.Values{
			"client_id":      {clientID},
			"service_ticket": {ticket},
			"grant_type":     {"https://connectapi.garmin.com/di-oauth2-service/oauth/grant/service_ticket"},
			"service_url":    {serviceURL},
		}, clientID)
		if err == nil {
			tok.ClientID = clientID
			c.token = tok
			return c.saveToken(tok)
		}
	}
	return fmt.Errorf("di token exchange failed for all client IDs")
}

func (c *Client) refreshToken(old *diToken) error {
	tok, err := c.doTokenRequest(url.Values{
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
	c.token = tok
	return c.saveToken(tok)
}

func (c *Client) doTokenRequest(params url.Values, clientID string) (*diToken, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":"))
	req, err := http.NewRequest(http.MethodPost, diAuthURL, strings.NewReader(params.Encode()))
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
		return nil, fmt.Errorf("di token request: status %d", resp.StatusCode)
	}

	var raw struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("di token decode: %w", err)
	}

	expiry := time.Duration(raw.ExpiresIn)*time.Second - 60*time.Second
	return &diToken{
		AccessToken:  raw.AccessToken,
		RefreshToken: raw.RefreshToken,
		ExpiresAt:    time.Now().Add(expiry),
	}, nil
}
