package garminconnect_test

import (
	"testing"
	"time"

	gc "github.com/barnes-c/go-garminconnect/garminconnect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGear(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	// Record cassette: fetch user summary to get the profile number, then gear.
	summary, err := c.UserSummary(t.Context(), testDate)
	require.NoError(t, err)

	gear, err := c.Gear(t.Context(), summary.UserProfileID)
	require.NoError(t, err)
	require.NotEmpty(t, gear)
	assert.NotZero(t, gear[0].MaxMeters)
	assert.NotNil(t, gear)
}

func TestGearStats(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	summary, err := c.UserSummary(t.Context(), testDate)
	require.NoError(t, err)

	gear, err := c.Gear(t.Context(), summary.UserProfileID)
	require.NoError(t, err)
	if len(gear) == 0 {
		return
	}

	stats, err := c.GearStats(t.Context(), gear[0].UUID)
	require.NoError(t, err)
	assert.NotEmpty(t, stats)
}

// TestGearLifecycle exercises every gear write against one throwaway gear item
// and manual activity, deleting both at the end so recording leaves the
// account unchanged.
func TestGearLifecycle(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	activityID, err := c.CreateManualActivity(ctx, gc.ManualActivity{
		Name:     "go-garminconnect gear test",
		TypeKey:  "running",
		Start:    testDate.Add(8 * time.Hour),
		TimeZone: "Europe/Zurich",
		Distance: 5000,
		Duration: 30 * time.Minute,
	})
	require.NoError(t, err)
	require.NotZero(t, activityID)
	require.NoError(t, c.SetActivityDescription(ctx, activityID, "created by TestGearLifecycle"))

	g, err := c.CreateGear(ctx, gc.GearDetail{
		GearType:     "SHOES",
		Brand:        "go-garminconnect",
		Model:        "Lifecycle",
		Name:         "gear test",
		FirstUseDate: testDate.Format(time.DateOnly),
		UsageType:    "DISTANCE",
	})
	require.NoError(t, err)
	require.NotEmpty(t, g.UUID)

	d, err := c.GearDetails(ctx, g.UUID)
	require.NoError(t, err)
	assert.NotEmpty(t, d.Name)
	d.Notes = "updated"
	d.AssociatedActivityTypes = []gc.GearActivityType{{ActivityTypeKey: "running", DefaultGear: true}}
	require.NoError(t, c.UpdateGear(ctx, d))

	require.NoError(t, c.AddGearToActivity(ctx, g.UUID, activityID))
	require.NoError(t, c.RemoveGearFromActivity(ctx, g.UUID, activityID))
	require.NoError(t, c.SetGearStatus(ctx, g.UUID, gc.GearRetired))

	require.NoError(t, c.DeleteGear(ctx, g.UUID))
	require.NoError(t, c.DeleteActivity(ctx, activityID))
}

// TestGearCatalogAndLists reads Garmin's gear catalog, then links a throwaway
// gear item to a throwaway activity to read the user's gear lists, and
// deletes both.
func TestGearCatalogAndLists(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	types, err := c.GearTypes(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, types)
	assert.NotEmpty(t, types[0].GearTypeName)

	models, err := c.SearchGearMakeModels(ctx, "pegasus")
	require.NoError(t, err)
	require.NotEmpty(t, models)
	assert.NotEmpty(t, models[0].MakeModelName)

	g, err := c.CreateGear(ctx, gc.GearDetail{
		GearType: "SHOES", Brand: "go-garminconnect", Model: "Lists", Name: "gear list test",
		FirstUseDate: testDate.Format(time.DateOnly), UsageType: "DISTANCE",
	})
	require.NoError(t, err)
	activityID, err := c.CreateManualActivity(ctx, gc.ManualActivity{
		Name: "go-garminconnect gear list test", TypeKey: "running",
		Start: testDate.Add(8 * time.Hour), TimeZone: "Europe/Zurich",
		Distance: 5000, Duration: 30 * time.Minute,
	})
	require.NoError(t, err)
	require.NoError(t, c.AddGearToActivity(ctx, g.UUID, activityID))

	userTypes, err := c.UserGearTypes(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, userTypes)

	list, err := c.GearList(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	assert.NotEmpty(t, list[0].UUID)

	linked, err := c.ActivityGear(ctx, activityID)
	require.NoError(t, err)
	require.NotEmpty(t, linked)
	assert.NotEmpty(t, linked[0].Name)

	require.NoError(t, c.DeleteGear(ctx, g.UUID))
	require.NoError(t, c.DeleteActivity(ctx, activityID))
}
