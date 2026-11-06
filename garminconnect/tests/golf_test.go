package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolfClubStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.GolfClubStats(t.Context(), 100)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].ClubTypeID)
}

func TestGolfUserStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.GolfUserStats(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestGolfClubTypes(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.GolfClubTypes(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].Name)
	assert.NotZero(t, out[0].LoftAngle)
}
