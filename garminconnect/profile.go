package garminconnect

import "encoding/json"

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
func (c *Client) UserProfile() (*UserProfile, error) {
	var out UserProfile
	if err := c.get("/userprofile-service/socialProfile", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserProfileSettings returns account and display settings for the authenticated user.
func (c *Client) UserProfileSettings() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/userprofile-service/userprofile/settings", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
