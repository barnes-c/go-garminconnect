package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNutritionStatus(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.NutritionStatus(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, out.CurrentStatus)
}
