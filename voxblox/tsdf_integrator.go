package voxblox

import (
	"math"
	"sync"
	"time"

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
	Config TsdfConfig
	Layer  *TsdfLayer
}

func NewSimpleTsdfIntegrator(
	config TsdfConfig,
	layer *TsdfLayer,
) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		Config: config,
		Layer:  layer,
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

func (i *SimpleTsdfIntegrator) updateTsdfVoxel(
	origin Point,
	pointG Point,
	globalVoxelIndex IndexType,
	color Color,
	weight float64,
	voxel *TsdfVoxel,
) {
	voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, i.Layer.VoxelSize)
	sdf := computeDistance(origin, pointG, voxelCenter)

	updatedWeight := weight

	// Weight drop-off
	dropOffEpsilon := i.Layer.VoxelSize
	if i.Config.ConstWeight && sdf < -dropOffEpsilon {
		updatedWeight = weight * (i.Config.TruncationDistance + sdf) /
			(i.Config.TruncationDistance - dropOffEpsilon)
		updatedWeight = math.Max(updatedWeight, 0.0)
	}

	// TODO: Sparsity compensation

	// Calculate the new weight
	voxelWeight := voxel.getWeight()
	newWeight := voxelWeight + weight
	if newWeight < kEpsilon {
		return
	}
	newWeight = math.Min(newWeight, i.Config.MaxWeight)

	// Calculate the new distance
	newSdf := (sdf*updatedWeight + voxel.getDistance()*voxelWeight) / newWeight

	// Blend colors
	if math.Abs(sdf) < i.Config.TruncationDistance {
		newColor := blendTwoColors(voxel.getColor(), voxelWeight, color, weight)
		voxel.setColor(newColor)
	}

	var newDistance float64
	if sdf > 0 {
		newDistance = math.Min(i.Config.TruncationDistance, newSdf)
	} else {
		newDistance = math.Max(-i.Config.TruncationDistance, newSdf)
	}

	voxel.setDistance(newDistance)
	voxel.setWeight(newWeight)
}

func calculateWeight(pointC Point) float64 {
	distZ := math.Abs(pointC[2])
	if distZ > kEpsilon {
		return 1.0 / (distZ * distZ)
	}
	return 0.0
}

func (i *SimpleTsdfIntegrator) integratePoints(
	pose Transformation,
	points []Point,
	colors []Color,
	wg *sync.WaitGroup,
) {
	for j, point := range points {
		var ray Ray
		if validateRay(&ray, point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowClearing) {
			// Transform the point into the global frame.
			ray.Origin = pose.Position
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				&ray,
				i.Layer.VoxelSizeInv,
				i.Config.TruncationDistance,
				i.Config.MaxRange,
				i.Config.AllowCarving,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				block, voxel := getBlockAndVoxelFromGlobalVoxelIndex(i.Layer, globalVoxelIdx)
				weight := 1.0
				if !i.Config.ConstWeight {
					weight = calculateWeight(point)
				}
				i.updateTsdfVoxel(
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					colors[j],
					weight,
					voxel,
				)
				block.setUpdated()
			}
		}
	}
	wg.Done()
}

// integratePointCloud integrates a point cloud into the TSDF Layer.
func (i *SimpleTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer timeTrack(time.Now(), "integratePointCloud")

	// Fill color buffer with white if empty
	// TODO: This is a hack. Handle empty color slice downstream.
	if len(pointCloud.Colors) == 0 {
		pointCloud.Colors = make([]Color, len(pointCloud.Points))
		for i := range pointCloud.Colors {
			pointCloud.Colors[i] = ColorWhite
		}
	}

	// TODO: Organise points to minimize mutex contention.

	nThreads := i.Config.Threads
	wg := sync.WaitGroup{}
	nPointsPerThread := len(pointCloud.Points) / nThreads
	for threadIdx := 0; threadIdx < nThreads; threadIdx++ {
		startIdx := threadIdx * nPointsPerThread
		endIdx := (threadIdx + 1) * nPointsPerThread
		if threadIdx == nThreads-1 {
			endIdx = len(pointCloud.Points)
		}
		wg.Add(1)
		go i.integratePoints(
			pose,
			pointCloud.Points[startIdx:endIdx],
			pointCloud.Colors[startIdx:endIdx],
			&wg)

	}
	wg.Wait()
}
