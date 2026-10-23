package garminconnect_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// syntheticGPX returns a short made-up track; withTimes makes it an activity
// rather than a route.
func syntheticGPX(name string, withTimes bool) []byte {
	var pts strings.Builder
	for i := range 30 {
		fmt.Fprintf(&pts, `<trkpt lat="%.4f" lon="%.4f"><ele>%d</ele>`, 46.0+float64(i)*0.0005, 7.0+float64(i)*0.0003, 500+i)
		if withTimes {
			fmt.Fprintf(&pts, `<time>2025-12-31T06:%02d:00Z</time>`, i)
		}
		pts.WriteString(`</trkpt>`)
	}
	return fmt.Appendf(nil, `<?xml version="1.0" encoding="UTF-8"?><gpx version="1.1" creator="go-garminconnect" `+
		`xmlns="http://www.topografix.com/GPX/1/1"><trk><name>%s</name><type>running</type><trkseg>%s</trkseg></trk></gpx>`, name, pts.String())
}

// TestCourseLifecycle imports, saves, edits, favorites and downloads a
// throwaway course, then deletes it so recording leaves the account unchanged.
func TestCourseLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	draft, err := c.ImportCourse(ctx, syntheticGPX("go-garminconnect course", false), "course.gpx")
	require.NoError(t, err)
	require.NotEmpty(t, draft.GeoPoints)
	draft.ActivityTypePK = 1

	course, err := c.CreateCourse(ctx, *draft)
	require.NoError(t, err)
	require.NotZero(t, course.CourseID)

	list, err := c.Courses(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	assert.NotZero(t, list[0].DistanceInMeters)

	got, err := c.Course(ctx, course.CourseID)
	require.NoError(t, err)
	assert.NotEmpty(t, got.GeoPoints)
	got.CourseName = "go-garminconnect course (renamed)"
	require.NoError(t, c.UpdateCourse(ctx, got))

	require.NoError(t, c.SetCourseFavorite(ctx, course.CourseID, true))
	favs, err := c.FavoriteCourses(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, favs)
	require.NoError(t, c.SetCourseFavorite(ctx, course.CourseID, false))

	gpx, err := c.DownloadCourseGPX(ctx, course.CourseID)
	require.NoError(t, err)
	assert.NotEmpty(t, gpx)
	fit, err := c.DownloadCourseFIT(ctx, course.CourseID)
	require.NoError(t, err)
	assert.NotEmpty(t, fit)

	require.NoError(t, c.DeleteCourse(ctx, course.CourseID))
}

// TestConvertActivityToCourse uploads a GPS activity, turns it into a course,
// then deletes both.
func TestConvertActivityToCourse(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	_, err := c.UploadActivity(ctx, syntheticGPX("go-garminconnect convert", true), "convert.gpx")
	require.NoError(t, err)
	// Uploads are processed asynchronously; poll until the activity appears.
	// Only wait between polls when recording against the live API.
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

	draft, err := c.ConvertActivityToCourse(ctx, activityID)
	require.NoError(t, err)
	assert.NotEmpty(t, draft.GeoPoints)
	assert.NotZero(t, draft.ActivityTypePK)

	course, err := c.CreateCourse(ctx, *draft)
	require.NoError(t, err)
	require.NotZero(t, course.CourseID)

	// Garmin answers 429 to deleting a course right after creating it.
	if os.Getenv("GARMIN_TOKEN") != "" {
		time.Sleep(10 * time.Second)
	}
	require.NoError(t, c.DeleteCourse(ctx, course.CourseID))
	require.NoError(t, c.DeleteActivity(ctx, activityID))
}
