package voxblox

import (
	"math"
	"runtime"

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
	Config Config
	layer  *TsdfLayer
}

func NewSimpleTsdfIntegrator(
	config Config,
	layer *TsdfLayer,
) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		Config: config,
		layer:  layer,
	}
}

func computeDistance(origin, pointG, voxelCenter Point) float64 {
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

func getVoxelWeight(pointC Point) float64 {
	distZ := math.Abs(pointC[2])
	if distZ > kEpsilon {
		return 1.0 / (distZ * distZ)
	}
	return 0.0
}

func (i *SimpleTsdfIntegrator) integratePoint(pose Transformation, points []Point) {
	for _, point := range points {
		ray := validateRay(point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowCarving)
		if ray.Valid {
			//Transform the point into the global frame.
			ray.Origin = pose.Position
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				ray,
				i.layer.getVoxelSizeInv(),
				i.Config.TruncationDistance,
				i.Config.MaxRange,
				i.Config.AllowClearing,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				voxel := allocateStorageAndGetVoxelPtr(i.layer, globalVoxelIdx)
				// TODO: weight drop off in Config
				weight := 1.0
				if !i.Config.ConstWeight {
					weight = getVoxelWeight(point)
				}
				// TODO: Voxel color
				updateTsdfVoxel(
					i.layer,
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					Color{},
					weight,
					i.Config.TruncationDistance,
					i.Config.MaxWeight,
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
