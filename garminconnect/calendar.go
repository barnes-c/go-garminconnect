package garminconnect

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// CalendarItem is one entry on the Garmin Connect calendar: an activity,
// scheduled workout, event, note, badge, and so on.
type CalendarItem struct {
	ID           int64   `json:"id"`
	ItemType     string  `json:"itemType"` // e.g. "activity", "workout", "event", "note"
	Title        string  `json:"title"`
	Date         string  `json:"date"` // YYYY-MM-DD
	SportTypeKey string  `json:"sportTypeKey"`
	WorkoutID    int64   `json:"workoutId"`
	Distance     float64 `json:"distance"` // meters
	Duration     float64 `json:"duration"` // seconds
	Calories     float64 `json:"calories"`
	IsRace       bool    `json:"isRace"`
	PrimaryEvent bool    `json:"primaryEvent"`
}

// CalendarMonth is the calendar for one month.
type CalendarMonth struct {
	Year             int            `json:"year"`
	Month            int            `json:"month"` // zero-based, as Garmin sends it
	NumOfDaysInMonth int            `json:"numOfDaysInMonth"`
	Items            []CalendarItem `json:"calendarItems"`
}

// CalendarWeek is the calendar for one week.
type CalendarWeek struct {
	StartDate string         `json:"startDate"`
	EndDate   string         `json:"endDate"`
	Items     []CalendarItem `json:"calendarItems"`
}

// CalendarMonth returns every calendar item in the given month.
func (c *Client) CalendarMonth(ctx context.Context, year int, month time.Month) (*CalendarMonth, error) {
	var out CalendarMonth
	if err := c.get(ctx, fmt.Sprintf("/calendar-service/year/%d/month/%d", year, month-1), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CalendarWeek returns every calendar item in the week containing day, with
// weeks starting on weekStart.
func (c *Client) CalendarWeek(ctx context.Context, day time.Time, weekStart time.Weekday) (*CalendarWeek, error) {
	path := fmt.Sprintf("/calendar-service/year/%d/month/%d/day/%d/start/%d", day.Year(), day.Month()-1, day.Day(), weekStart)
	var out CalendarWeek
	if err := c.get(ctx, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// MoveCalendarItem moves a calendar item (e.g. itemType "note" or "event") to
// another day.
func (c *Client) MoveCalendarItem(ctx context.Context, itemType string, id int64, to time.Time) error {
	path := fmt.Sprintf("/calendar-service/calendarItem/%s/%d/date/%s", itemType, id, date(to))
	return c.put(ctx, path, map[string]any{"date": date(to)}, nil)
}

// CalendarEvent is a race or other event on the calendar.
type CalendarEvent struct {
	ID            int64               `json:"id,omitempty"`
	EventName     string              `json:"eventName"`
	Date          string              `json:"date"`      // YYYY-MM-DD
	EventType     string              `json:"eventType"` // e.g. "running"; Garmin can't update or move an event without one
	Race          bool                `json:"race"`
	Privacy       *EventPrivacy       `json:"eventPrivacy,omitempty"`
	Customization *EventCustomization `json:"eventCustomization,omitempty"`
}

// EventPrivacy controls who can see an event.
type EventPrivacy struct {
	Label          string `json:"label"` // e.g. "PRIVATE"
	IsShareable    bool   `json:"isShareable"`
	IsDiscoverable bool   `json:"isDiscoverable"`
}

// EventCustomization holds the user's settings for an event.
type EventCustomization struct {
	IsPrimaryEvent  bool `json:"isPrimaryEvent"`
	IsTrainingEvent bool `json:"isTrainingEvent"`
}

// CalendarEvents returns the events between start and end.
func (c *Client) CalendarEvents(ctx context.Context, start, end time.Time) ([]CalendarEvent, error) {
	params := url.Values{"startDate": {date(start)}, "endDate": {date(end)}}
	var out []CalendarEvent
	if err := c.get(ctx, "/calendar-service/events", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UpcomingEvents returns up to limit events in the next days days.
func (c *Client) UpcomingEvents(ctx context.Context, days, limit int) ([]CalendarEvent, error) {
	params := url.Values{"numDaysForward": {fmt.Sprintf("%d", days)}, "limit": {fmt.Sprintf("%d", limit)}}
	var out []CalendarEvent
	if err := c.get(ctx, "/calendar-service/events/upcoming", params, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PrimaryEvent returns the event marked as primary, or nil if there is none.
func (c *Client) PrimaryEvent(ctx context.Context) (*CalendarEvent, error) {
	var out CalendarEvent
	if err := c.get(ctx, "/calendar-service/event/primary", nil, &out); err != nil {
		var ae *APIError
		if errors.As(err, &ae) && ae.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

// CalendarEvent returns a single event.
func (c *Client) CalendarEvent(ctx context.Context, id int64) (*CalendarEvent, error) {
	var out CalendarEvent
	if err := c.get(ctx, fmt.Sprintf("/calendar-service/event/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCalendarEvent adds an event and returns it as stored by Garmin.
func (c *Client) CreateCalendarEvent(ctx context.Context, e CalendarEvent) (*CalendarEvent, error) {
	var out CalendarEvent
	if err := c.post(ctx, "/calendar-service/event", e, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCalendarEvent saves e, identified by e.ID.
func (c *Client) UpdateCalendarEvent(ctx context.Context, e *CalendarEvent) error {
	return c.put(ctx, fmt.Sprintf("/calendar-service/event/%d", e.ID), e, nil)
}

// DeleteCalendarEvent removes an event.
func (c *Client) DeleteCalendarEvent(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/calendar-service/event/%d", id))
}

// CalendarNote is a free-text note on a calendar day.
type CalendarNote struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"noteName"`
	Content string `json:"content"`
	Date    string `json:"date"` // YYYY-MM-DD
}

// CreateCalendarNote adds a note to the given day.
func (c *Client) CreateCalendarNote(ctx context.Context, day time.Time, name, content string) (*CalendarNote, error) {
	var out CalendarNote
	body := CalendarNote{Name: name, Content: content, Date: date(day)}
	if err := c.post(ctx, "/calendar-service/note", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCalendarNote saves n, identified by n.ID.
func (c *Client) UpdateCalendarNote(ctx context.Context, n *CalendarNote) error {
	return c.put(ctx, fmt.Sprintf("/calendar-service/note/%d", n.ID), n, nil)
}

// DeleteCalendarNote removes a note.
func (c *Client) DeleteCalendarNote(ctx context.Context, id int64) error {
	return c.del(ctx, fmt.Sprintf("/calendar-service/note/%d", id))
}
