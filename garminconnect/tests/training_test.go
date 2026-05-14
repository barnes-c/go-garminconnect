package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrainingReadiness(t *testing.T) {
	c, stop := newVCRClient(t, "training_readiness")
	defer stop()

	entries, err := c.TrainingReadiness(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	assert.NotEmpty(t, entries[0].CalendarDate)
}

func TestTrainingStatus(t *testing.T) {
	c, stop := newVCRClient(t, "training_status")
	defer stop()

	entries, err := c.TrainingStatus(testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestMaxMetrics(t *testing.T) {
	c, stop := newVCRClient(t, "max_metrics")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.MaxMetrics(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestEnduranceScore(t *testing.T) {
	c, stop := newVCRClient(t, "endurance_score")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.EnduranceScore(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestRacePredictions(t *testing.T) {
	c, stop := newVCRClient(t, "race_predictions")
	defer stop()

	preds, err := c.RacePredictions()
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, preds)
}

func TestHillScore(t *testing.T) {
	c, stop := newVCRClient(t, "hill_score")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	entries, err := c.HillScore(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestLactateThreshold(t *testing.T) {
	c, stop := newVCRClient(t, "lactate_threshold")
	defer stop()

	out, err := c.LactateThreshold()
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.Greater(t, len(out), 0)
}

func TestFitnessAge(t *testing.T) {
	c, stop := newVCRClient(t, "fitness_age")
	defer stop()

	out, err := c.FitnessAge(testDate)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRunningTolerance(t *testing.T) {
	c, stop := newVCRClient(t, "running_tolerance")
	defer stop()

	start := testDate.AddDate(0, -1, 0)
	out, err := c.RunningTolerance(start, testDate)
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestCyclingFTP(t *testing.T) {
	c, stop := newVCRClient(t, "cycling_ftp")
	defer stop()

	out, err := c.CyclingFTP()
	skipAPIError(t, err)
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}
