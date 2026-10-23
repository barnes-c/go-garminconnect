package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// NutritionFoodLog returns the food log entries for the given date.
func (c *Client) NutritionFoodLog(ctx context.Context, d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/nutrition-service/food/logs/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// NutritionMeals returns meal data for the given date.
func (c *Client) NutritionMeals(ctx context.Context, d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/nutrition-service/meals/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// NutritionSettings returns nutrition goal settings for the given date.
func (c *Client) NutritionSettings(ctx context.Context, d time.Time) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/nutrition-service/settings/%s", date(d)), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// NutritionStatus says whether the user tracks nutrition in Garmin Connect or
// MyFitnessPal.
type NutritionStatus struct {
	CurrentStatus    string `json:"currentStatus"` // e.g. "NUTRITION_NONE"
	HasUsedNutrition bool   `json:"hasUsedNutrition"`
	HasUsedMFP       bool   `json:"hasUsedMFP"`
}

// NutritionStatus returns the user's nutrition tracking status.
func (c *Client) NutritionStatus(ctx context.Context) (*NutritionStatus, error) {
	var out NutritionStatus
	if err := c.get(ctx, "/nutrition-service/user/nutritionCurrentStatus", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
