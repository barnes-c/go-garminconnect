package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTrainingPlanLifecycle schedules a catalog plan and deletes it again,
// which also removes the workouts and notes it put on the calendar. The start
// date is testDate, since the sanitizer rewrites later URL dates.
func TestTrainingPlanLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	catalog, err := c.SearchTrainingPlans(ctx, 1, 3)
	require.NoError(t, err)
	require.NotEmpty(t, catalog)
	plan := catalog[0]
	assert.NotEmpty(t, plan.Name)
	assert.NotZero(t, plan.DurationInWeeks)

	tasks, err := c.TrainingPlanTasks(ctx, plan.TrainingPlanID)
	require.NoError(t, err)
	require.NotEmpty(t, tasks)
	assert.NotZero(t, tasks[0].WeekID)

	scheduled, err := c.ScheduleTrainingPlan(ctx, plan.TrainingPlanID, testDate)
	require.NoError(t, err)
	require.NotZero(t, scheduled.TrainingPlanID)
	assert.NotEmpty(t, scheduled.TrainingStatus.StatusKey)

	mine, err := c.TrainingPlans(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, mine)
	_, err = c.ActiveTrainingPlans(ctx, testDate)
	require.NoError(t, err)

	require.NoError(t, c.DeleteTrainingPlan(ctx, scheduled.TrainingPlanID))
}
