# go-garminconnect

[![GitHub Release](https://img.shields.io/github/v/release/barnes-c/go-garminconnect)](https://github.com/barnes-c/go-garminconnect/releases/latest)
[![Build Status](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml/badge.svg)](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml)
[![golangci-lint](https://github.com/barnes-c/go-garminconnect/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/barnes-c/go-garminconnect/actions/workflows/golangci-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/barnes-c/go-garminconnect.svg)](https://pkg.go.dev/github.com/barnes-c/go-garminconnect/garminconnect)

Go client library for the Garmin Connect API.

## Installation

```bash
go get github.com/barnes-c/go-garminconnect
```

Requires Go 1.24+.

## Quick start

```go
ctx := context.Background()
client := garminconnect.NewClient(os.Getenv("HOME") + "/.garminconnect/tokens.json")
if err := client.Login(ctx, "user@example.com", "password"); err != nil {
    log.Fatal(err)
}

summary, err := client.UserSummary(ctx,time.Now())
if err != nil {
 log.Fatal(err)
}
fmt.Printf("Steps today: %d\n", summary.TotalSteps)
```

## Authentication

`Login` loads a cached token from disk, refreshes it if expired, or performs a full SSO login.

If the account has MFA enabled, provide a callback via `WithMFAPrompt` that returns the verification code. Without it, `Login` returns `ErrMFARequired`.

```go
client := garminconnect.NewClient(os.Getenv("HOME")+"/.garminconnect/tokens.json",
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
```

Other options: `WithHTTPClient(hc)`, `WithToken(accessToken)`, `WithDisplayName(name)`.

## API

The full method reference lives on [**pkg.go.dev**](https://pkg.go.dev/github.com/barnes-c/go-garminconnect/garminconnect), generated from the source. The client covers:

|         Area         |                                                                       What it exposes                                                                        |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Wellness**         | Steps, Body Battery, stress, floors, hydration, respiration, SpO₂, intensity minutes, blood pressure                                                         |
| **Heart rate**       | Intraday and resting heart rate                                                                                                                              |
| **Sleep**            | Sleep stages, sleep statistics, and HRV                                                                                                                      |
| **Activities**       | Listing, detail, splits, HR/power zones, exercise sets, weather, personal records, edit/delete, and download & upload                                        |
| **Workouts**         | Saved workouts, scheduling, and download & upload                                                                                                            |
| **Training**         | Readiness & status, VO₂ Max, endurance/hill scores, race predictions, running tolerance, lactate threshold, fitness age, cycling FTP, HR/power zone settings |
| **Body composition** | Weigh-ins and body composition history                                                                                                                       |
| **Goals**            | Goals, badges, and challenges                                                                                                                                |
| **Devices & gear**   | Registered devices, settings, solar data, and gear tracking                                                                                                  |
| **Profile**          | User profile, settings, and unit system                                                                                                                      |
| **Women's health**   | Menstrual cycle and pregnancy data                                                                                                                           |
| **Nutrition**        | Food log, meals, and goals                                                                                                                                   |
| **Golf**             | Scorecards and shot data                                                                                                                                     |

Activity and workout files support FIT, GPX, TCX, KML, and CSV (see the `Format*` constants).

## Error handling

```go
import "errors"

acts, err := client.Activities(context.Background(), 10)
switch {
case errors.Is(err, garminconnect.ErrUnauthorized):
    // automatic token refresh failed — call Login again to re-authenticate
case errors.Is(err, garminconnect.ErrMFARequired):
    // account requires MFA — provide WithMFAPrompt on client creation
case errors.Is(err, garminconnect.ErrRateLimit):
    // back off and retry
case errors.Is(err, garminconnect.ErrNoData):
    // no records for the query
}

var apiErr *garminconnect.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.StatusCode, apiErr.Path)
}
```

## Testing

Tests replay recorded HTTP interactions — no credentials needed.

```bash
make test
```

To re-record cassettes against a live account:

```bash
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./internal/record
# Only record missing cassettes:
go run ./internal/record --missing
```

The tool logs in once and records one cassette per test. PII is stripped **inline** as each cassette is written (the recorder's `BeforeSaveHook` calls `internal/sanitize`), so cassettes are never persisted unsanitized.

A scheduled [smoke-test workflow](.github/workflows/smoke.yml) runs `internal/smoke` weekly against the live API to catch drift the cassettes can't see. It can also be run locally:

```bash
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret go run ./internal/smoke
```
