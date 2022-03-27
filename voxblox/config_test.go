package voxblox

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigValid(t *testing.T) {
	config, err := ReadConfig("../voxblox.yaml")
	assert.Nil(t, err, "Error reading config")
	assert.Equal(t, 0.05, config.VoxelSize, "voxel size should be 0.05")
	assert.Equal(t, 16, config.VoxelsPerSide, "voxels per side should be 16")
	assert.Equal(t, 0.1, config.MinRange, "min range should be 0.1")
	assert.Equal(t, 5.0, config.MaxRange, "max range should be 5.0")
	assert.Equal(
		t,
		runtime.NumCPU(),
		config.Threads,
		"num workers should be equal to number of cores",
	)
}
