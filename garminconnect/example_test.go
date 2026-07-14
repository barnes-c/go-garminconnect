package garminconnect_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/barnes-c/go-garminconnect/garminconnect"
)

// Log in and fetch today's activity summary. Login loads a cached token from
// tokenFile, refreshes it if expired, or performs a full SSO login.
func Example() {
	ctx := context.Background()
	tokenFile := filepath.Join(os.Getenv("HOME"), ".garminconnect", "tokens.json")

	client := garminconnect.NewClient(tokenFile)
	if err := client.Login(ctx, os.Getenv("GARMIN_EMAIL"), os.Getenv("GARMIN_PASSWORD")); err != nil {
		log.Fatal(err)
	}

	summary, err := client.UserSummary(ctx, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Steps today: %d\n", summary.TotalSteps)
}

// Accounts with multi-factor authentication need a prompt callback that
// returns the verification code; without one, Login returns ErrMFARequired.
func ExampleWithMFAPrompt() {
	client := garminconnect.NewClient("tokens.json",
		garminconnect.WithMFAPrompt(func() (string, error) {
			fmt.Print("MFA code: ")
			var code string
			_, err := fmt.Scan(&code)
			return code, err
		}),
	)
	if err := client.Login(context.Background(), "user@example.com", "password"); err != nil {
		log.Fatal(err)
	}
}

// List recent activities, distinguishing error kinds with errors.Is and
// errors.As.
func ExampleClient_Activities() {
	client := garminconnect.NewClient("tokens.json")

	activities, err := client.Activities(context.Background(), 0, 10)
	switch {
	case errors.Is(err, garminconnect.ErrUnauthorized):
		log.Fatal("token refresh failed — call Login again")
	case errors.Is(err, garminconnect.ErrRateLimit):
		log.Fatal("rate limited — back off and retry")
	case err != nil:
		var apiErr *garminconnect.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("garmin returned %d for %s", apiErr.StatusCode, apiErr.Path)
		}
		log.Fatal(err)
	}
	for _, a := range activities {
		fmt.Println(a.ActivityName)
	}
}

// Download the most recent activity as a GPX file.
func ExampleClient_DownloadActivity() {
	ctx := context.Background()
	client := garminconnect.NewClient("tokens.json")

	last, err := client.LastActivity(ctx)
	if err != nil {
		log.Fatal(err)
	}
	data, err := client.DownloadActivity(ctx, last.ActivityID, garminconnect.FormatGPX)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("activity.gpx", data, 0o644); err != nil {
		log.Fatal(err)
	}
}
