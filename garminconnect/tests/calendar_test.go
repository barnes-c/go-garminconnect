package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalendarLifecycle adds a note and an event, reads them back through
// every calendar view, edits and moves them, then deletes both so recording
// leaves the account unchanged. All dates are on or before testDate, since
// the sanitizer rewrites later URL dates.
func TestCalendarLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	note, err := c.CreateCalendarNote(ctx, testDate, "go-garminconnect note", "created by TestCalendarLifecycle")
	require.NoError(t, err)
	require.NotZero(t, note.ID)

	ev, err := c.CreateCalendarEvent(ctx, gc.CalendarEvent{
		EventName: "go-garminconnect event",
		Date:      testDate.Format(time.DateOnly),
		EventType: "running",
	})
	require.NoError(t, err)
	require.NotZero(t, ev.ID)

	month, err := c.CalendarMonth(ctx, testDate.Year(), testDate.Month())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(month.Items), 2)
	assert.NotEmpty(t, month.Items[0].ItemType)

	week, err := c.CalendarWeek(ctx, testDate, time.Monday)
	require.NoError(t, err)
	assert.NotEmpty(t, week.Items)

	events, err := c.CalendarEvents(ctx, testDate.AddDate(0, 0, -30), testDate)
	require.NoError(t, err)
	require.NotEmpty(t, events)

	_, err = c.UpcomingEvents(ctx, 90, 10)
	require.NoError(t, err)

	got, err := c.CalendarEvent(ctx, ev.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, got.EventType)
	got.EventName = "go-garminconnect event (renamed)"
	got.Customization = &gc.EventCustomization{IsPrimaryEvent: true}
	require.NoError(t, c.UpdateCalendarEvent(ctx, got))

	primary, err := c.PrimaryEvent(ctx)
	require.NoError(t, err)
	require.NotNil(t, primary)
	assert.NotZero(t, primary.ID)

	note.Content = "edited"
	require.NoError(t, c.UpdateCalendarNote(ctx, note))
	require.NoError(t, c.MoveCalendarItem(ctx, "note", note.ID, testDate.AddDate(0, 0, -1)))
	require.NoError(t, c.MoveCalendarItem(ctx, "event", ev.ID, testDate.AddDate(0, 0, -1)))

	require.NoError(t, c.DeleteCalendarEvent(ctx, ev.ID))
	require.NoError(t, c.DeleteCalendarNote(ctx, note.ID))
}
