package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

func ReadConfig(filename string) (Config, error) {
	config := new(Config)

	// Read the file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return *config, err
	}

	// parse the bytes to yaml
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return *config, err
	}

	return *config, nil
}
