package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevices(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	devices, err := c.Devices(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, devices)
}

func TestLastUsedDevice(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	d, err := c.LastUsedDevice(t.Context())
	require.NoError(t, err)
	require.NotNil(t, d)
	assert.NotEmpty(t, d.LastUsedDeviceName)
}

func TestPrimaryTrainingDevice(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	d, err := c.PrimaryTrainingDevice(t.Context())
	require.NoError(t, err)
	require.NotNil(t, d)
	assert.NotEmpty(t, d.RegisteredDevices)
}

func TestDeviceSettings(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	devices, err := c.Devices(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	out, err := c.DeviceSettings(t.Context(), devices[0].DeviceID)
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
}

func TestDeviceSolarData(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	devices, err := c.Devices(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	out, err := c.DeviceSolarData(t.Context(), devices[0].DeviceID, "2026-01-01", "2026-01-01")
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
}

// The tests below need registered devices and sensors, so they are recorded
// against a real account.

func TestHistoricalDevices(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.HistoricalDevices(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotZero(t, out[0].DeviceID)
}

func TestSensors(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.Sensors(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].SensorType)
}

func TestUserDevice(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()
	ctx := t.Context()

	devices, err := c.Devices(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, devices)
	out, err := c.UserDevice(ctx, devices[0].DeviceID)
	require.NoError(t, err)
	assert.NotZero(t, out.LastUploadTimestamp)
}

func TestDeviceMessages(t *testing.T) {
	c, stop := newVCRClient(t)
	defer stop()

	out, err := c.DeviceMessages(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, out)
	assert.NotEmpty(t, out[0].MessageType)
}
