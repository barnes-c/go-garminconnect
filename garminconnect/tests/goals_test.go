package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoals(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	goals, err := c.Goals(t.Context(), "active", 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, goals)
}

func TestEarnedBadges(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	badges, err := c.EarnedBadges(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, badges)
	assert.NotEmpty(t, badges[0].EarnedDate)
	assert.NotEmpty(t, badges)
}

func TestAvailableBadges(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	badges, err := c.AvailableBadges(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, badges)
}

func TestBadgeCatalog(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	attrs, err := c.BadgeAttributes(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, attrs.Categories)
	assert.NotEmpty(t, attrs.Categories[0].Key)

	available, err := c.AvailableBadges(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, available)
	badge, err := c.BadgeDetail(ctx, available[0].BadgeID)
	require.NoError(t, err)
	assert.NotEmpty(t, badge.BadgeKey)

	challenges, err := c.AvailableVirtualChallenges(ctx, 1, 5)
	require.NoError(t, err)
	require.NotEmpty(t, challenges)
	assert.NotEmpty(t, challenges[0].UUID)

	types, err := c.PersonalRecordTypes(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, types)
	assert.NotEmpty(t, types[0].Key)
}

func TestRecentBadges(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	latest, err := c.LatestEarnedBadges(ctx, 5)
	require.NoError(t, err)
	require.NotEmpty(t, latest)
	assert.NotEmpty(t, latest[0].EarnedDate)

	// Badges are only earned on days with activity; see dataDate.
	daily, err := c.DailyBadges(ctx, dataDate(t, "GARMIN_SUMMARY_DATE"))
	require.NoError(t, err)
	require.NotEmpty(t, daily)
	assert.NotZero(t, daily[0].BadgeID)
}

// TestActivityAchievements reads the records and badges of a throwaway
// activity, then deletes it.
func TestActivityAchievements(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	id, err := c.CreateManualActivity(ctx, gc.ManualActivity{
		Name: "go-garminconnect achievements test", TypeKey: "running",
		Start: testDate.Add(8 * time.Hour), TimeZone: "Europe/Zurich",
		Distance: 5000, Duration: 30 * time.Minute,
	})
	require.NoError(t, err)

	_, err = c.ActivityPersonalRecords(ctx, id)
	require.NoError(t, err)
	_, err = c.ActivityBadges(ctx, id)
	require.NoError(t, err)

	require.NoError(t, c.DeleteActivity(ctx, id))
}

// TestDailyGoalLifecycle reads the daily goals, changes one and restores it.
func TestDailyGoalLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	steps, err := c.StepsGoal(ctx, testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, steps.GoalType)

	goals, err := c.EffectiveGoals(ctx, testDate)
	require.NoError(t, err)
	require.NotEmpty(t, goals)
	// A fixed goal type: the cassette scrubs goalType, so one read from the
	// response would change the URL on replay. Only the value comes from the
	// response, and request bodies don't affect replay matching.
	orig := gc.DailyGoal{GoalType: "INTENSITY_MINUTES", GoalValue: 150}
	for _, g := range goals {
		if g.GoalType == orig.GoalType {
			orig.GoalValue = g.GoalValue
		}
	}
	require.NoError(t, c.SetDailyGoal(ctx, testDate, gc.DailyGoal{GoalType: orig.GoalType, GoalValue: orig.GoalValue + 10}))
	require.NoError(t, c.SetDailyGoal(ctx, testDate, orig))
}
