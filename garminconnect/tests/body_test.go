package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeighIns(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.WeighIns(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDailyWeighIns(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.DailyWeighIns(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestLatestWeight(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.LatestWeight(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotZero(t, out.Weight)
	assert.NotEmpty(t, out.SourceType)
}

func TestBodyComposition(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	bc, err := c.BodyComposition(t.Context(), testDate, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, bc)
}
