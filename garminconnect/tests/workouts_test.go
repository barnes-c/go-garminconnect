package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkouts(t *testing.T) {
	c, stop := newVCRClient(t, "workouts")
	defer stop()

	workouts, err := c.Workouts(0, 5)
	require.NoError(t, err)
	if len(workouts) == 0 {
		t.Skip("no workouts in cassette")
	}
	assert.NotZero(t, workouts[0].WorkoutID)
	assert.NotEmpty(t, workouts[0].WorkoutName)
}

func TestWorkout(t *testing.T) {
	c, stop := newVCRClient(t, "workout_detail")
	defer stop()

	// Record cassette: fetch list first to get a real ID, then the detail.
	workouts, err := c.Workouts(0, 1)
	require.NoError(t, err)
	if len(workouts) == 0 {
		t.Skip("no workouts in cassette")
	}

	out, err := c.Workout(workouts[0].WorkoutID)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestScheduledWorkouts(t *testing.T) {
	c, stop := newVCRClient(t, "scheduled_workouts")
	defer stop()

	sw, err := c.ScheduledWorkouts(0, 5)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, sw)
}
