package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Workout is a saved workout definition.
type Workout struct {
	WorkoutID   int64  `json:"workoutId"`
	WorkoutName string `json:"workoutName"`
	Description string `json:"description"`
	SportType   struct {
		SportTypeKey string `json:"sportTypeKey"`
	} `json:"sportType"`
	CreatedDate string `json:"createdDate"`
	UpdatedDate string `json:"updatedDate"`
}

// ScheduledWorkout links a workout to a specific calendar date.
type ScheduledWorkout struct {
	ScheduledWorkoutID int64  `json:"scheduledWorkoutId"`
	WorkoutID          int64  `json:"workoutId"`
	Date               string `json:"date"`
}

// Workouts returns saved workouts with pagination.
func (c *Client) Workouts(ctx context.Context, start, limit int) ([]Workout, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out []Workout
	if err := c.get(ctx, "/workout-service/workouts", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Workout returns full details for a single saved workout.
func (c *Client) Workout(ctx context.Context, id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/workout-service/workout/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWorkout permanently deletes a saved workout.
func (c *Client) DeleteWorkout(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/workout-service/workout/%d", id))
}

// ScheduledWorkouts returns the calendar for the given year and month.
func (c *Client) ScheduledWorkouts(ctx context.Context, year, month int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.get(ctx, fmt.Sprintf("/calendar-service/year/%d/month/%d", year, month-1), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ScheduleWorkout assigns a saved workout to a calendar date (YYYY-MM-DD).
func (c *Client) ScheduleWorkout(ctx context.Context, workoutID int64, date string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{"date": date}
	if err := c.post(ctx, fmt.Sprintf("/calendar-service/schedule/workout/%d", workoutID), body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnscheduleWorkout removes a workout from the calendar.
func (c *Client) UnscheduleWorkout(ctx context.Context, scheduledWorkoutID int64) error {
	return c.del(ctx, fmt.Sprintf("/calendar-service/schedule/workout/%d", scheduledWorkoutID))
}

// DownloadWorkout returns the FIT file for a saved workout.
func (c *Client) DownloadWorkout(ctx context.Context, id int64) ([]byte, error) {
	return c.getBytes(ctx, fmt.Sprintf("/download-service/files/workout/%d", id), nil)
}

// UploadWorkout uploads a workout FIT file and returns the server response.
func (c *Client) UploadWorkout(ctx context.Context, data []byte, filename string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.upload(ctx, "/upload-service/upload", data, filename, &out); err != nil {
		return nil, err
	}
	return out, nil
}
