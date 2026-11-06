package garminconnect

import (
	"context"
	"encoding/json"
)

// UserProfile holds public profile information for a Garmin Connect user.
type UserProfile struct {
	DisplayName          string `json:"displayName"`
	FullName             string `json:"fullName"`
	UserProfilePK        int    `json:"profileId"`
	ProfileImageURL      string `json:"profileImageUrlMedium"`
	ProfileImageURLLarge string `json:"profileImageUrlLarge"`
	ProfileImageURLSmall string `json:"profileImageUrlSmall"`
	Location             string `json:"location"`
	Biography            string `json:"bio"`
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

// UpdateUserSettings changes the authenticated user's settings. Garmin merges
// the patch into the current settings, so only the keys to change need to be
// sent, using the user-settings JSON layout, e.g.
//
//	map[string]any{"userData": map[string]any{"timeFormat": "time_twenty_four_hr"}}
func (c *Client) UpdateUserSettings(ctx context.Context, patch map[string]any) error {
	return c.put(ctx, "/userprofile-service/userprofile/user-settings", patch, nil)
}

// PersonalInformation holds the authenticated user's account details and
// biometric profile.
type PersonalInformation struct {
	UserInfo         PersonalUserInfo `json:"userInfo"`
	BiometricProfile BiometricProfile `json:"biometricProfile"`
	TimeZone         string           `json:"timeZone"`
	Locale           string           `json:"locale"`
	Gender           *string          `json:"gender"`
	BirthDate        string           `json:"birthDate"`
}

// PersonalUserInfo holds the account details within PersonalInformation.
type PersonalUserInfo struct {
	BirthDate   string  `json:"birthDate"`
	GenderType  *string `json:"genderType"`
	Email       string  `json:"email"`
	Locale      string  `json:"locale"`
	TimeZone    string  `json:"timeZone"`
	Age         int     `json:"age"`
	CountryCode string  `json:"countryCode"`
}

// BiometricProfile holds the user's physiological settings. Fields are null
// until measured or entered.
type BiometricProfile struct {
	UserID                    int64    `json:"userId"`
	Height                    *float64 `json:"height"`
	Weight                    *float64 `json:"weight"`
	VO2Max                    *float64 `json:"vo2Max"`
	VO2MaxCycling             *float64 `json:"vo2MaxCycling"`
	VO2MaxRowing              *float64 `json:"vo2MaxRowing"`
	LactateThresholdHeartRate *float64 `json:"lactateThresholdHeartRate"`
	ActivityClass             *float64 `json:"activityClass"`
	FunctionalThresholdPower  *float64 `json:"functionalThresholdPower"`
	CriticalSwimSpeed         *float64 `json:"criticalSwimSpeed"`
}

// PersonalInformation returns the authenticated user's account details and
// biometric profile.
func (c *Client) PersonalInformation(ctx context.Context) (*PersonalInformation, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out PersonalInformation
	if err := c.get(ctx, "/userprofile-service/userprofile/personal-information/"+name, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserProfileBase holds the authenticated user's core account record.
type UserProfileBase struct {
	UserProfilePK              int64   `json:"userProfilePk"`
	UserName                   string  `json:"userName"`
	FirstName                  string  `json:"firstName"`
	LastName                   string  `json:"lastName"`
	BirthDate                  string  `json:"birthDate"`
	Gender                     *string `json:"gender"`
	EmailAddress               string  `json:"emailAddress"`
	CreateDate                 string  `json:"createDate"`
	MeasurementSystemPK        int     `json:"measurementSystemPk"`
	GlucoseMeasurementUnitID   int     `json:"glucoseMeasurementUnitId"`
	HydrationMeasurementUnitID int     `json:"hydrationMeasurementUnitId"`
	TimeZonePK                 int     `json:"timeZonePk"`
	DecimalFormat              int     `json:"decimalFormat"`
	TimeFormat                 int     `json:"timeFormat"`
	FormatLocalePK             int     `json:"formatLocalePk"`
	DayOfWeekPK                int     `json:"dayOfWeekPk"`
	GarminGlobalID             string  `json:"garminGlobalId"`
	DisplayName                string  `json:"displayName"`
	TOCAcceptedDate            string  `json:"tocAcceptedDate"`
	GarminGUID                 string  `json:"garminGUID"`
	CountryCode                string  `json:"countryCode"`
	CountryCodeVerified        bool    `json:"countryCodeVerified"`
	GolfDistanceUnitID         int     `json:"golfDistanceUnitId"`
	PhoneNumber                string  `json:"phoneNumber"`
}

// UserProfileBase returns the authenticated user's core account record.
func (c *Client) UserProfileBase(ctx context.Context) (*UserProfileBase, error) {
	var out UserProfileBase
	if err := c.get(ctx, "/userprofile-service/userprofile/userProfileBase", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PrivacySettings holds the authenticated user's sharing rules per category.
type PrivacySettings struct {
	Activity PrivacyCategory `json:"activityPrivacySettings"`
	Wellness PrivacyCategory `json:"wellnessPrivacySettings"`
	Profile  PrivacyCategory `json:"profilePrivacySettings"`
}

// PrivacyCategory holds the rule for one privacy category and its
// subcategories (e.g. "running" under activities).
type PrivacyCategory struct {
	CategoryDetails struct {
		CategoryAccessRule   string  `json:"categoryAccessRule"`
		InheritedFromDefault bool    `json:"inheritedFromDefault"`
		UpdatedTimestampGMT  *string `json:"updatedTimestampGMT"`
	} `json:"categoryDetails"`
	Subcategories []PrivacySubcategory `json:"subcategoryPrivacySettings"`
}

// PrivacySubcategory holds the access rule for one privacy subcategory.
type PrivacySubcategory struct {
	Type                            string `json:"privacySubcategoryType"`
	AssociatedParentActivityTypeIDs []int  `json:"associatedParentActivityTypeIds"`
	Details                         struct {
		AccessRule            string  `json:"accessRule"`
		UpdatedTimestampGMT   *string `json:"updatedTimestampGMT"`
		InheritedFromCategory bool    `json:"inheritedFromCategory"`
		InheritedFromDefault  bool    `json:"inheritedFromDefault"`
	} `json:"subcategoryDetails"`
}

// PrivacySettings returns the authenticated user's privacy settings.
func (c *Client) PrivacySettings(ctx context.Context) (*PrivacySettings, error) {
	var out PrivacySettings
	if err := c.get(ctx, "/privacysettings-service/privacy-settings", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UsageIndicators summarizes which devices and features the authenticated
// user has, which Garmin Connect uses to decide what to show.
type UsageIndicators struct {
	DeviceBasedIndicators           DeviceIndicators `json:"deviceBasedIndicators"`
	HistoricalDeviceBasedIndicators DeviceIndicators `json:"historicalDeviceBasedIndicators"`
	PerformanceBasedIndicators      struct {
		WellnessIndicators struct {
			FloorsAscended       *float64 `json:"floorsAscended"`
			FloorsDescended      *float64 `json:"floorsDescended"`
			TotalSteps           *int     `json:"totalSteps"`
			LastDailySummaryDate *string  `json:"lastDailySummaryDate"`
		} `json:"wellnessIndicators"`
		// ActivityIndicators is keyed by activity type, e.g. "running".
		ActivityIndicators map[string]struct {
			LastActivityID   *int64  `json:"lastActivityId"`
			LastActivityDate *string `json:"lastActivityDate"`
		} `json:"activityIndicators"`
	} `json:"performanceBasedIndicators"`
	ProfileBasedIndicators struct {
		IsPulseOxAllowed bool `json:"isPulseOxAllowed"`
	} `json:"profileBasedIndicators"`
}

// DeviceIndicators counts the user's devices by kind.
type DeviceIndicators struct {
	GolfDeviceCount                  int  `json:"golfDeviceCount"`
	OutdoorDeviceCount               int  `json:"outdoorDeviceCount"`
	HasAllDayHeartRateCapableDevice  bool `json:"hasAllDayHeartRateCapableDevice"`
	IsTacxOnlyUser                   bool `json:"isTacxOnlyUser"`
	DedicatedBikeComputerDeviceCount int  `json:"dedicatedBikeComputerDeviceCount"`
	WellnessDeviceCount              int  `json:"wellnessDeviceCount"`
	WeightScaleDeviceCount           int  `json:"weightScaleDeviceCount"`
	FitnessDeviceCount               int  `json:"fitnessDeviceCount"`
}

// UsageIndicators returns the authenticated user's device and feature usage
// summary.
func (c *Client) UsageIndicators(ctx context.Context) (*UsageIndicators, error) {
	var out UsageIndicators
	if err := c.get(ctx, "/web-gateway/snapshot/usageIndicators", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserFocus is the training focus the user picked in Garmin Connect.
type UserFocus struct {
	UserFocus       string  `json:"userFocus"`
	CreateTimestamp *string `json:"createTimestamp"`
}

// UserFocus returns the authenticated user's training focus.
func (c *Client) UserFocus(ctx context.Context) (*UserFocus, error) {
	var out UserFocus
	if err := c.get(ctx, "/userfocus-service/focus", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserLocation is the authenticated user's account country.
type UserLocation struct {
	UserProfileID                int64  `json:"userProfileId"`
	CountryCode                  string `json:"countryCode"`
	CountryCodeVerified          bool   `json:"countryCodeVerified"`
	CountryCodeVerifiedTimestamp string `json:"countryCodeVerifiedTimestamp"`
}

// UserLocation returns the authenticated user's account country.
func (c *Client) UserLocation(ctx context.Context) (*UserLocation, error) {
	var out UserLocation
	if err := c.get(ctx, "/userprofile-service/userprofile/location", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// FollowingsCount returns how many users the authenticated user follows.
func (c *Client) FollowingsCount(ctx context.Context) (int, error) {
	var out int
	if err := c.get(ctx, "/follower-service/follow/followings/count", nil, &out); err != nil {
		return 0, err
	}
	return out, nil
}
