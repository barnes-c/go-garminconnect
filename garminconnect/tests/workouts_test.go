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

	sports, err := c.WorkoutSportTypes(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, sports)
	assert.NotEmpty(t, sports[0].SportTypeKey)
	count, err := c.WorkoutCount(ctx)
	require.NoError(t, err)
	assert.NotZero(t, count.TotalCount)

	// Scheduled on or before testDate so the period query below covers it.
	sched, err := c.ScheduleWorkout(ctx, id, testDate.AddDate(0, 0, -3))
	require.NoError(t, err)
	schedID := sched.WorkoutScheduleID
	require.NotZero(t, schedID)
	require.NoError(t, c.RescheduleWorkout(ctx, schedID, testDate.AddDate(0, 0, -2)))

	got, err := c.ScheduledWorkout(ctx, schedID)
	require.NoError(t, err)
	assert.NotZero(t, got.Workout.WorkoutID)
	assert.NotEmpty(t, got.CalendarDate)
	period, err := c.WorkoutSchedule(ctx, testDate.AddDate(0, 0, -7), testDate)
	require.NoError(t, err)
	require.NotEmpty(t, period)
	assert.NotZero(t, period[0].WorkoutScheduleID)

	require.NoError(t, c.UnscheduleWorkout(ctx, schedID))

	require.NoError(t, c.DeleteWorkout(ctx, cp.WorkoutID))
	require.NoError(t, c.DeleteWorkout(ctx, id))
}

func TestWorkoutTypes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.WorkoutTypes(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out.StepTypes)
	assert.NotEmpty(t, out.StepTypes[0].StepTypeKey)
	assert.NotEmpty(t, out.SportTypes)
	assert.NotEmpty(t, out.ConditionTypes)
	assert.NotEmpty(t, out.TargetTypes)
}

func TestStrengthBenchmarks(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.StrengthBenchmarks(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].BenchmarkKey)
}
