package garminconnect_test

import (
	"testing"
	"time"

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

// TestBodyWritesLifecycle records a weigh-in and a blood pressure reading,
// edits the reading's note, then deletes both so recording leaves the account
// unchanged.
func TestBodyWritesLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()
	day := testDate.Format(time.DateOnly)

	_, err := c.AddWeighIn(ctx, 75.5, "kg", day+"T08:00:00.00")
	require.NoError(t, err)

	zurich, err := time.LoadLocation("Europe/Zurich")
	require.NoError(t, err)
	at := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 8, 0, 0, 0, zurich)
	_, err = c.SetBloodPressure(ctx, 120, 80, 60, at, "created by TestBodyWritesLifecycle")
	require.NoError(t, err)

	bp, err := c.BloodPressure(ctx, testDate, testDate)
	require.NoError(t, err)
	require.NotEmpty(t, bp.Days)
	require.NotEmpty(t, bp.Days[0].Measurements)
	m := bp.Days[0].Measurements[0]
	assert.NotZero(t, m.Systolic)
	require.NotZero(t, m.Version)
	require.NoError(t, c.SetBloodPressureNotes(ctx, m, "edited"))
	require.NoError(t, c.DeleteBloodPressure(ctx, day, m.Version))

	n, err := c.DeleteWeighIns(ctx, testDate)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
}
