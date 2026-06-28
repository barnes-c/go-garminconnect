package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrainingReadiness(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	entries, err := c.TrainingReadiness(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	assert.NotEmpty(t, entries[0].CalendarDate)
}

func TestTrainingStatus(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	entries, err := c.TrainingStatus(t.Context(), testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestMaxMetrics(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.MaxMetrics(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestEnduranceScore(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.EnduranceScore(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestRacePredictions(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	preds, err := c.RacePredictions(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, preds)
}

func TestHillScore(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.HillScore(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestLactateThreshold(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.LactateThreshold(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestFitnessAge(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.FitnessAge(t.Context(), testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRunningTolerance(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.RunningTolerance(t.Context(), start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	if len(out) == 0 {
		t.Skip("no running tolerance data in cassette")
	}
	assert.NotEmpty(t, out)
}

func TestCyclingFTP(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.CyclingFTP(t.Context())
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}
