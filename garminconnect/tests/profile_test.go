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
	assert.NotZero(t, out.UserProfilePK)
	assert.NotEmpty(t, out.ProfileImageURL)
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

func TestUpdateUserSettings(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	s, err := c.UserSettings(ctx)
	skipAPIError(t, err)
	require.NoError(t, err)
	orig := s.UserData.TimeFormat
	next := "time_twenty_four_hr"
	if orig == next {
		next = "time_twelve_hr"
	}
	set := func(v string) error {
		return c.UpdateUserSettings(ctx, map[string]any{"userData": map[string]any{"timeFormat": v}})
	}
	require.NoError(t, set(next))
	require.NoError(t, set(orig))
}

func TestPersonalInformation(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.PersonalInformation(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.UserInfo.Email)
	assert.NotZero(t, out.BiometricProfile.UserID)
}

func TestUserProfileBase(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserProfileBase(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotZero(t, out.UserProfilePK)
	assert.NotEmpty(t, out.DisplayName)
}

func TestPrivacySettings(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.PrivacySettings(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.Activity.CategoryDetails.CategoryAccessRule)
	require.NotEmpty(t, out.Activity.Subcategories)
	assert.NotEmpty(t, out.Activity.Subcategories[0].Details.AccessRule)
}

func TestUsageIndicators(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UsageIndicators(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.PerformanceBasedIndicators.ActivityIndicators)
}

func TestUserFocus(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserFocus(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.UserFocus)
}

func TestUserLocation(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.UserLocation(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out.CountryCode)
	assert.NotZero(t, out.UserProfileID)
}

func TestFollowingsCount(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	_, err := c.FollowingsCount(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
}
