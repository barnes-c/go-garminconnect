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

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Gear represents a piece of equipment tracked in Garmin Connect.
type Gear struct {
	GearPK          int64  `json:"gearPk"`
	UUID            string `json:"uuid"`
	GearTypeName    string `json:"gearTypeName"`
	DisplayName     string `json:"displayName"`
	CustomMakeModel string `json:"customMakeModel"`
	MaxMeters       int    `json:"maxMeters"`
	NotifiedMeters  int    `json:"notifiedAtMeters"`
	DateBegin       string `json:"dateBegin"`
	DateEnd         string `json:"dateEnd"`
	GearStatusName  string `json:"gearStatusName"`
}

// Gear returns all gear registered to the given user profile number.
func (c *Client) Gear(userProfileNumber int) ([]Gear, error) {
	params := url.Values{"userProfilePk": {fmt.Sprintf("%d", userProfileNumber)}}
	var out []Gear
	if err := c.get("/gear-service/gear/filterGear", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearStats returns usage statistics for the given gear UUID.
func (c *Client) GearStats(gearUUID string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/gear-service/gear/stats/%s", gearUUID), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearActivities returns activities that used the given gear.
func (c *Client) GearActivities(gearUUID string, start, limit int) (map[string]json.RawMessage, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/gear-service/gear/%s/activities", gearUUID), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GearDefaults returns the default gear assigned per activity type for the user.
func (c *Client) GearDefaults(userProfileNumber int) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/gear-service/gear/user/%d/activityTypes", userProfileNumber), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
