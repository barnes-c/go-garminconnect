package garminconnect_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkouts(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	workouts, err := c.Workouts(t.Context(), 0, 5)
	require.NoError(t, err)
	assert.NotNil(t, workouts)
	if len(workouts) > 0 {
		assert.NotZero(t, workouts[0].WorkoutID)
		assert.NotEmpty(t, workouts[0].WorkoutName)
	}
}

func TestWorkout(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	workouts, err := c.Workouts(t.Context(), 0, 1)
	require.NoError(t, err)
	if len(workouts) == 0 {
		return
	}

	out, err := c.Workout(t.Context(), workouts[0].WorkoutID)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestScheduledWorkouts(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Same endpoint as the deprecated ScheduledWorkouts.
	month, err := c.CalendarMonth(t.Context(), testDate.Year(), testDate.Month())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, month.Items)
	assert.NotEmpty(t, month.Items[0].ItemType)
}

// TestWorkoutLifecycle exercises every workout write against throwaway
// workouts, deleting them at the end so recording leaves the account unchanged.
func TestWorkoutLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	running := map[string]any{"sportTypeId": 1, "sportTypeKey": "running"}
	created, err := c.CreateWorkout(ctx, map[string]any{
		"workoutName": "go-garminconnect workout test",
		"sportType":   running,
		"workoutSegments": []any{map[string]any{
			"segmentOrder": 1,
			"sportType":    running,
			"workoutSteps": []any{map[string]any{
				"type":              "ExecutableStepDTO",
				"stepOrder":         1,
				"stepType":          map[string]any{"stepTypeId": 3, "stepTypeKey": "interval"},
				"endCondition":      map[string]any{"conditionTypeId": 2, "conditionTypeKey": "time"},
				"endConditionValue": 600,
				"targetType":        map[string]any{"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"},
			}},
		}},
	})
	require.NoError(t, err)
	var id int64
	require.NoError(t, json.Unmarshal(created["workoutId"], &id))
	require.NotZero(t, id)

	w, err := c.Workout(ctx, id)
	require.NoError(t, err)
	w["workoutName"] = json.RawMessage(`"go-garminconnect workout test (renamed)"`)
	require.NoError(t, c.UpdateWorkout(ctx, id, w))

	cp, err := c.CopyWorkout(ctx, id, "go-garminconnect workout copy")
	require.NoError(t, err)
	require.NotZero(t, cp.WorkoutID)
	assert.NotEmpty(t, cp.WorkoutName)

	sched, err := c.ScheduleWorkout(ctx, id, "2026-01-05")
	require.NoError(t, err)
	var schedID int64
	require.NoError(t, json.Unmarshal(sched["workoutScheduleId"], &schedID))
	require.NotZero(t, schedID)
	require.NoError(t, c.RescheduleWorkout(ctx, schedID, "2026-01-06"))
	require.NoError(t, c.UnscheduleWorkout(ctx, schedID))

	require.NoError(t, c.DeleteWorkout(ctx, cp.WorkoutID))
	require.NoError(t, c.DeleteWorkout(ctx, id))
}
