package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
)

type TsdfIntegrator interface {
	integratePointCloud(
		pose Transformation,
		pointCloud PointCloud,
		freeSpacePoints bool,
	)
}

type SimpleTsdfIntegrator struct {
	VoxelCarving       bool
	TruncationDistance float64
	MinDistance        float64
	MaxDistance        float64
	Layer              *TsdfLayer
}

func NewSimpleTsdfIntegrator(
	voxelCarving bool,
	truncationDistance float64,
	minDistance float64,
	maxDistance float64,
	layer *TsdfLayer,
) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		VoxelCarving:       voxelCarving,
		TruncationDistance: truncationDistance,
		MinDistance:        minDistance,
		MaxDistance:        maxDistance,
		Layer:              layer,
	}
}

func (i *SimpleTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
	freeSpacePoints bool,
) {
	// Create a new Ray-caster object.
	rayCaster := NewRayCaster(i.Layer.VoxelSizeInv, i.MaxDistance, i.TruncationDistance)
	_ = rayCaster // TODO

	// Integrate the point cloud.
	for _, point := range pointCloud.Points {
		if point.X == 0 && point.Y == 0 && point.Z == 0 {
			continue
		}
		origin := pose.Position
		pointG := pose.TransformPoint(point)

		_, _ = origin, pointG // TODO

	}
}

func isPointValid(i *SimpleTsdfIntegrator, point vec3.T, freeSpacePoint bool) (bool, bool) {
	rayDistance := point.Length()
	_ = rayDistance // TODO
	return false, false
}
