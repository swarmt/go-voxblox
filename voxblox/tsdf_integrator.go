package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"math"
	"math/rand"
	"sync"
	"time"
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

func (i *SimpleTsdfIntegrator) updateTsdfVoxel(
	origin Point,
	pointG Point,
	globalVoxelIndex IndexType,
	color Color,
	weight float64,
	voxel *TsdfVoxel,
) {
	voxelSize := i.layer.VoxelSize

	voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, voxelSize)
	sdf := computeDistance(origin, pointG, voxelCenter)

	updatedWeight := weight

	// Weight drop-off
	dropOffEpsilon := voxelSize
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

	// TODO: Color blending

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
	wg *sync.WaitGroup,
) {
	for _, point := range points {
		ray := validateRay(point, i.Config.MinRange, i.Config.MaxRange, i.Config.AllowClearing)

		if ray.Valid {
			// Transform the point into the global frame.
			ray.Origin = pose.Position
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				ray,
				i.layer.VoxelSizeInv,
				i.Config.TruncationDistance,
				i.Config.MaxRange,
				i.Config.AllowCarving,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				voxel := allocateStorageAndGetVoxelPtr(i.layer, globalVoxelIdx)
				weight := 1.0
				if !i.Config.ConstWeight {
					weight = calculateWeight(point)
				}
				// TODO: Voxel color
				i.updateTsdfVoxel(
					ray.Origin,
					ray.Point,
					globalVoxelIdx,
					Color{},
					weight,
					voxel,
				)
			}
		}
	}
	wg.Done()
}

// shufflePoints returns a shuffled slice of points.
func shufflePoints(points []Point) []Point {
	shuffled := make([]Point, len(points))
	perm := rand.Perm(len(points))
	for i, v := range perm {
		shuffled[v] = points[i]
	}
	return shuffled
}

// integratePointCloud integrates a point cloud into the TSDF layer.
func (i *SimpleTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer timeTrack(time.Now(), "integratePointCloud")

	// Shuffle points to minimize mutex contention.
	// TODO: better way to do this?
	points := shufflePoints(pointCloud.Points)

	nThreads := i.Config.IntegratorThreads
	wg := &sync.WaitGroup{}
	nPointsPerThread := len(points) / nThreads
	for threadIdx := 0; threadIdx < nThreads; threadIdx++ {
		startIdx := threadIdx * nPointsPerThread
		endIdx := (threadIdx + 1) * nPointsPerThread
		if threadIdx == nThreads-1 {
			endIdx = len(points)
		}
		wg.Add(1)
		go i.integratePoints(pose, points[startIdx:endIdx], wg)
	}
	wg.Wait()
}
