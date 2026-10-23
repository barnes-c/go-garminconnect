package garminconnect_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileDownloads uploads a throwaway GPS activity and creates a workout,
// downloads both in every supported format, then deletes them. Downloaded
// bodies are scrubbed to a placeholder in the cassette, so only non-emptiness
// is checked.
func TestFileDownloads(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	_, err := c.UploadActivity(ctx, syntheticGPX("go-garminconnect download", true), "download.gpx")
	require.NoError(t, err)
	var activityID int64
	for range 10 {
		acts, err := c.Activities(ctx, 0, 1)
		require.NoError(t, err)
		if len(acts) > 0 {
			activityID = acts[0].ActivityID
			break
		}
		if os.Getenv("GARMIN_TOKEN") != "" {
			time.Sleep(2 * time.Second)
		}
	}
	require.NotZero(t, activityID)

	for _, f := range []gc.DownloadFormat{gc.FormatOriginal, gc.FormatGPX, gc.FormatTCX, gc.FormatKML, gc.FormatCSV} {
		data, err := c.DownloadActivity(ctx, activityID, f)
		require.NoError(t, err, f)
		assert.NotEmpty(t, data, f)
	}

	running := map[string]any{"sportTypeId": 1, "sportTypeKey": "running"}
	created, err := c.CreateWorkout(ctx, map[string]any{
		"workoutName": "go-garminconnect download test",
		"sportType":   running,
		"workoutSegments": []any{map[string]any{
			"segmentOrder": 1,
			"sportType":    running,
			"workoutSteps": []any{map[string]any{
				"type":              "ExecutableStepDTO",
				"stepOrder":         1,
				"stepType":          map[string]any{"stepTypeId": 3, "stepTypeKey": "interval"},
				"endCondition":      map[string]any{"conditionTypeId": 2, "conditionTypeKey": "time"},
				"endConditionValue": 600,
				"targetType":        map[string]any{"workoutTargetTypeId": 1, "workoutTargetTypeKey": "no.target"},
			}},
		}},
	})
	require.NoError(t, err)
	var workoutID int64
	require.NoError(t, json.Unmarshal(created["workoutId"], &workoutID))
	fit, err := c.DownloadWorkout(ctx, workoutID)
	require.NoError(t, err)
	assert.NotEmpty(t, fit)

	require.NoError(t, c.DeleteWorkout(ctx, workoutID))
	require.NoError(t, c.DeleteActivity(ctx, activityID))
}
