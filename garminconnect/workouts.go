package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"time"
)

// Workout is a saved workout definition.
type Workout struct {
	WorkoutID   int64            `json:"workoutId"`
	WorkoutName string           `json:"workoutName"`
	Description string           `json:"description"`
	SportType   WorkoutSportType `json:"sportType"`
	CreatedDate string           `json:"createdDate"`
	UpdatedDate string           `json:"updatedDate"`
}

// WorkoutSportType identifies the sport a workout targets.
type WorkoutSportType struct {
	SportTypeID  int    `json:"sportTypeId"`
	SportTypeKey string `json:"sportTypeKey"`
	DisplayOrder int    `json:"displayOrder"`
}

// ScheduledWorkout is a workout placed on a calendar date.
type ScheduledWorkout struct {
	WorkoutScheduleID    int64   `json:"workoutScheduleId"`
	Workout              Workout `json:"workout"`
	CalendarDate         string  `json:"calendarDate"`
	CreatedDate          string  `json:"createdDate"`
	OwnerID              int64   `json:"ownerId"`
	AssociatedActivityID *int64  `json:"associatedActivityId"`
	NameChanged          bool    `json:"nameChanged"`
	Protected            bool    `json:"protected"`
	ITP                  bool    `json:"itp"`
	Race                 bool    `json:"race"`
}

// WorkoutTypes lists the IDs and keys Garmin's workout JSON uses for steps,
// sports, end conditions, intensities, targets, equipment and swim strokes.
type WorkoutTypes struct {
	StepTypes []struct {
		StepTypeID   int    `json:"stepTypeId"`
		StepTypeKey  string `json:"stepTypeKey"`
		DisplayOrder int    `json:"displayOrder"`
	} `json:"workoutStepTypes"`
	SportTypes     []WorkoutSportType `json:"workoutSportTypes"`
	ConditionTypes []struct {
		ConditionTypeID  int    `json:"conditionTypeId"`
		ConditionTypeKey string `json:"conditionTypeKey"`
		DisplayOrder     int    `json:"displayOrder"`
		Displayable      bool   `json:"displayable"`
	} `json:"workoutConditionTypes"`
	IntensityTypes []struct {
		IntensityTypeID  int    `json:"intensityTypeId"`
		IntensityTypeKey string `json:"intensityTypeKey"`
		DisplayOrder     int    `json:"displayOrder"`
	} `json:"workoutIntensityTypes"`
	TargetTypes []struct {
		WorkoutTargetTypeID  int    `json:"workoutTargetTypeId"`
		WorkoutTargetTypeKey string `json:"workoutTargetTypeKey"`
		DisplayOrder         int    `json:"displayOrder"`
	} `json:"workoutTargetTypes"`
	EquipmentTypes []struct {
		EquipmentTypeID  int    `json:"equipmentTypeId"`
		EquipmentTypeKey string `json:"equipmentTypeKey"`
		DisplayOrder     int    `json:"displayOrder"`
	} `json:"workoutEquipmentTypes"`
	StrokeTypes []struct {
		StrokeTypeID  int    `json:"strokeTypeId"`
		StrokeTypeKey string `json:"strokeTypeKey"`
		DisplayOrder  int    `json:"displayOrder"`
	} `json:"workoutStrokeTypes"`
	SwimInstructionTypes []struct {
		InstructionTypeID  int    `json:"instructionTypeId"`
		InstructionTypeKey string `json:"instructionTypeKey"`
		DisplayOrder       int    `json:"displayOrder"`
	} `json:"workoutSwimInstructionTypes"`
	DrillTypes []struct {
		DrillTypeID  int    `json:"drillTypeId"`
		DrillTypeKey string `json:"drillTypeKey"`
		DisplayOrder int    `json:"displayOrder"`
	} `json:"workoutDrillTypes"`
}

// WorkoutCount counts the saved workouts by sport.
type WorkoutCount struct {
	TotalCount   int `json:"totalCount"`
	RunningCount int `json:"runningCount"`
	CyclingCount int `json:"cyclingCount"`
	OtherCount   int `json:"otherCount"`
}

// StrengthBenchmark is an exercise usable as a one-rep-max based target.
type StrengthBenchmark struct {
	BenchmarkKey         string `json:"benchmarkKey"`
	BenchmarkCategoryKey string `json:"benchmarkCategoryKey"`
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

// ScheduleWorkout assigns a saved workout to a calendar day.
func (c *Client) ScheduleWorkout(ctx context.Context, workoutID int64, d time.Time) (*ScheduledWorkout, error) {
	var out ScheduledWorkout
	body := map[string]any{"date": date(d)}
	if err := c.post(ctx, fmt.Sprintf("/workout-service/schedule/%d", workoutID), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ScheduledWorkout returns one scheduled workout.
func (c *Client) ScheduledWorkout(ctx context.Context, scheduledWorkoutID int64) (*ScheduledWorkout, error) {
	var out ScheduledWorkout
	if err := c.get(ctx, fmt.Sprintf("/workout-service/schedule/%d", scheduledWorkoutID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WorkoutSchedule returns the workouts scheduled between start and end,
// including those from adaptive training plans.
func (c *Client) WorkoutSchedule(ctx context.Context, start, end time.Time) ([]ScheduledWorkout, error) {
	params := url.Values{"startDate": {date(start)}, "endDate": {date(end)}}
	var out []ScheduledWorkout
	if err := c.get(ctx, "/workout-service/schedule/workoutScheduleByPeriod", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UnscheduleWorkout removes a workout from the calendar.
func (c *Client) UnscheduleWorkout(ctx context.Context, scheduledWorkoutID int64) error {
	return c.del(ctx, fmt.Sprintf("/workout-service/schedule/%d", scheduledWorkoutID))
}

// RescheduleWorkout moves a scheduled workout to another calendar day.
func (c *Client) RescheduleWorkout(ctx context.Context, scheduledWorkoutID int64, d time.Time) error {
	path := fmt.Sprintf("/workout-service/schedule/%d", scheduledWorkoutID)
	// Garmin wants the whole schedule record back (a bare date is a 500), with
	// the embedded workout reduced to its ID.
	var sched map[string]json.RawMessage
	if err := c.get(ctx, path, nil, &sched); err != nil {
		return err
	}
	var w struct {
		WorkoutID int64 `json:"workoutId"`
	}
	if err := json.Unmarshal(sched["workout"], &w); err != nil {
		return fmt.Errorf("scheduled workout %d: %w", scheduledWorkoutID, err)
	}
	ref, err := json.Marshal(w)
	if err != nil {
		return err
	}
	sched["workout"] = ref
	sched["calendarDate"], err = json.Marshal(date(d))
	if err != nil {
		return err
	}
	return c.put(ctx, path, sched, nil)
}

// CreateWorkout saves a new workout and returns it as stored by Garmin. The
// definition uses Garmin's workout JSON (workoutName, sportType,
// workoutSegments, ...), e.g. the map returned by Workout.
func (c *Client) CreateWorkout(ctx context.Context, workout any) (map[string]json.RawMessage, error) {
	var out map[string]json.RawMessage
	if err := c.post(ctx, "/workout-service/workout", workout, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateWorkout replaces the workout with the given definition, typically the
// map from Workout with some fields changed.
func (c *Client) UpdateWorkout(ctx context.Context, id int64, workout map[string]json.RawMessage) error {
	body := maps.Clone(workout)
	idJSON, err := json.Marshal(id)
	if err != nil {
		return err
	}
	body["workoutId"] = idJSON
	return c.put(ctx, fmt.Sprintf("/workout-service/workout/%d", id), body, nil)
}

// CopyWorkout duplicates a saved workout under a new name and returns the copy.
func (c *Client) CopyWorkout(ctx context.Context, id int64, name string) (*Workout, error) {
	var out Workout
	if err := c.post(ctx, fmt.Sprintf("/workout-service/workout/copy/%d", id), plainText(name), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DownloadWorkout returns the FIT file for a saved workout.
func (c *Client) DownloadWorkout(ctx context.Context, id int64) ([]byte, error) {
	return c.getBytes(ctx, fmt.Sprintf("/workout-service/workout/FIT/%d", id), nil)
}

// WorkoutTypes returns the step, sport, condition, target and other type IDs
// used in workout definitions.
func (c *Client) WorkoutTypes(ctx context.Context) (*WorkoutTypes, error) {
	var out WorkoutTypes
	if err := c.get(ctx, "/workout-service/workout/types", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// WorkoutSportTypes returns the sports that have saved workouts.
func (c *Client) WorkoutSportTypes(ctx context.Context) ([]WorkoutSportType, error) {
	var out []WorkoutSportType
	if err := c.get(ctx, "/workout-service/workouts/sport-types", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// WorkoutCount returns how many workouts are saved, by sport.
func (c *Client) WorkoutCount(ctx context.Context) (*WorkoutCount, error) {
	var out WorkoutCount
	if err := c.get(ctx, "/workout-service/workouts/count", url.Values{"includeAtp": {"false"}}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// StrengthBenchmarks returns the exercises usable as one-rep-max targets.
func (c *Client) StrengthBenchmarks(ctx context.Context) ([]StrengthBenchmark, error) {
	var out []StrengthBenchmark
	if err := c.get(ctx, "/workout-service/benchmarks", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
