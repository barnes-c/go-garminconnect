package garminconnect

import (
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
func (c *Client) Workouts(start, limit int) ([]Workout, error) {
	params := url.Values{
		"start": {fmt.Sprintf("%d", start)},
		"limit": {fmt.Sprintf("%d", limit)},
	}
	var out []Workout
	if err := c.get("/workout-service/workouts", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Workout returns full details for a single saved workout.
func (c *Client) Workout(id int64) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.get(fmt.Sprintf("/workout-service/workout/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteWorkout permanently deletes a saved workout.
func (c *Client) DeleteWorkout(id int64) error {
	return c.del(fmt.Sprintf("/workout-service/workout/%d", id))
}

// ScheduledWorkouts returns the calendar for the given year and month.
func (c *Client) ScheduledWorkouts(year, month int) (json.RawMessage, error) {
	var out json.RawMessage
	if err := c.get(fmt.Sprintf("/calendar-service/year/%d/month/%d", year, month-1), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ScheduleWorkout assigns a saved workout to a calendar date (YYYY-MM-DD).
func (c *Client) ScheduleWorkout(workoutID int64, date string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	body := map[string]any{"date": date}
	if err := c.post(fmt.Sprintf("/calendar-service/schedule/workout/%d", workoutID), body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnscheduleWorkout removes a workout from the calendar.
func (c *Client) UnscheduleWorkout(scheduledWorkoutID int64) error {
	return c.del(fmt.Sprintf("/calendar-service/schedule/workout/%d", scheduledWorkoutID))
}

// DownloadWorkout returns the FIT file for a saved workout.
func (c *Client) DownloadWorkout(id int64) ([]byte, error) {
	return c.getBytes(fmt.Sprintf("/download-service/files/workout/%d", id), nil)
}

// UploadWorkout uploads a workout FIT file and returns the server response.
func (c *Client) UploadWorkout(data []byte, filename string) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.upload("/upload-service/upload", data, filename, &out); err != nil {
		return nil, err
	}
	return out, nil
}
