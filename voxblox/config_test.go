package voxblox

import "testing"

func TestReadConfig(t *testing.T) {
	config, err := ReadConfig("../voxblox.yaml")
	if err != nil {
		t.Errorf("Error reading config: %v", err)
	}
	if config.VoxelSize != 0.05 {
		t.Errorf("Voxel size should be 0.05, but is %f", config.VoxelSize)
	}
	if config.VoxelsPerSide != 16 {
		t.Errorf("Voxels per side should be 16, but is %d", config.VoxelsPerSide)
	}
	if config.MinRange != 0.1 {
		t.Errorf("Min range should be 0.1, but is %f", config.MinRange)
	}
	if config.MaxRange != 5.0 {
		t.Errorf("Max range should be 5.0, but is %f", config.MaxRange)
	}

}
