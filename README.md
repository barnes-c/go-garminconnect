# go-garminconnect

[![Garmin Connect](https://img.shields.io/badge/Garmin%20Connect-007CC3?logo=garmin&logoColor=white)](https://connect.garmin.com)
[![GitHub Release](https://img.shields.io/github/v/release/barnes-c/go-garminconnect)](https://github.com/barnes-c/go-garminconnect/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/barnes-c/go-garminconnect.svg)](https://pkg.go.dev/github.com/barnes-c/go-garminconnect/garminconnect)
[![CI](https://img.shields.io/github/actions/workflow/status/barnes-c/go-garminconnect/ci.yml?branch=main&label=CI&logo=github)](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml)

Go client for the Garmin Connect API.

```bash
go get github.com/barnes-c/go-garminconnect
```

## Usage

```go
ctx := context.Background()
client := garminconnect.NewClient(os.Getenv("HOME") + "/.garminconnect/tokens.json")
if err := client.Login(ctx, "user@example.com", "password"); err != nil {
    log.Fatal(err)
}

summary, err := client.UserSummary(ctx, time.Now())
if err != nil {
    log.Fatal(err)
}
fmt.Println("Steps today:", summary.TotalSteps)
```

`Login` reuses the cached token, refreshes it when expired, and only falls back
to a full SSO login when needed. For accounts with MFA, pass
`garminconnect.WithMFAPrompt(func() (string, error) { ... })` to `NewClient`.

It covers wellness, sleep, heart rate, activities, workouts, courses, training
metrics and plans, body composition, goals and badges, devices and gear,
profile, women's health, nutrition, golf and reports. See
[pkg.go.dev](https://pkg.go.dev/github.com/barnes-c/go-garminconnect/garminconnect)
for the full API.

## Used by

- [garmin_exporter](https://github.com/barnes-c/garmin_exporter): OTel-native Prometheus exporter for Garmin Connect
- [garmin-datasource](https://github.com/barnes-c/garmin-datasource): Grafana data source for Garmin Connect
