# go-garminconnect

[![Build Status](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml/badge.svg)](https://github.com/barnes-c/go-garminconnect/actions/workflows/ci.yml)
![golangci-lint](https://github.com/barnes-c/go-garminconnect/actions/workflows/golangci-lint.yml/badge.svg)
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
client := garminconnect.NewClient("~/.garminconnect/tokens.json")
if err := client.Login("user@example.com", "password"); err != nil {
    log.Fatal(err)
}

summary, err := client.UserSummary(time.Now())
fmt.Printf("Steps today: %d\n", summary.TotalSteps)
```

## Authentication

`Login` loads a cached token from disk, refreshes it if expired, or performs a full SSO login. MFA is not supported.

Options: `WithHTTPClient(hc)`, `WithToken(accessToken)`, `WithDisplayName(name)`.

## API

### Wellness & daily health

| Method | Description |
|---|---|
| `UserSummary(date)` | Steps, calories, active minutes, stress |
| `Steps(date)` | Intraday step entries |
| `DailySteps(start, end)` | Daily step totals over a date range |
| `WeeklySteps(end, weeks)` | Weekly step aggregates |
| `BodyBattery(start, end)` | Body Battery readings |
| `BodyBatteryEvents(date)` | Body Battery charge/drain events |
| `AllDayStress(date)` | Stress measurements throughout the day |
| `WeeklyStress(end, weeks)` | Weekly stress aggregates |
| `Floors(date)` | Floors climbed |
| `Hydration(date)` | Hydration log |
| `AddHydration(ml, timestamp, date)` | Log a hydration entry |
| `Respiration(date)` | Respiration rate |
| `SpO2(date)` | Blood oxygen saturation |
| `IntensityMinutes(date)` | Moderate and vigorous intensity minutes |
| `WeeklyIntensityMinutes(start, end)` | Weekly intensity minute aggregates |
| `BloodPressure(start, end)` | Blood pressure readings |
| `SetBloodPressure(sys, dia, pulse, ts, notes)` | Log a blood pressure reading |
| `DeleteBloodPressure(date, version)` | Delete a blood pressure entry |
| `AllDayEvents(date)` | All wellness events for a day |
| `LifestyleData(date)` | Lifestyle summary |

### Heart rate

| Method | Description |
|---|---|
| `HeartRates(date)` | Intraday heart rate readings |
| `RestingHeartRate(start, end)` | Resting HR over a date range |

### Sleep

| Method | Description |
|---|---|
| `SleepData(date)` | Sleep stages and quality for a night |
| `HRVData(date)` | HRV measurements during sleep |

### Activities

| Method | Description |
|---|---|
| `Activities(limit)` | Most recent N activities |
| `ActivitiesByDate(start, end, type)` | Activities in a date range, optional type filter |
| `LastActivity()` | Single most recent activity |
| `ActivityCount()` | Total activity count |
| `ActivityDetail(id)` | Full activity detail |
| `ActivitySplits(id)` | Lap/split summaries |
| `ActivityTypedSplits(id)` | Sport-specific split data |
| `ActivitySplitSummaries(id)` | Split summary statistics |
| `ActivityHRZones(id)` | Time in heart rate zones |
| `ActivityPowerZones(id)` | Time in power zones |
| `ActivityExerciseSets(id)` | Strength training exercise sets |
| `ActivityWeather(id)` | Weather recorded during the activity |
| `PersonalRecords()` | Personal bests |
| `SetActivityName(id, name)` | Rename an activity |
| `SetActivityType(id, typeID, parentID, key)` | Change activity sport type |
| `DeleteActivity(id)` | Delete an activity |
| `DownloadActivity(id, format)` | Download FIT, GPX, TCX, KML, or CSV |
| `UploadActivity(data, filename)` | Upload a FIT, GPX, or TCX file |

Download format constants: `FormatOriginal`, `FormatTCX`, `FormatGPX`, `FormatKML`, `FormatCSV`.

### Workouts

| Method | Description |
|---|---|
| `Workouts(start, limit)` | Saved workouts |
| `Workout(id)` | Single workout detail |
| `DeleteWorkout(id)` | Delete a workout |
| `ScheduledWorkouts(year, month)` | Scheduled workout calendar entries |
| `ScheduleWorkout(workoutID, date)` | Add a workout to the calendar |
| `UnscheduleWorkout(scheduledID)` | Remove a workout from the calendar |
| `DownloadWorkout(id)` | Download workout as FIT |
| `UploadWorkout(data, filename)` | Upload a workout file |

### Training metrics

| Method | Description |
|---|---|
| `TrainingReadiness(date)` | Training readiness score |
| `TrainingStatus(date)` | Training status and load |
| `MaxMetrics(start, end)` | VO2 Max and other max metrics |
| `EnduranceScore(start, end)` | Endurance score trend |
| `HillScore(start, end)` | Hill score trend |
| `RacePredictions()` | Predicted race finish times |
| `LactateThreshold()` | Lactate threshold data |
| `FitnessAge(date)` | Fitness age estimate |
| `RunningTolerance(start, end)` | Running load tolerance |
| `CyclingFTP()` | Functional threshold power |

### Body composition

| Method | Description |
|---|---|
| `BodyComposition(start, end)` | Weight and body composition over a range |
| `WeighIns(start, end)` | All weigh-in entries in a range |
| `DailyWeighIns(date)` | Weigh-ins for a single day |
| `AddWeighIn(kg, unitKey, timestamp)` | Log a weigh-in |
| `DeleteWeighIn(date, weightPK)` | Delete a weigh-in entry |

### Goals & achievements

| Method | Description |
|---|---|
| `Goals(status, start, limit)` | Goals filtered by status |
| `EarnedBadges()` | Earned badges |
| `AvailableBadges()` | Available badges |
| `AdHocChallenges(start, limit)` | Ad-hoc challenges |
| `BadgeChallenges(start, limit)` | Badge challenges |
| `AvailableBadgeChallenges(start, limit)` | Joinable badge challenges |
| `InProgressVirtualChallenges(start, limit)` | Active virtual challenges |

### Devices & gear

| Method | Description |
|---|---|
| `Devices()` | Registered devices |
| `DeviceSettings(deviceID)` | Settings for a specific device |
| `LastUsedDevice()` | Most recently synced device |
| `PrimaryTrainingDevice()` | Primary training device |
| `DeviceSolarData(deviceID, start, end)` | Solar charging data |
| `Gear(userProfileNumber)` | Gear items |
| `GearStats(gearUUID)` | Usage stats for a gear item |
| `GearActivities(gearUUID, start, limit)` | Activities using a gear item |
| `GearDefaults(userProfileNumber)` | Default gear assignments |

### Profile

| Method | Description |
|---|---|
| `UserProfile()` | Display name, location, join date |
| `UserProfileSettings()` | Account settings |
| `DisplayName()` | Cached display name (set on login) |

### Women's health

| Method | Description |
|---|---|
| `MenstrualData(date)` | Menstrual cycle data for a day |
| `MenstrualCalendar(start, end)` | Cycle data over a date range |
| `PregnancySummary()` | Pregnancy snapshot |

### Nutrition

| Method | Description |
|---|---|
| `NutritionFoodLog(date)` | Food log for a day |
| `NutritionMeals(date)` | Meal breakdown |
| `NutritionSettings(date)` | Nutrition goal settings |

### Golf

| Method | Description |
|---|---|
| `GolfSummary(start, limit)` | Scorecard summaries |
| `GolfScorecard(scorecardID)` | Full scorecard detail |
| `GolfShotData(scorecardID, holes)` | Shot-by-shot data for selected holes |

## Error handling

```go
import "errors"

acts, err := client.Activities(10)
switch {
case errors.Is(err, garminconnect.ErrUnauthorized):
    // token expired — call Login again
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
GARMIN_EMAIL=you@example.com GARMIN_PASSWORD=secret bash tools/record_cassettes.sh
# Only record missing cassettes:
bash tools/record_cassettes.sh --missing
```

The script logs in once, records one cassette per test, then runs `tools/sanitize_cassettes.py` to strip PII before committing.
