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
	assert.NotEmpty(t, d)
}

func TestPrimaryTrainingDevice(t *testing.T) {
	c, stop := newVCRClient(t, "primary_training_device")
	defer stop()

	d, err := c.PrimaryTrainingDevice(t.Context())
	require.NoError(t, err)
	assert.NotEmpty(t, d)
}
