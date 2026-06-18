package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActiveTrainingPlan(t *testing.T) {
	c, stop := newVCRClient(t, "active_training_plan")
	defer stop()

	// A 204 (no active plan) is a valid response and yields a nil map.
	_, err := c.ActiveTrainingPlan(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
}

func TestCompletedTrainingPlans(t *testing.T) {
	c, stop := newVCRClient(t, "completed_training_plans")
	defer stop()

	out, err := c.CompletedTrainingPlans(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestTrainingPlanTypes(t *testing.T) {
	c, stop := newVCRClient(t, "training_plan_types")
	defer stop()

	out, err := c.TrainingPlanTypes(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].WorkoutPlanTypeID)
	assert.NotEmpty(t, out[0].Name)
	assert.NotEmpty(t, out[0].TrainingType)
}
