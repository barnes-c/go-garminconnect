package garminconnect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevices(t *testing.T) {
	c, stop := newVCRClient(t, "devices")
	defer stop()

	devices, err := c.Devices()
	require.NoError(t, err)
	assert.NotEmpty(t, devices)
}

func TestLastUsedDevice(t *testing.T) {
	c, stop := newVCRClient(t, "last_used_device")
	defer stop()

	d, err := c.LastUsedDevice()
	require.NoError(t, err)
	assert.NotEmpty(t, d)
}

func TestPrimaryTrainingDevice(t *testing.T) {
	c, stop := newVCRClient(t, "primary_training_device")
	defer stop()

	d, err := c.PrimaryTrainingDevice()
	require.NoError(t, err)
	assert.NotEmpty(t, d)
}
