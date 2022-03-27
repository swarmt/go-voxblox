package voxblox

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/ungerik/go3d/float64/quaternion"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// ROS
	RosMaster        string       `yaml:"ros_master"`
	TopicPointCloud2 string       `yaml:"topic_pointcloud2"`
	TopicTransform   string       `yaml:"topic_transform"`
	Translation      Point        `yaml:"translation"`
	Rotation         quaternion.T `yaml:"rotation"`

	// TSDF configuration.
	VoxelSize                   float64 `yaml:"voxel_size"`
	VoxelsPerSide               int     `yaml:"voxels_per_side"`
	MinRange                    float64 `yaml:"min_range"`
	MaxRange                    float64 `yaml:"max_range"`
	truncationDistance          float64
	AllowCarving                bool    `yaml:"allow_carving"`
	AllowClearing               bool    `yaml:"allow_clearing"`
	MaxWeight                   float64 `yaml:"max_weight"`
	WeightConstant              bool    `yaml:"weight_constant"`
	WeightDropOff               bool    `yaml:"weight_dropoff"`
	StartVoxelSubsamplingFactor float64 `yaml:"start_voxel_subsampling_factor"`
	MaxConsecutiveRayCollisions int     `yaml:"max_consecutive_ray_collisions"`
	Threads                     int     `yaml:"integrator_threads"`

	// Mesh configuration.
	UseColor  bool    `yaml:"use_color"`
	MinWeight float64 `yaml:"min_weight"`
}

// ReadConfig reads a yaml config file and returns a Config struct.
func ReadConfig(filename string) (Config, error) {
	config := new(Config)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return *config, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return *config, err
	}

	if config.VoxelSize <= 0 {
		return *config, fmt.Errorf("voxel size must be positive")
	}

	if config.VoxelsPerSide <= 0 {
		return *config, fmt.Errorf("voxels per side must be positive")
	}

	if config.MinRange < 0 {
		return *config, fmt.Errorf("min range must be positive")
	}

	if config.MinRange > config.MaxRange {
		return *config, fmt.Errorf("max range must be greater than min range")
	}

	if config.MaxWeight <= 0 {
		return *config, fmt.Errorf("max weight must be positive")
	}

	if config.StartVoxelSubsamplingFactor < 1.0 {
		return *config, fmt.Errorf("start voxel subsampling factor must be 1.0 or greater")
	}

	if config.MaxConsecutiveRayCollisions <= 0 {
		return *config, fmt.Errorf("max consecutive ray collisions must be positive")
	}

	if config.MinWeight < 0 {
		return *config, fmt.Errorf("min weight must be positive")
	}

	if config.Threads <= 0 {
		config.Threads = runtime.NumCPU()
	}

	return *config, nil
}
