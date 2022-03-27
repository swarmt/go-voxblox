package voxblox

type Config struct {
	// TSDF configuration.
	VoxelSize                   float64 `yaml:"voxel_size"`
	VoxelsPerSide               int     `yaml:"block_size"`
	MinRange                    float64 `yaml:"min_range"`
	MaxRange                    float64 `yaml:"max_range"`
	TruncationDistance          float64 `yaml:"truncation_distance"`
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
