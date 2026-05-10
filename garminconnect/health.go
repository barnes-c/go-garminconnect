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
	"time"
)

// MenstrualData returns menstrual cycle data for the given date.
func (c *Client) MenstrualData(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/women-health-service/menstrualcycle/dayview/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// MenstrualCalendar returns menstrual cycle data between start and end dates.
func (c *Client) MenstrualCalendar(start, end time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/women-health-service/menstrualcycle/calendar/%s/%s", date(start), date(end)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PregnancySummary returns the current pregnancy snapshot for the authenticated user.
func (c *Client) PregnancySummary() (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get("/women-health-service/pregnancy/snapshot", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
