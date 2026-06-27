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

// UserProfile returns detailed profile information for the authenticated user.
func (c *Client) UserProfile(ctx context.Context) (*UserProfile, error) {
	var out UserProfile
	if err := c.get(ctx, "/userprofile-service/socialProfile", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserProfileSettings returns account and display settings for the authenticated user.
func (c *Client) UserProfileSettings(ctx context.Context) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/userprofile-service/userprofile/settings", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnitSystem returns the authenticated user's measurement system,
// e.g. "metric" or "statute_us".
func (c *Client) UnitSystem(ctx context.Context) (string, error) {
	var out struct {
		UserData struct {
			MeasurementSystem string `json:"measurementSystem"`
		} `json:"userData"`
	}
	if err := c.get(ctx, "/userprofile-service/userprofile/user-settings", nil, &out); err != nil {
		return "", err
	}
	return out.UserData.MeasurementSystem, nil
}
