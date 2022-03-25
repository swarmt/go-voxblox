package voxblox

import (
	"sync"
	"time"

	"github.com/dustinxie/lockfree"
)

type TsdfIntegrator interface {
	IntegratePointCloud(pose Transformation, cloud PointCloud)
}

type SimpleTsdfIntegrator struct {
	Config TsdfConfig
	Layer  *TsdfLayer
}

// NewSimpleTsdfIntegrator creates a new SimpleTsdfIntegrator.
func NewSimpleTsdfIntegrator(config TsdfConfig, layer *TsdfLayer) *SimpleTsdfIntegrator {
	config.TruncationDistance = config.VoxelSize * 4
	return &SimpleTsdfIntegrator{
		Config: config,
		Layer:  layer,
	}
}

// IntegratePointCloud integrates a point cloud into the TSDF Layer.
func (i *SimpleTsdfIntegrator) IntegratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer TimeTrack(time.Now(), "Integrate Simple")

	wg := sync.WaitGroup{}
	for _, pC := range splitPointCloud(&pointCloud, i.Config.Threads) {
		wg.Add(1)
		go i.integratePoints(pose, pC, &wg)
	}
	wg.Wait()
}

func (i *SimpleTsdfIntegrator) integratePoints(
	pose Transformation,
	pointCloud PointCloud,
	wg *sync.WaitGroup,
) {
	for j, point := range pointCloud.Points {
		var ray Ray
		if validateRay(&ray, point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowClearing) {
			// Transform the point into the global frame.
			ray.Origin = pose.Translation
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				&ray,
				i.Layer.VoxelSizeInv,
				i.Config.TruncationDistance,
				i.Config.MaxRange,
				i.Config.AllowCarving,
				true,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				block, voxel := getBlockAndVoxelFromGlobalVoxelIndex(i.Layer, globalVoxelIdx)
				weight := 1.0
				if !i.Config.WeightConstant {
					weight = calculateWeight(point)
				}
				updateTsdfVoxel(
					i.Layer,
					i.Config,
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					pointCloud.Colors[j],
					weight,
					voxel,
				)
				block.setUpdated()
			}
		}
	}
	wg.Done()
}

type MergedTsdfIntegrator struct {
	Config TsdfConfig
	Layer  *TsdfLayer
}

// NewMergedTsdfIntegrator creates a new MergedTsdfIntegrator.
func NewMergedTsdfIntegrator(config TsdfConfig, layer *TsdfLayer) *MergedTsdfIntegrator {
	config.TruncationDistance = config.VoxelSize * 4
	return &MergedTsdfIntegrator{
		Config: config,
		Layer:  layer,
	}
}

// bundleRays decimates the point cloud by the voxel size.
func bundleRays(voxelSizeInv float64, pointCloud PointCloud) map[IndexType]int {
	voxelMap := make(map[IndexType]int)
	for j, point := range pointCloud.Points {
		voxelIndex := getGridIndexFromPoint(point, voxelSizeInv)
		voxelMap[voxelIndex] = j
	}
	return voxelMap
}

// IntegratePointCloud integrates a point cloud into the TSDF Layer.
func (i *MergedTsdfIntegrator) IntegratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer TimeTrack(time.Now(), "Integrate Merged")

	voxelMap := bundleRays(i.Layer.VoxelSizeInv, pointCloud)

	// Filter the point cloud by the voxel map
	filteredPoints := make([]Point, 0, len(pointCloud.Points))
	filteredColors := make([]Color, 0, len(pointCloud.Colors))
	for _, pointIndex := range voxelMap {
		filteredPoints = append(
			filteredPoints,
			pointCloud.Points[pointIndex],
		)
		filteredColors = append(
			filteredColors,
			pointCloud.Colors[pointIndex],
		)
	}
	pointCloud.Points = filteredPoints
	pointCloud.Colors = filteredColors

	wg := sync.WaitGroup{}
	for _, pC := range splitPointCloud(&pointCloud, i.Config.Threads) {
		wg.Add(1)
		go i.integratePoints(pose, pC, &wg)
	}
	wg.Wait()
}

func (i *MergedTsdfIntegrator) integratePoints(
	pose Transformation,
	pointCloud PointCloud,
	wg *sync.WaitGroup,
) {
	for j, point := range pointCloud.Points {
		var ray Ray
		if validateRay(&ray, point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowClearing) {
			// Transform the point into the global frame.
			ray.Origin = pose.Translation
			ray.Point = pose.transformPoint(point)

			// TODO: Merge weights

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				&ray,
				i.Layer.VoxelSizeInv,
				i.Config.TruncationDistance,
				i.Config.MaxRange,
				i.Config.AllowCarving,
				true,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				block, voxel := getBlockAndVoxelFromGlobalVoxelIndex(i.Layer, globalVoxelIdx)
				weight := 1.0
				if !i.Config.WeightConstant {
					weight = calculateWeight(point)
				}
				updateTsdfVoxel(
					i.Layer,
					i.Config,
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					pointCloud.Colors[j],
					weight,
					voxel,
				)
				block.setUpdated()
			}
		}
	}
	wg.Done()
}

type FastTsdfIntegrator struct {
	Config                 TsdfConfig
	Layer                  *TsdfLayer
	startVoxelApproxSet    lockfree.HashMap
	voxelObservedApproxSet lockfree.HashMap
}

func NewFastTsdfIntegrator(config TsdfConfig, layer *TsdfLayer) *FastTsdfIntegrator {
	config.TruncationDistance = config.VoxelSize * 4
	return &FastTsdfIntegrator{
		Config:                 config,
		Layer:                  layer,
		startVoxelApproxSet:    lockfree.NewHashMap(),
		voxelObservedApproxSet: lockfree.NewHashMap(),
	}
}

// voxelInStartApproxSet returns true if the voxel is in the approximate set.
// Adds it to the approximate set if it is not already there.
func (i *FastTsdfIntegrator) voxelInStartApproxSet(voxelIndex IndexType) bool {
	voxelIndexString := indexToString(voxelIndex)
	if _, ok := i.startVoxelApproxSet.Get(voxelIndexString); ok {
		return true
	}
	i.startVoxelApproxSet.Set(voxelIndexString, struct{}{})
	return false
}

// clearStartVoxelApproxSet clears the approximate set.
func (i *FastTsdfIntegrator) clearStartVoxelApproxSet() {
	i.startVoxelApproxSet = lockfree.NewHashMap()
}

// voxelInObservedApproxSet returns true if the voxel is in the approximate set.
// Adds it to the approximate set if it is not already there.
func (i *FastTsdfIntegrator) voxelInObservedApproxSet(voxelIndex IndexType) bool {
	voxelIndexString := indexToString(voxelIndex)
	if _, ok := i.voxelObservedApproxSet.Get(voxelIndexString); ok {
		return true
	}
	i.voxelObservedApproxSet.Set(voxelIndexString, struct{}{})
	return false
}

// clearObservedVoxelApproxSet clears the approximate set.
func (i *FastTsdfIntegrator) clearObservedVoxelApproxSet() {
	i.voxelObservedApproxSet = lockfree.NewHashMap()
}

// IntegratePointCloud integrates a point cloud into the TSDF Layer.
func (i *FastTsdfIntegrator) IntegratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer TimeTrack(time.Now(), "Integrate Fast")

	resetCounter := 0
	resetCounter++
	if resetCounter >= i.Config.ClearChecksEveryNFrames {
		resetCounter = 0
		i.clearStartVoxelApproxSet()
		i.clearObservedVoxelApproxSet()
	}

	wg := sync.WaitGroup{}
	for _, pC := range splitPointCloud(&pointCloud, i.Config.Threads) {
		wg.Add(1)
		go i.integratePoints(pose, pC, &wg)
	}
	wg.Wait()
}

func (i *FastTsdfIntegrator) integratePoints(
	pose Transformation,
	pointCloud PointCloud,
	wg *sync.WaitGroup,
) {
	for j, point := range pointCloud.Points {
		var ray Ray
		if !validateRay(&ray, point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowClearing) {
			continue
		}

		// Transform the point into the global frame.
		ray.Origin = pose.Translation
		ray.Point = pose.transformPoint(point)

		// Checks to see if another ray in this scan has already started 'close'
		// to this location. If it has then we skip ray casting this point. We
		// measure if a start location is 'close' to another points by inserting
		// the point into a set of voxels. This voxel set has a resolution
		// start_voxel_subsampling_factor times higher than the voxel size.
		globalVoxelIndex := getGridIndexFromPoint(ray.Point, i.Config.StartVoxelSubsamplingFactor*i.Layer.VoxelSizeInv)

		// Continue if the voxel is already in the set.
		if i.voxelInStartApproxSet(globalVoxelIndex) {
			continue
		}

		// Create a new Ray-caster.
		rayCaster := NewRayCaster(
			&ray,
			i.Layer.VoxelSizeInv,
			i.Config.TruncationDistance,
			i.Config.MaxRange,
			i.Config.AllowCarving,
			false,
		)

		consecutiveRayCollisions := 0

		for rayCaster.nextRayIndex(&globalVoxelIndex) {
			// Check if the current voxel has been seen by any ray cast this scan.
			// If it has increment the consecutive_ray_collisions counter, otherwise
			// reset it. If the counter reaches a threshold we stop casting as the
			// ray is deemed to be contributing too little new information.
			if i.voxelInObservedApproxSet(globalVoxelIndex) {
				consecutiveRayCollisions++
			}
			if consecutiveRayCollisions >= i.Config.MaxConsecutiveRayCollisions {
				break
			}

			block, voxel := getBlockAndVoxelFromGlobalVoxelIndex(i.Layer, globalVoxelIndex)
			weight := 1.0
			if !i.Config.WeightConstant {
				weight = calculateWeight(point)
			}
			updateTsdfVoxel(
				i.Layer,
				i.Config,
				ray.Origin,
				ray.Point,
				globalVoxelIndex,
				pointCloud.Colors[j],
				weight,
				voxel,
			)
			block.setUpdated()
		}
	}
	wg.Done()
}
