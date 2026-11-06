package garminconnect_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReportLifecycle adds a blood pressure reading so the report has data,
// exports a blood pressure report, waits for it and downloads it, then
// deletes the reading.
func TestReportLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	at := time.Date(testDate.Year(), testDate.Month(), testDate.Day(), 8, 0, 0, 0, time.UTC)
	_, err := c.SetBloodPressure(ctx, 120, 80, 60, at, "created by TestReportLifecycle")
	skipAPIError(t, err)
	require.NoError(t, err)

	export, err := c.RequestReport(ctx, "bloodPressure", testDate.AddDate(0, 0, -27), testDate,
		map[string]any{"reportSubCategory": "FourWeeks", "reportFormat": "PDF"})
	require.NoError(t, err)
	require.NotEmpty(t, export.ReportID)

	// One poll after a generous wait: replay can't loop on the status, which
	// the sanitizer scrubs to "TEST".
	if os.Getenv("GARMIN_TOKEN") != "" {
		time.Sleep(10 * time.Second)
	}
	st, err := c.ReportStatus(ctx, export.ReportID)
	require.NoError(t, err)
	assert.NotEmpty(t, st.Progress.Status)
	if os.Getenv("GARMIN_TOKEN") != "" {
		require.Equal(t, "READY", st.Progress.Status)
	}

	pdf, err := c.DownloadReport(ctx, export.ReportID)
	require.NoError(t, err)
	assert.NotEmpty(t, pdf)

	bp, err := c.BloodPressure(ctx, testDate, testDate)
	require.NoError(t, err)
	require.NotEmpty(t, bp.Days)
	require.NotEmpty(t, bp.Days[0].Measurements)
	m := bp.Days[0].Measurements[0]
	require.NoError(t, c.DeleteBloodPressure(ctx, testDate.Format(time.DateOnly), m.Version))
}
