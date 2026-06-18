package garminconnect

import (
	"context"
	"encoding/json"
	"net/url"
)

// ActiveTrainingPlan returns the athlete's currently-active annual training
// plan (ATP). Returns a nil map if there is no active plan (HTTP 204).
func (c *Client) ActiveTrainingPlan(ctx context.Context) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, "/atp-api/atp/athlete/active", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CompletedTrainingPlans returns the athlete's completed annual training plans (ATP).
func (c *Client) CompletedTrainingPlans(ctx context.Context) ([]map[string]json.RawMessage, error) {
	var out []map[string]json.RawMessage
	if err := c.get(ctx, "/atp-api/atp/athlete/completed", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TrainingPlanType describes one annual training plan (ATP) type in the catalog.
type TrainingPlanType struct {
	WorkoutPlanTypeID  int      `json:"workoutPlanTypeId"`
	Name               string   `json:"name"`
	TrainingType       string   `json:"trainingType"`
	TrainingSubtype    string   `json:"trainingSubtype"`
	TrainingLevel      []string `json:"trainingLevel"`
	TrainingVersion    *string  `json:"trainingVersion"`
	PlanDistanceMeters int      `json:"planDistanceMeters"`
	WeeklyWorkoutsMin  int      `json:"weeklyWorkoutsMin"`
	WeeklyWorkoutsMax  int      `json:"weeklyWorkoutsMax"`
	PlanDurationMin    int      `json:"planDurationMin"`
	PlanDurationMax    int      `json:"planDurationMax"`
}

// TrainingPlanTypes returns the catalog of annual training plan (ATP) types.
func (c *Client) TrainingPlanTypes(ctx context.Context) ([]TrainingPlanType, error) {
	params := url.Values{"lang": {"en-US"}}
	var out []TrainingPlanType
	if err := c.get(ctx, "/atp-api/atp/types/", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}
