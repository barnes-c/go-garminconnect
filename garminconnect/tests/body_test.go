package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeighIns(t *testing.T) {
	c, stop := newVCRClient(t, "weigh_ins")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.WeighIns(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestDailyWeighIns(t *testing.T) {
	c, stop := newVCRClient(t, "daily_weigh_ins")
	defer stop()

	out, err := c.DailyWeighIns(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, out)
}

func TestBodyComposition(t *testing.T) {
	c, stop := newVCRClient(t, "body_composition")
	defer stop()

	bc, err := c.BodyComposition(testDate, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotNil(t, bc)
}
