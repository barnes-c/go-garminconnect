package garminconnect

import (
	"context"
	"encoding/json"
)

// UserProfile holds public profile information for a Garmin Connect user.
type UserProfile struct {
	DisplayName          string `json:"displayName"`
	FullName             string `json:"fullName"`
	UserProfilePK        int    `json:"userProfilePK"`
	ProfileImageURL      string `json:"profileImageUrl"`
	ProfileImageURLLarge string `json:"profileImageUrlLarge"`
	ProfileImageURLSmall string `json:"profileImageUrlSmall"`
	Location             string `json:"location"`
	Biography            string `json:"biography"`
}

// UserProfileSettings holds account and display settings for the authenticated
// user. Format fields reuse MeasurementFormat and FirstDayOfWeek.
type UserProfileSettings struct {
	DisplayName               string            `json:"displayName"`
	PreferredLocale           string            `json:"preferredLocale"`
	MeasurementSystem         string            `json:"measurementSystem"`
	FirstDayOfWeek            FirstDayOfWeek    `json:"firstDayOfWeek"`
	NumberFormat              string            `json:"numberFormat"`
	TimeFormat                MeasurementFormat `json:"timeFormat"`
	DateFormat                MeasurementFormat `json:"dateFormat"`
	PowerFormat               MeasurementFormat `json:"powerFormat"`
	HeartRateFormat           MeasurementFormat `json:"heartRateFormat"`
	TimeZone                  string            `json:"timeZone"`
	HydrationMeasurementUnit  string            `json:"hydrationMeasurementUnit"`
	HydrationContainers       []json.RawMessage `json:"hydrationContainers"`
	GolfDistanceUnit          string            `json:"golfDistanceUnit"`
	GolfElevationUnit         *string           `json:"golfElevationUnit"`
	GolfSpeedUnit             *string           `json:"golfSpeedUnit"`
	AvailableTrainingDays     []string          `json:"availableTrainingDays"`
	PreferredLongTrainingDays []string          `json:"preferredLongTrainingDays"`
}

// UserProfileSettings returns account and display settings for the authenticated user.
func (c *Client) UserProfileSettings(ctx context.Context) (*UserProfileSettings, error) {
	var out UserProfileSettings
	if err := c.get(ctx, "/userprofile-service/userprofile/settings", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserSettings holds the authenticated user's account-level preferences from
// the user-settings endpoint. Fields that Garmin may report as null are
// pointers. Numeric field types are inferred from field names; verify against a
// live response if a field matters.
type UserSettings struct {
	ID               int64         `json:"id"`
	UserData         UserData      `json:"userData"`
	UserSleep        UserSleep     `json:"userSleep"`
	ConnectDate      *string       `json:"connectDate"`
	SourceType       *string       `json:"sourceType"`
	UserSleepWindows []SleepWindow `json:"userSleepWindows"`
}

// UserData holds measurement preferences and biometric settings within UserSettings.
type UserData struct {
	Gender                           string            `json:"gender"`
	Weight                           float64           `json:"weight"`
	Height                           float64           `json:"height"`
	TimeFormat                       string            `json:"timeFormat"`
	BirthDate                        string            `json:"birthDate"`
	MeasurementSystem                string            `json:"measurementSystem"`
	ActivityLevel                    *int              `json:"activityLevel"`
	Handedness                       string            `json:"handedness"`
	PowerFormat                      MeasurementFormat `json:"powerFormat"`
	HeartRateFormat                  MeasurementFormat `json:"heartRateFormat"`
	FirstDayOfWeek                   FirstDayOfWeek    `json:"firstDayOfWeek"`
	VO2MaxRunning                    *float64          `json:"vo2MaxRunning"`
	VO2MaxCycling                    *float64          `json:"vo2MaxCycling"`
	LactateThresholdSpeed            *float64          `json:"lactateThresholdSpeed"`
	LactateThresholdHeartRate        *int              `json:"lactateThresholdHeartRate"`
	DiveNumber                       *int              `json:"diveNumber"`
	IntensityMinutesCalcMethod       string            `json:"intensityMinutesCalcMethod"`
	ModerateIntensityMinutesHrZone   int               `json:"moderateIntensityMinutesHrZone"`
	VigorousIntensityMinutesHrZone   int               `json:"vigorousIntensityMinutesHrZone"`
	HydrationMeasurementUnit         string            `json:"hydrationMeasurementUnit"`
	HydrationContainers              []json.RawMessage `json:"hydrationContainers"`
	HydrationAutoGoalEnabled         bool              `json:"hydrationAutoGoalEnabled"`
	FirstbeatMaxStressScore          *float64          `json:"firstbeatMaxStressScore"`
	FirstbeatCyclingLtTimestamp      *int64            `json:"firstbeatCyclingLtTimestamp"`
	FirstbeatRunningLtTimestamp      *int64            `json:"firstbeatRunningLtTimestamp"`
	ThresholdHeartRateAutoDetected   bool              `json:"thresholdHeartRateAutoDetected"`
	FTPAutoDetected                  bool              `json:"ftpAutoDetected"`
	TrainingStatusPausedDate         *string           `json:"trainingStatusPausedDate"`
	WeatherLocation                  WeatherLocation   `json:"weatherLocation"`
	GolfDistanceUnit                 string            `json:"golfDistanceUnit"`
	GolfElevationUnit                *string           `json:"golfElevationUnit"`
	GolfSpeedUnit                    *string           `json:"golfSpeedUnit"`
	ExternalBottomTime               *float64          `json:"externalBottomTime"`
	AvailableTrainingDays            []string          `json:"availableTrainingDays"`
	PreferredLongTrainingDays        []string          `json:"preferredLongTrainingDays"`
	VirtualCaddieDataSource          *string           `json:"virtualCaddieDataSource"`
	NumberDivesAutomatically         *bool             `json:"numberDivesAutomatically"`
	FirstbeatRowingLtTimestamp       *int64            `json:"firstbeatRowingLtTimestamp"`
	LactateThresholdRowingPace       *float64          `json:"lactateThresholdRowingPace"`
	LactateThresholdHeartRateRowing  *int              `json:"lactateThresholdHeartRateRowing"`
	LactateThresholdHeartRateCycling *int              `json:"lactateThresholdHeartRateCycling"`
}

// MeasurementFormat describes how a metric (e.g. power or heart rate) is
// formatted for display.
type MeasurementFormat struct {
	FormatID      int     `json:"formatId"`
	FormatKey     string  `json:"formatKey"`
	MinFraction   int     `json:"minFraction"`
	MaxFraction   int     `json:"maxFraction"`
	GroupingUsed  bool    `json:"groupingUsed"`
	DisplayFormat *string `json:"displayFormat"`
}

// FirstDayOfWeek identifies the user's configured first day of the week.
type FirstDayOfWeek struct {
	DayID              int    `json:"dayId"`
	DayName            string `json:"dayName"`
	SortOrder          int    `json:"sortOrder"`
	IsPossibleFirstDay bool   `json:"isPossibleFirstDay"`
}

// WeatherLocation holds the user's fixed weather location, if one is set.
type WeatherLocation struct {
	UseFixedLocation *bool    `json:"useFixedLocation"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	LocationName     *string  `json:"locationName"`
	ISOCountryCode   *string  `json:"isoCountryCode"`
	PostalCode       *string  `json:"postalCode"`
}

// UserSleep holds the user's default sleep and wake times.
type UserSleep struct {
	SleepTime        int  `json:"sleepTime"`
	DefaultSleepTime bool `json:"defaultSleepTime"`
	WakeTime         int  `json:"wakeTime"`
	DefaultWakeTime  bool `json:"defaultWakeTime"`
}

// SleepWindow is one configured sleep window, in seconds from midnight.
type SleepWindow struct {
	SleepWindowFrequency              string `json:"sleepWindowFrequency"`
	StartSleepTimeSecondsFromMidnight int    `json:"startSleepTimeSecondsFromMidnight"`
	EndSleepTimeSecondsFromMidnight   int    `json:"endSleepTimeSecondsFromMidnight"`
}

// UserProfile returns detailed profile information for the authenticated user.
func (c *Client) UserProfile(ctx context.Context) (*UserProfile, error) {
	var out UserProfile
	if err := c.get(ctx, "/userprofile-service/socialProfile", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserSettings returns the authenticated user's account-level preferences.
func (c *Client) UserSettings(ctx context.Context) (*UserSettings, error) {
	var out UserSettings
	if err := c.get(ctx, "/userprofile-service/userprofile/user-settings", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnitSystem returns the authenticated user's measurement system,
// e.g. "metric" or "statute_us".
func (c *Client) UnitSystem(ctx context.Context) (string, error) {
	s, err := c.UserSettings(ctx)
	if err != nil {
		return "", err
	}
	return s.UserData.MeasurementSystem, nil
}
