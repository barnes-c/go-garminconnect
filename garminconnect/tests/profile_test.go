package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserProfile(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserProfile(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestUserProfileSettings(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserProfileSettings(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotEmpty(t, out.MeasurementSystem)
}

func TestUserSettings(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserSettings(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotEmpty(t, out.UserData.MeasurementSystem)
}

func TestUnitSystem(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UnitSystem(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}
