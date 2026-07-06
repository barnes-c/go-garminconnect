// Command smoke runs a live smoke test against the Garmin Connect API. It
// logs in, calls one read-only endpoint per API area, and exits non-zero if
// any check fails — catching drift in the private API that the recorded
// cassettes can't see.
//
// A 4xx response is reported as SKIP, mirroring skipAPIError in the test
// suite: the endpoint isn't available on this account, which is not drift.
//
// Usage:
//
//	GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./internal/smoke
//
// Credentials may be omitted if a valid cached token exists in the token file.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"
	"github.com/barnes-c/go-garminconnect/internal/login"
)

const (
	tokenFile = ".garmin_token.json"
	delay     = time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	warnRefreshExpiry()
	c, err := login.Client(tokenFile)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	// The display name is PII and appears in API URLs (and therefore in
	// APIError messages) — have GitHub redact it from public logs.
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		fmt.Printf("::add-mask::%s\n", c.DisplayName())
	}
	fmt.Println("==> Logged in")

	ctx := context.Background()
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	monthAgo := today.AddDate(0, -1, 0)

	checks := []struct {
		name string
		fn   func() error
	}{
		{"UserProfile", func() error { _, err := c.UserProfile(ctx); return err }},
		{"UserSummary", func() error { _, err := c.UserSummary(ctx, yesterday); return err }},
		{"Steps", func() error { _, err := c.Steps(ctx, yesterday); return err }},
		{"HeartRates", func() error { _, err := c.HeartRates(ctx, yesterday); return err }},
		{"SleepData", func() error { _, err := c.SleepData(ctx, yesterday); return err }},
		{"Activities", func() error { _, err := c.Activities(ctx, 0, 1); return err }},
		{"ActivityTypes", func() error { _, err := c.ActivityTypes(ctx); return err }},
		{"Devices", func() error { _, err := c.Devices(ctx); return err }},
		{"Workouts", func() error { _, err := c.Workouts(ctx, 0, 1); return err }},
		{"BodyComposition", func() error { _, err := c.BodyComposition(ctx, monthAgo, today); return err }},
		{"TrainingStatus", func() error { _, err := c.TrainingStatus(ctx, today); return err }},
		{"EarnedBadges", func() error { _, err := c.EarnedBadges(ctx); return err }},
	}

	var failed int
	for i, check := range checks {
		err := check.fn()
		switch {
		case err == nil:
			fmt.Printf("--- PASS %s\n", check.name)
		case skippable(err):
			fmt.Printf("--- SKIP %s: %v\n", check.name, err)
		default:
			failed++
			fmt.Printf("--- FAIL %s: %v\n", check.name, err)
		}
		if i < len(checks)-1 {
			time.Sleep(delay)
		}
	}
	if failed > 0 {
		return fmt.Errorf("%d check(s) failed", failed)
	}
	return nil
}

// warnRefreshExpiry inspects the cached token before login rotates it. In CI
// the file is seeded from the GARMIN_TOKEN_JSON secret, whose refresh token
// expires 30 days after it was minted locally — warn while there is still
// time to re-seed it.
func warnRefreshExpiry() {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return
	}
	var tok struct {
		RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	}
	if json.Unmarshal(data, &tok) != nil || tok.RefreshExpiresAt.IsZero() {
		return
	}
	if left := time.Until(tok.RefreshExpiresAt); left < 10*24*time.Hour {
		fmt.Printf("::warning::Garmin refresh token expires in %.0f days — re-seed the GARMIN_TOKEN_JSON secret (log in locally, then: gh secret set GARMIN_TOKEN_JSON < .garmin_token.json)\n", left.Hours()/24)
	}
}

func skippable(err error) bool {
	if errors.Is(err, gc.ErrNoData) {
		return true
	}
	var apiErr *gc.APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode >= 400 && apiErr.StatusCode < 500
}
