package garminconnect

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// TrainingPlan is a Garmin training plan: a catalog template from
// SearchTrainingPlans, or the user's scheduled copy of one.
type TrainingPlan struct {
	TrainingPlanID    int64  `json:"trainingPlanId"`
	Category          string `json:"trainingPlanCategory"` // e.g. "STATIC"
	Name              string `json:"name"`
	Description       string `json:"description"`
	DurationInWeeks   int    `json:"durationInWeeks"`
	AvgWeeklyWorkouts int    `json:"avgWeeklyWorkouts"`
	StartDate         string `json:"startDate"` // set once scheduled
	EndDate           string `json:"endDate"`
	TrainingType      struct {
		TypeKey string `json:"typeKey"` // e.g. "running"
	} `json:"trainingType"`
	TrainingLevel struct {
		LevelKey string `json:"levelKey"`
	} `json:"trainingLevel"`
	TrainingStatus struct {
		StatusKey string `json:"statusKey"` // e.g. "Scheduled"
	} `json:"trainingStatus"`
}

// TrainingPlanTask is one day's entry in a plan: a workout, a note, or both.
type TrainingPlanTask struct {
	TrainingPlanID int64 `json:"trainingPlanId"`
	WeekID         int   `json:"weekId"`      // 1-based week of the plan
	DayOfWeekID    int   `json:"dayOfWeekId"` // day within the week
	WorkoutOrder   int   `json:"workoutOrder"`
	Workout        *struct {
		WorkoutID     int64  `json:"workoutId"`
		WorkoutName   string `json:"workoutName"`
		ScheduledDate string `json:"scheduledDate"`
		RestDay       bool   `json:"restDay"`
	} `json:"taskWorkout"`
	Note *struct {
		NoteID       int64  `json:"noteId"`
		Note         string `json:"note"`
		CalendarDate string `json:"calendarDate"`
	} `json:"taskNote"`
}

type trainingPlanList struct {
	TrainingPlanList []TrainingPlan `json:"trainingPlanList"`
}

// TrainingPlans returns the user's training plans.
func (c *Client) TrainingPlans(ctx context.Context) ([]TrainingPlan, error) {
	var out trainingPlanList
	if err := c.get(ctx, "/trainingplan-service/trainingplan/plans", nil, &out); err != nil {
		return nil, err
	}
	return out.TrainingPlanList, nil
}

// ActiveTrainingPlans returns the user's plans that are active on day.
func (c *Client) ActiveTrainingPlans(ctx context.Context, day time.Time) ([]TrainingPlan, error) {
	var out trainingPlanList
	params := url.Values{"startDate": {date(day)}}
	if err := c.get(ctx, "/trainingplan-service/trainingplan/active-plans", params, &out); err != nil {
		return nil, err
	}
	return out.TrainingPlanList, nil
}

// SearchTrainingPlans returns catalog plans that can be scheduled, starting
// at the 1-based position start.
func (c *Client) SearchTrainingPlans(ctx context.Context, start, limit int) ([]TrainingPlan, error) {
	params := url.Values{"start": {fmt.Sprint(start)}, "limit": {fmt.Sprint(limit)}, "locale": {"en"}}
	var out trainingPlanList
	// Garmin only accepts a form-encoded body here; the filters go in the query.
	if err := c.post(ctx, "/trainingplan-service/trainingplan/search?"+params.Encode(), url.Values{}, &out); err != nil {
		return nil, err
	}
	return out.TrainingPlanList, nil
}

// TrainingPlanTasks returns every day's task of a plan.
func (c *Client) TrainingPlanTasks(ctx context.Context, planID int64) ([]TrainingPlanTask, error) {
	var out []TrainingPlanTask
	params := url.Values{"localeKey": {"en"}}
	if err := c.get(ctx, fmt.Sprintf("/trainingplan-service/trainingplan/tasks/%d", planID), params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ScheduleTrainingPlan starts a catalog plan on the given day, adding its
// workouts and notes to the calendar, and returns the user's scheduled copy,
// which has its own TrainingPlanID.
func (c *Client) ScheduleTrainingPlan(ctx context.Context, planID int64, start time.Time) (*TrainingPlan, error) {
	var out TrainingPlan
	params := url.Values{"startDate": {date(start)}}
	if err := c.get(ctx, fmt.Sprintf("/trainingplan-service/trainingplan/schedule/%d", planID), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteTrainingPlan removes a scheduled plan and the calendar items it added.
func (c *Client) DeleteTrainingPlan(ctx context.Context, scheduledPlanID int64) error {
	return c.del(ctx, fmt.Sprintf("/trainingplan-service/trainingplan/trainingplan/%d", scheduledPlanID))
}
