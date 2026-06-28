package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevices(t *testing.T) {
	c, stop := newVCRClient(t, "devices")
	defer stop()

	devices, err := c.Devices(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, devices)
}

func TestLastUsedDevice(t *testing.T) {
	c, stop := newVCRClient(t, "last_used_device")
	defer stop()

	d, err := c.LastUsedDevice(t.Context())
	require.NoError(t, err)
	require.NotNil(t, d)
	assert.NotEmpty(t, d.LastUsedDeviceName)
}

func TestPrimaryTrainingDevice(t *testing.T) {
	c, stop := newVCRClient(t, "primary_training_device")
	defer stop()

	d, err := c.PrimaryTrainingDevice(t.Context())
	require.NoError(t, err)
	require.NotNil(t, d)
	assert.NotEmpty(t, d.RegisteredDevices)
}

func TestDeviceSettings(t *testing.T) {
	c, stop := newVCRClient(t, "device_settings")
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
	c, stop := newVCRClient(t, "device_solar_data")
	defer stop()

	devices, err := c.Devices(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	out, err := c.DeviceSolarData(t.Context(), devices[0].DeviceID, "2026-01-01", "2026-01-01")
	skipAPIError(t, err)
	require.NoError(t, err)
	require.NotNil(t, out)
}
