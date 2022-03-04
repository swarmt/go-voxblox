package voxblox

import "github.com/ungerik/go3d/float64/vec3"

type RayCaster struct {
	VoxelSizeInv       float64
	MaxDistance        float64
	TruncationDistance float64
}

func (r *RayCaster) Cast(
	origin vec3.T,
	point vec3.T,
	clearing bool,
	carving bool,
	fromOrigin bool,
) {

}

func NewRayCaster(
	voxelSizeInv float64,
	maxDistance float64,
	truncationDistance float64,
) *RayCaster {
	return &RayCaster{
		VoxelSizeInv:       voxelSizeInv,
		MaxDistance:        maxDistance,
		TruncationDistance: truncationDistance,
	}
}
