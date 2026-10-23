package garminconnect

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
)

// CourseSummary is a course as listed by Courses and FavoriteCourses.
type CourseSummary struct {
	CourseID              int64   `json:"courseId"`
	CourseName            string  `json:"courseName"`
	CourseDescription     string  `json:"courseDescription"`
	DistanceInMeters      float64 `json:"distanceInMeters"`
	ElevationGainInMeters float64 `json:"elevationGainInMeters"`
	ElevationLossInMeters float64 `json:"elevationLossInMeters"`
	StartLatitude         float64 `json:"startLatitude"`
	StartLongitude        float64 `json:"startLongitude"`
	SourceTypeID          int     `json:"sourceTypeId"`
	Favorite              bool    `json:"favorite"`
	CreatedDate           int64   `json:"createdDate"` // epoch milliseconds
	UpdatedDate           int64   `json:"updatedDate"` // epoch milliseconds
}

// GeoPoint is one point of a course track.
type GeoPoint struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Elevation *float64 `json:"elevation"` // meters; nil when unknown
	Distance  float64  `json:"distance"`  // meters from the start
}

// Course is a full course: a saved one from Course, or an unsaved draft from
// ImportCourse or ConvertActivityToCourse.
type Course struct {
	CourseID           int64      `json:"courseId,omitempty"`
	CourseName         string     `json:"courseName"`
	Description        string     `json:"description,omitempty"`
	RulePK             int        `json:"rulePK"` // privacy: 1 public, 2 private, 4 group
	ActivityTypePK     int        `json:"activityTypePk"`
	SourceTypeID       int        `json:"sourceTypeId"`
	SourcePK           int64      `json:"sourcePk,omitempty"` // e.g. the activity a course was converted from
	StartPoint         *GeoPoint  `json:"startPoint"`
	GeoPoints          []GeoPoint `json:"geoPoints"`
	DistanceMeter      float64    `json:"distanceMeter"`
	ElevationGainMeter float64    `json:"elevationGainMeter"`
	ElevationLossMeter float64    `json:"elevationLossMeter"`
}

// Courses returns the user's saved courses.
func (c *Client) Courses(ctx context.Context) ([]CourseSummary, error) {
	name, err := c.displayNamePath()
	if err != nil {
		return nil, err
	}
	var out []CourseSummary
	if err := c.get(ctx, "/course-service/course/owner/"+name, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FavoriteCourses returns the courses marked as favorite.
func (c *Client) FavoriteCourses(ctx context.Context) ([]CourseSummary, error) {
	var out []CourseSummary
	if err := c.get(ctx, "/course-service/course/favorites", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Course returns a saved course with its track.
func (c *Client) Course(ctx context.Context, id int64) (*Course, error) {
	var out Course
	if err := c.get(ctx, fmt.Sprintf("/course-service/course/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ImportCourse parses a GPX, TCX or FIT file into an unsaved course draft.
// Set ActivityTypePK on the draft, then save it with CreateCourse.
func (c *Client) ImportCourse(ctx context.Context, data []byte, filename string) (*Course, error) {
	var out Course
	if err := c.upload(ctx, "/course-service/course/import", data, filename, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ConvertActivityToCourse builds an unsaved course draft from a GPS activity.
// Save it with CreateCourse.
func (c *Client) ConvertActivityToCourse(ctx context.Context, activityID int64) (*Course, error) {
	var out Course
	if err := c.get(ctx, fmt.Sprintf("/course-service/course/convert/activity/%d", activityID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCourse saves a course draft and returns the stored course. Fields an
// imported draft leaves empty get defaults: StartPoint is the first point,
// RulePK is private and SourceTypeID is 1. ActivityTypePK must be set.
func (c *Client) CreateCourse(ctx context.Context, course Course) (*Course, error) {
	if course.StartPoint == nil && len(course.GeoPoints) > 0 {
		course.StartPoint = &course.GeoPoints[0]
	}
	if course.RulePK == 0 {
		course.RulePK = 2
	}
	if course.SourceTypeID == 0 {
		course.SourceTypeID = 1
	}
	var out Course
	if err := c.post(ctx, "/course-service/course", course, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCourse saves course, identified by course.CourseID.
func (c *Client) UpdateCourse(ctx context.Context, course *Course) error {
	path := fmt.Sprintf("/course-service/course/%d", course.CourseID)
	// Garmin rejects partial records, and the record has fields Course
	// doesn't model, so overlay course onto the current raw record.
	var current map[string]json.RawMessage
	if err := c.get(ctx, path, nil, &current); err != nil {
		return err
	}
	data, err := json.Marshal(course)
	if err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	maps.Copy(current, fields)
	return c.put(ctx, path, current, nil)
}

// DeleteCourse permanently deletes a course.
func (c *Client) DeleteCourse(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/course-service/course/%d", id))
}

// SetCourseFavorite marks or unmarks a course as favorite.
func (c *Client) SetCourseFavorite(ctx context.Context, id int64, favorite bool) error {
	if favorite {
		return c.post(ctx, "/course-service/course/favorites", map[string]any{"courseId": id}, nil)
	}
	return c.del(ctx, fmt.Sprintf("/course-service/course/favorites/course/%d", id))
}

// DownloadCourseGPX returns a course as a GPX file.
func (c *Client) DownloadCourseGPX(ctx context.Context, id int64) ([]byte, error) {
	return c.getBytes(ctx, fmt.Sprintf("/course-service/course/gpx/%d", id), nil)
}

// DownloadCourseFIT returns a course as a FIT file, including elevation.
func (c *Client) DownloadCourseFIT(ctx context.Context, id int64) ([]byte, error) {
	return c.getBytes(ctx, fmt.Sprintf("/course-service/course/fit/%d/0", id), url.Values{"elevation": {"true"}})
}
