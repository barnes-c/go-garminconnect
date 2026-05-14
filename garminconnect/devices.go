package garminconnect

import (
	"encoding/json"
	"fmt"
)

// Device holds registration information for a Garmin device.
type Device struct {
	DeviceID           int64  `json:"deviceId"`
	ProductDisplayName string `json:"productDisplayName"`
	DisplayName        string `json:"displayName"`
	UnitID             int64  `json:"unitId"`
	DeviceStatus       string `json:"deviceStatus"`
	ActiveForGoals     bool   `json:"activeForGoals"`
	ImageURL           string `json:"imageUrl"`
	RegistrationDate   struct {
		LocalRegistrationAppDate string `json:"localRegistrationAppDate"`
	} `json:"registrationDate"`
}

// Devices returns all Garmin devices registered to the authenticated user.
func (c *Client) Devices() ([]Device, error) {
	var out []Device
	if err := c.get("/device-service/deviceregistration/devices", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceSettings returns configuration settings for the given device ID.
func (c *Client) DeviceSettings(deviceID int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/device-service/deviceservice/device-info/settings/%d", deviceID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// LastUsedDevice returns information about the most recently synced device.
func (c *Client) LastUsedDevice() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/device-service/deviceservice/mylastused", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PrimaryTrainingDevice returns information about the user's primary training device.
func (c *Client) PrimaryTrainingDevice() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/web-gateway/device-info/primary-training-device", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceSolarData returns solar charging data for the given device between start and end dates.
func (c *Client) DeviceSolarData(deviceID int64, start, end string) ([]map[string]json.RawMessage, error) {
	var out []map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/web-gateway/solar/%d/%s/%s", deviceID, start, end), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
