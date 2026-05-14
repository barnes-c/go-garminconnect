package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGear(t *testing.T) {
	c, stop := newVCRClient(t, "gear")
	defer stop()

	// Record cassette: fetch user summary to get the profile number, then gear.
	summary, err := c.UserSummary(testDate)
	require.NoError(t, err)

	gear, err := c.Gear(summary.UserProfileID)
	require.NoError(t, err)
	assert.NotNil(t, gear)
}

func TestGearStats(t *testing.T) {
	c, stop := newVCRClient(t, "gear_stats")
	defer stop()

	summary, err := c.UserSummary(testDate)
	require.NoError(t, err)

	gear, err := c.Gear(summary.UserProfileID)
	require.NoError(t, err)
	if len(gear) == 0 {
		return
	}

	stats, err := c.GearStats(gear[0].UUID)
	require.NoError(t, err)
	assert.NotEmpty(t, stats)
}
