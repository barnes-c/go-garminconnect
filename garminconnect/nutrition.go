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
