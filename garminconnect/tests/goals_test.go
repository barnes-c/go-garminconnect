package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoals(t *testing.T) {
	c, stop := newVCRClient(t, "goals")
	defer stop()

	goals, err := c.Goals(t.Context(), "active", 0, 10)
	require.NoError(t, err)
	assert.NotNil(t, goals)
}

func TestEarnedBadges(t *testing.T) {
	c, stop := newVCRClient(t, "earned_badges")
	defer stop()

	badges, err := c.EarnedBadges(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, badges)
}

func TestAvailableBadges(t *testing.T) {
	c, stop := newVCRClient(t, "available_badges")
	defer stop()

	badges, err := c.AvailableBadges(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, badges)
}
