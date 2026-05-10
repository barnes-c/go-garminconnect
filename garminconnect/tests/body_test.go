package garminconnect_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyComposition(t *testing.T) {
	c, stop := newVCRClient(t, "body_composition")
	defer stop()

	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	bc, err := c.BodyComposition(start, testDate)
	require.NoError(t, err)

	assert.Equal(t, "2026-05-01", bc.StartDate)
	assert.Equal(t, "2026-05-10", bc.EndDate)
	assert.Equal(t, 78500.0, bc.TotalAverage.Weight)
	assert.Equal(t, 24.1, bc.TotalAverage.Bmi)
	assert.Equal(t, 18.5, bc.TotalAverage.BodyFat)
	assert.Len(t, bc.DateWeightList, 1)
}
