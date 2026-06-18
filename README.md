# go-garminconnect

[![Build Status](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml/badge.svg)](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml)
[![golangci-lint](https://github.com/barnes-c/go-garminconnect/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/barnes-c/go-garminconnect/actions/workflows/golangci-lint.yml)
[![GitHub Release](https://img.shields.io/github/v/release/barnes-c/go-garminconnect)](https://github.com/barnes-c/go-garminconnect/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/barnes-c/go-garminconnect)](https://goreportcard.com/report/github.com/barnes-c/go-garminconnect)

Go client library for the Garmin Connect API.

## Installation

```bash
go get github.com/barnes-c/go-garminconnect
```

Requires Go 1.25+.

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

### Wellness & daily health

| Method                                              | Description                               |
|-----------------------------------------------------|-------------------------------------------|
| `UserSummary(ctx, date)`                            | Steps, calories, active minutes, stress   |
| `Steps(ctx, date)`                                  | Intraday step entries                     |
| `StepsData(ctx, date)`                              | Intraday step data in 15-minute intervals |
| `DailySteps(ctx, start, end)`                       | Daily step totals over a date range       |
| `WeeklySteps(ctx, end, weeks)`                      | Weekly step aggregates                    |
| `BodyBattery(ctx, start, end)`                      | Body Battery readings                     |
| `BodyBatteryEvents(ctx, date)`                      | Body Battery charge/drain events          |
| `AllDayStress(ctx, date)`                           | Stress measurements throughout the day    |
| `WeeklyStress(ctx, end, weeks)`                     | Weekly stress aggregates                  |
| `Floors(ctx, date)`                                 | Floors climbed                            |
| `Hydration(ctx, date)`                              | Hydration log                             |
| `AddHydration(ctx, ml, timestamp, date)`            | Log a hydration entry                     |
| `Respiration(ctx, date)`                            | Respiration rate                          |
| `SpO2(ctx, date)`                                   | Blood oxygen saturation                   |
| `IntensityMinutes(ctx, date)`                       | Moderate and vigorous intensity minutes   |
| `WeeklyIntensityMinutes(ctx, start, end)`           | Weekly intensity minute aggregates        |
| `BloodPressure(ctx, start, end)`                    | Blood pressure readings                   |
| `SetBloodPressure(ctx, sys, dia, pulse, ts, notes)` | Log a blood pressure reading              |
| `DeleteBloodPressure(ctx, date, version)`           | Delete a blood pressure entry             |
| `AllDayEvents(ctx, date)`                           | All wellness events for a day             |
| `LifestyleData(ctx, date)`                          | Lifestyle summary                         |

### Heart rate

| Method                              | Description                  |
|-------------------------------------|------------------------------|
| `HeartRates(ctx, date)`             | Intraday heart rate readings |
| `RestingHeartRate(ctx, start, end)` | Resting HR over a date range |

### Sleep

| Method                        | Description                                       |
|-------------------------------|---------------------------------------------------|
| `SleepData(ctx, date)`        | Sleep stages and quality for a night              |
| `DailySleepData(ctx, date)`   | Sleep data for a night (sleep-service endpoint)   |
| `SleepStats(ctx, start, end)` | Daily sleep statistics over a range (max 28 days) |
| `HRVData(ctx, date)`          | HRV measurements during sleep                     |

### Activities

| Method                                            | Description                                      |
|---------------------------------------------------|--------------------------------------------------|
| `Activities(ctx, limit)`                          | Most recent N activities                         |
| `ActivitiesByDate(ctx, start, end, type)`         | Activities in a date range, optional type filter |
| `LastActivity(ctx)`                               | Single most recent activity                      |
| `ActivityCount(ctx)`                              | Total activity count                             |
| `ActivityDetail(ctx, id)`                         | Full activity detail                             |
| `ActivitySplits(ctx, id)`                         | Lap/split summaries                              |
| `ActivityTypedSplits(ctx, id)`                    | Sport-specific split data                        |
| `ActivitySplitSummaries(ctx, id)`                 | Split summary statistics                         |
| `ActivityHRZones(ctx, id)`                        | Time in heart rate zones                         |
| `ActivityPowerZones(ctx, id)`                     | Time in power zones                              |
| `ActivityExerciseSets(ctx, id)`                   | Strength training exercise sets                  |
| `ActivityWeather(ctx, id)`                        | Weather recorded during the activity             |
| `PersonalRecords(ctx)`                            | Personal bests                                   |
| `SetActivityName(ctx, id, name)`                  | Rename an activity                               |
| `SetActivityType(ctx, id, typeID, parentID, key)` | Change activity sport type                       |
| `DeleteActivity(ctx, id)`                         | Delete an activity                               |
| `DownloadActivity(ctx, id, format)`               | Download FIT, GPX, TCX, KML, or CSV              |
| `UploadActivity(ctx, data, filename)`             | Upload a FIT, GPX, or TCX file                   |

Download format constants: `FormatOriginal`, `FormatTCX`, `FormatGPX`, `FormatKML`, `FormatCSV`.

### Workouts

| Method                                  | Description                        |
|-----------------------------------------|------------------------------------|
| `Workouts(ctx, start, limit)`           | Saved workouts                     |
| `Workout(ctx, id)`                      | Single workout detail              |
| `DeleteWorkout(ctx, id)`                | Delete a workout                   |
| `ScheduledWorkouts(ctx, year, month)`   | Scheduled workout calendar entries |
| `ScheduleWorkout(ctx, workoutID, date)` | Add a workout to the calendar      |
| `UnscheduleWorkout(ctx, scheduledID)`   | Remove a workout from the calendar |
| `DownloadWorkout(ctx, id)`              | Download workout as FIT            |
| `UploadWorkout(ctx, data, filename)`    | Upload a workout file              |

### Training metrics

| Method                              | Description                   |
|-------------------------------------|-------------------------------|
| `TrainingReadiness(ctx, date)`      | Training readiness score      |
| `TrainingStatus(ctx, date)`         | Training status and load      |
| `MaxMetrics(ctx, start, end)`       | VO2 Max and other max metrics |
| `EnduranceScore(ctx, start, end)`   | Endurance score trend         |
| `HillScore(ctx, start, end)`        | Hill score trend              |
| `RacePredictions(ctx)`              | Predicted race finish times   |
| `LactateThreshold(ctx)`             | Lactate threshold data        |
| `FitnessAge(ctx, date)`             | Fitness age estimate          |
| `RunningTolerance(ctx, start, end)` | Running load tolerance        |
| `CyclingFTP(ctx)`                   | Functional threshold power    |

### Body composition

| Method                                    | Description                              |
|-------------------------------------------|------------------------------------------|
| `BodyComposition(ctx, start, end)`        | Weight and body composition over a range |
| `WeighIns(ctx, start, end)`               | All weigh-in entries in a range          |
| `DailyWeighIns(ctx, date)`                | Weigh-ins for a single day               |
| `AddWeighIn(ctx, kg, unitKey, timestamp)` | Log a weigh-in                           |
| `DeleteWeighIn(ctx, date, weightPK)`      | Delete a weigh-in entry                  |

### Goals & achievements

| Method                                           | Description               |
|--------------------------------------------------|---------------------------|
| `Goals(ctx, status, start, limit)`               | Goals filtered by status  |
| `EarnedBadges(ctx)`                              | Earned badges             |
| `AvailableBadges(ctx)`                           | Available badges          |
| `AdHocChallenges(ctx, start, limit)`             | Ad-hoc challenges         |
| `BadgeChallenges(ctx, start, limit)`             | Badge challenges          |
| `AvailableBadgeChallenges(ctx, start, limit)`    | Joinable badge challenges |
| `InProgressVirtualChallenges(ctx, start, limit)` | Active virtual challenges |

### Devices & gear

| Method                                        | Description                    |
|-----------------------------------------------|--------------------------------|
| `Devices(ctx)`                                | Registered devices             |
| `DeviceSettings(ctx, deviceID)`               | Settings for a specific device |
| `LastUsedDevice(ctx)`                         | Most recently synced device    |
| `PrimaryTrainingDevice(ctx)`                  | Primary training device        |
| `DeviceSolarData(ctx, deviceID, start, end)`  | Solar charging data            |
| `Gear(ctx, userProfileNumber)`                | Gear items                     |
| `GearStats(ctx, gearUUID)`                    | Usage stats for a gear item    |
| `GearActivities(ctx, gearUUID, start, limit)` | Activities using a gear item   |
| `GearDefaults(ctx, userProfileNumber)`        | Default gear assignments       |

### Profile

| Method                     | Description                             |
|----------------------------|-----------------------------------------|
| `UserProfile(ctx)`         | Display name, location, join date       |
| `UserProfileSettings(ctx)` | Account settings                        |
| `DisplayName(ctx)`         | Cached display name (ctx, set on login) |

### Women's health

| Method                               | Description                    |
|--------------------------------------|--------------------------------|
| `MenstrualData(ctx, date)`           | Menstrual cycle data for a day |
| `MenstrualCalendar(ctx, start, end)` | Cycle data over a date range   |
| `PregnancySummary(ctx)`              | Pregnancy snapshot             |

### Nutrition

| Method                         | Description             |
|--------------------------------|-------------------------|
| `NutritionFoodLog(ctx, date)`  | Food log for a day      |
| `NutritionMeals(ctx, date)`    | Meal breakdown          |
| `NutritionSettings(ctx, date)` | Nutrition goal settings |

### Golf

| Method                                  | Description                          |
|-----------------------------------------|--------------------------------------|
| `GolfSummary(ctx, start, limit)`        | Scorecard summaries                  |
| `GolfScorecard(ctx, scorecardID)`       | Full scorecard detail                |
| `GolfShotData(ctx, scorecardID, holes)` | Shot-by-shot data for selected holes |

## Error handling

```go
import "errors"

acts, err := client.Activities(context.Background(), 10)
switch {
case errors.Is(err, garminconnect.ErrUnauthorized):
    // token expired — call Login again
case errors.Is(err, garminconnect.ErrMFARequired):
    // account requires MFA — provide WithMFAPrompt on client creation
case errors.Is(err, garminconnect.ErrRateLimit):
    // back off and retry
case errors.Is(err, garminconnect.ErrNoData):
    // no records for the query
}

if apiErr, ok := errors.AsType[*garminconnect.APIError](err); ok {
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
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret bash tools/record_cassettes.sh
# Only record missing cassettes:
bash tools/record_cassettes.sh --missing
```

The script logs in once, records one cassette per test, then runs `tools/sanitize_cassettes.py` to strip PII before committing.
