// Copyright Christopher Barnes
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if err := c.get("/user-service/userprofile/v2/information/basic", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UserProfileSettings returns account and display settings for the authenticated user.
func (c *Client) UserProfileSettings() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/userprofile-service/userprofileservice/userprofile/v2/settings", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
