package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Device holds registration information for a Garmin device.
type Device struct {
	DeviceID               int64  `json:"deviceId"`
	ProductDisplayName     string `json:"productDisplayName"`
	DisplayName            string `json:"displayName"`
	UnitID                 int64  `json:"unitId"`
	DeviceStatus           string `json:"deviceStatus"`
	ActiveForGoals         bool   `json:"activeForGoals"`
	ImageURL               string `json:"imageUrl"`
	SerialNumber           string `json:"serialNumber"`
	CurrentFirmwareVersion string `json:"currentFirmwareVersion"`
	RegisteredDate         int64  `json:"registeredDate"` // epoch milliseconds
	RegistrationDate       struct {
		LocalRegistrationAppDate string `json:"localRegistrationAppDate"`
	} `json:"registrationDate"`
}

// Devices returns all Garmin devices registered to the authenticated user.
func (c *Client) Devices(ctx context.Context) ([]Device, error) {
	var out []Device
	if err := c.get(ctx, "/device-service/deviceregistration/devices", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceAlarm is a single alarm configured on a device.
type DeviceAlarm struct {
	AlarmMode    string   `json:"alarmMode"`
	AlarmTime    int      `json:"alarmTime"`
	AlarmDays    []string `json:"alarmDays"`
	AlarmSound   string   `json:"alarmSound"`
	AlarmID      int64    `json:"alarmId"`
	ChangeState  string   `json:"changeState"`
	Backlight    string   `json:"backlight"`
	Enabled      *bool    `json:"enabled"`
	AlarmMessage *string  `json:"alarmMessage"`
}

// DeviceLanguage is a supported display language on a device.
type DeviceLanguage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DeviceSettings holds configuration for a device. Garmin returns ~150 fields,
// most of them device-specific and often null; the common settings are typed
// here and the rest are dropped.
type DeviceSettings struct {
	DeviceID                       int64            `json:"deviceId"`
	TimeFormat                     string           `json:"timeFormat"`
	DateFormat                     string           `json:"dateFormat"`
	MeasurementUnits               string           `json:"measurementUnits"`
	AllUnits                       string           `json:"allUnits"`
	Alarms                         []DeviceAlarm    `json:"alarms"`
	MultipleAlarmEnabled           bool             `json:"multipleAlarmEnabled"`
	SupportedLanguages             []DeviceLanguage `json:"supportedLanguages"`
	Language                       int              `json:"language"`
	SupportedAudioPromptDialects   []string         `json:"supportedAudioPromptDialects"`
	AutoSyncStepsBeforeSync        int              `json:"autoSyncStepsBeforeSync"`
	AutoSyncMinutesBeforeSync      int              `json:"autoSyncMinutesBeforeSync"`
	DNDEnabled                     bool             `json:"dndEnabled"`
	StartOfWeek                    string           `json:"startOfWeek"`
	IntensityMinutesCalcMethod     string           `json:"intensityMinutesCalcMethod"`
	ModerateIntensityMinutesHrZone int              `json:"moderateIntensityMinutesHrZone"`
	VigorousIntensityMinutesHrZone int              `json:"vigorousIntensityMinutesHrZone"`
}

// DeviceSettings returns configuration settings for the given device ID.
func (c *Client) DeviceSettings(ctx context.Context, deviceID int64) (*DeviceSettings, error) {
	var out DeviceSettings
	if err := c.get(ctx, fmt.Sprintf("/device-service/deviceservice/device-info/settings/%d", deviceID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// LastUsedDevice holds details about the most recently synced device.
type LastUsedDevice struct {
	UserDeviceID                 int64  `json:"userDeviceId"`
	UserProfileNumber            int64  `json:"userProfileNumber"`
	ApplicationNumber            int    `json:"applicationNumber"`
	LastUsedDeviceApplicationKey string `json:"lastUsedDeviceApplicationKey"`
	LastUsedDeviceName           string `json:"lastUsedDeviceName"`
	LastUsedDeviceUploadTime     int64  `json:"lastUsedDeviceUploadTime"`
	ImageURL                     string `json:"imageUrl"`
	Released                     bool   `json:"released"`
}

// LastUsedDevice returns information about the most recently synced device.
func (c *Client) LastUsedDevice(ctx context.Context) (*LastUsedDevice, error) {
	var out LastUsedDevice
	if err := c.get(ctx, "/device-service/deviceservice/mylastused", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeviceRef is a minimal reference to a device by ID.
type DeviceRef struct {
	DeviceID int64 `json:"deviceId"`
}

// DeviceWeight describes a device's training-priority weighting.
type DeviceWeight struct {
	DisplayName            string `json:"displayName"`
	DeviceID               int64  `json:"deviceId"`
	ImageURL               string `json:"imageUrl"`
	Weight                 int    `json:"weight"`
	PrimaryTrainingCapable bool   `json:"primaryTrainingCapable"`
	LHABackupCapable       bool   `json:"lhaBackupCapable"`
	PrimaryWearableDevice  bool   `json:"primaryWearableDevice"`
}

// DeviceWeightGroup is a set of weighted devices with a count. Only one of the
// count fields is populated, depending on the group.
type DeviceWeightGroup struct {
	DeviceWeights              []DeviceWeight `json:"deviceWeights"`
	WearableDeviceCount        int            `json:"wearableDeviceCount,omitempty"`
	PrimaryTrainingDeviceCount int            `json:"primaryTrainingDeviceCount,omitempty"`
}

// RegisteredDevice holds identity and firmware details for a registered device.
// Garmin also returns ~250 capability flags per device, which are omitted here.
type RegisteredDevice struct {
	DisplayName            string `json:"displayName"`
	ProductDisplayName     string `json:"productDisplayName"`
	DeviceID               int64  `json:"deviceId"`
	UnitID                 int64  `json:"unitId"`
	SerialNumber           string `json:"serialNumber"`
	PartNumber             string `json:"partNumber"`
	ProductSku             string `json:"productSku"`
	DeviceStatus           string `json:"deviceStatus"`
	CurrentFirmwareVersion string `json:"currentFirmwareVersion"`
	RegisteredDate         int64  `json:"registeredDate"`
	Primary                bool   `json:"primary"`
	ImageURL               string `json:"imageUrl"`
}

// PrimaryTrainingDevice holds the user's primary training device configuration.
type PrimaryTrainingDevice struct {
	PrimaryDevice             DeviceRef          `json:"PrimaryTrainingDevice"`
	WearableDevices           DeviceWeightGroup  `json:"WearableDevices"`
	TrainingStatusOnlyDevices DeviceWeightGroup  `json:"TrainingStatusOnlyDevices"`
	PrimaryTrainingDevices    DeviceWeightGroup  `json:"PrimaryTrainingDevices"`
	RegisteredDevices         []RegisteredDevice `json:"RegisteredDevices"`
}

// PrimaryTrainingDevice returns information about the user's primary training device.
func (c *Client) PrimaryTrainingDevice(ctx context.Context) (*PrimaryTrainingDevice, error) {
	var out PrimaryTrainingDevice
	if err := c.get(ctx, "/web-gateway/device-info/primary-training-device", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SolarDailyData is one day of solar charging data. Field names are inferred;
// verify against a solar-capable device's response, as the list is empty for
// devices without a solar panel.
type SolarDailyData struct {
	CalendarDate     string  `json:"calendarDate"`
	SolarUtilization float64 `json:"solarUtilization"`
}

// DeviceSolarData holds solar charging data for a device over a date range.
type DeviceSolarData struct {
	SolarDailyDataDTOs []SolarDailyData `json:"solarDailyDataDTOs"`
}

// DeviceSolarData returns solar charging data for the given device between start
// and end dates. The endpoint wraps the data in {"deviceSolarInput": ...}.
func (c *Client) DeviceSolarData(ctx context.Context, deviceID int64, start, end string) (*DeviceSolarData, error) {
	params := url.Values{"singleDayView": {strconv.FormatBool(start == end)}}
	var wrapper struct {
		DeviceSolarInput DeviceSolarData `json:"deviceSolarInput"`
	}
	if err := c.get(ctx, fmt.Sprintf("/web-gateway/solar/%d/%s/%s", deviceID, start, end), params, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.DeviceSolarInput, nil
}
