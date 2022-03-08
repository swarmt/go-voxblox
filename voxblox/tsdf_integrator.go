package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"math"
	"runtime"
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
	MaxWeight          float64
	Layer              *TsdfLayer
}

func NewSimpleTsdfIntegrator(
	voxelCarving bool,
	truncationDistance float64,
	minDistance float64,
	maxDistance float64,
	maxWeight float64,
	layer *TsdfLayer,
) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		VoxelCarving:       voxelCarving,
		TruncationDistance: truncationDistance,
		MinDistance:        minDistance,
		MaxDistance:        maxDistance,
		MaxWeight:          maxWeight,
		Layer:              layer,
	}
}

func computeDistance(origin Point, pointG Point, voxelCenter Point) float64 {
	vVoxelOrigin := voxelCenter.Sub(&origin)
	vPointOrigin := pointG.Sub(&origin)

	distG := vPointOrigin.Length()
	distGv := vec3.Dot(vPointOrigin, vVoxelOrigin) / distG
	sdf := distG - distGv
	return sdf
}

func updateTsdfVoxel(
	layer *TsdfLayer,
	origin Point,
	pointG Point,
	globalVoxelIndex IndexType,
	color Color,
	weight float64,
	truncationDistance float64,
	maxWeight float64,
	voxel *TsdfVoxel,
) {
	voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, layer.getVoxelSize())
	sdf := computeDistance(origin, pointG, voxelCenter)

	updatedWeight := weight

	// TODO: Weight drop off

	// TODO: Sparsity compensation

	// Calculate the new weight
	voxelWeight := voxel.getWeight()
	newWeight := voxelWeight + weight
	if newWeight < kEpsilon {
		return
	}
	newWeight = math.Min(newWeight, maxWeight)

	// Calculate the new distance
	newSdf := (sdf*updatedWeight + voxel.getDistance()*voxelWeight) / newWeight

	// TODO: Color blending

	var newDistance float64
	if sdf > 0 {
		newDistance = math.Min(truncationDistance, newSdf)
	} else {
		newDistance = math.Max(-truncationDistance, newSdf)
	}

	voxel.setDistance(newDistance)
	voxel.setWeight(newWeight)
}

func getVoxelWeight(pointC Point, useConstWeight bool) float64 {
	if useConstWeight {
		return 1.0
	}
	distZ := math.Abs(pointC[2])
	if distZ > kEpsilon {
		return 1.0 / (distZ * distZ)
	}
	return 0.0
}

func (i *SimpleTsdfIntegrator) integratePoint(pose Transformation, points []Point) {
	for _, point := range points {
		ray := validateRay(point, i.MinDistance, i.MaxDistance, i.VoxelCarving)
		if ray.Valid {
			//Transform the point into the global frame.
			ray.Origin = pose.Position
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			// TODO: Allow clearing from config
			rayCaster := NewRayCaster(
				ray,
				i.Layer.getVoxelSizeInv(),
				i.TruncationDistance,
				i.MaxDistance,
				true,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				voxel := allocateStorageAndGetVoxelPtr(i.Layer, globalVoxelIdx)
				// TODO: weight drop off in config
				weight := getVoxelWeight(point, false)
				// TODO: Voxel color
				updateTsdfVoxel(
					i.Layer,
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					Color{},
					weight,
					i.TruncationDistance,
					i.MaxWeight,
					voxel,
				)
			}
		}
	}
}

func (i *SimpleTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	// Integrate the point cloud points with multiple cores
	// TODO: Configurable number of cores
	numCores := runtime.NumCPU()
	numPointsPerCore := len(pointCloud.Points) / numCores
	for coreIdx := 0; coreIdx < numCores; coreIdx++ {
		startIdx := coreIdx * numPointsPerCore
		endIdx := (coreIdx + 1) * numPointsPerCore
		if coreIdx == numCores-1 {
			endIdx = len(pointCloud.Points)
		}
		go i.integratePoint(pose, pointCloud.Points[startIdx:endIdx])
	}

}
