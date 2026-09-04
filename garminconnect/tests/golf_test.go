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
	assert.NotEmpty(t, out)
}

func TestGolfUserStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.GolfUserStats(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}
