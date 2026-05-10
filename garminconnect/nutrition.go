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

// NutritionFoodLog returns the food log entries for the given date.
func (c *Client) NutritionFoodLog(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/nutrition-service/food/logs/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// NutritionMeals returns meal data for the given date.
func (c *Client) NutritionMeals(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/nutrition-service/meals/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// NutritionSettings returns nutrition goal settings for the given date.
func (c *Client) NutritionSettings(d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/nutrition-service/settings/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
