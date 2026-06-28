package garminconnect_test

import (
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

	sw, err := c.ScheduledWorkouts(t.Context(), testDate.Year(), int(testDate.Month()))
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, sw)
}
