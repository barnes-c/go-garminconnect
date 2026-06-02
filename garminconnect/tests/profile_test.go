package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserProfile(t *testing.T) {
	c, stop := newVCRClient(t, "user_profile")
	defer stop()

	out, err := c.UserProfile(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestUserProfileSettings(t *testing.T) {
	c, stop := newVCRClient(t, "user_profile_settings")
	defer stop()

	out, err := c.UserProfileSettings(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}
