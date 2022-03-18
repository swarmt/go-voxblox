package voxblox

import (
	"math"
	"sync"

	"github.com/ungerik/go3d/float64/vec3"
)

type Ray struct {
	Origin   vec3.T
	Point    vec3.T
	Length   float64
	Clearing bool
}

type RayCaster struct {
	TruncationDistance float64
	MaxDistance        float64
	AllowCarving       bool
	lengthInSteps      int
	stepSigns          IndexType
	currentIndex       IndexType
	endIndex           IndexType
	currentStep        int
	tToNextBoundary    Point
	tStepSize          Point
	startScaled        Point
	endScaled          Point
}

func calculateWeight(pointC Point) float64 {
	distZ := math.Abs(pointC[2])
	if distZ > kEpsilon {
		return 1.0 / (distZ * distZ)
	}
	return 0.0
}

func computeDistance(origin, pointG, voxelCenter Point) float64 {
	vVoxelOrigin := voxelCenter.Sub(&origin)
	vPointOrigin := pointG.Sub(&origin)
	distG := vPointOrigin.Length()
	distGv := vec3.Dot(vPointOrigin, vVoxelOrigin) / distG
	sdf := distG - distGv
	return sdf
}

func (r *RayCaster) SetUp(startScaled, endScaled Point) {
	r.startScaled = startScaled
	r.endScaled = endScaled
	r.currentIndex = getGridIndexFromScaledPoint(startScaled)
	r.endIndex = getGridIndexFromScaledPoint(endScaled)
	diffIndex := SubIndex(r.endIndex, r.currentIndex)

	r.currentStep = 0

	r.lengthInSteps = int(
		math.Abs(
			float64(diffIndex[0]),
		) + math.Abs(
			float64(diffIndex[1]),
		) + math.Abs(
			float64(diffIndex[2]),
		))

	rayScaled := endScaled.Sub(&startScaled)

	// Get the signs of the components of the scaled ray.
	r.stepSigns = IndexType{
		Sgn(rayScaled[0]),
		Sgn(rayScaled[1]),
		Sgn(rayScaled[2]),
	}
	correctedStep := Point{
		math.Max(0, float64(r.stepSigns[0])),
		math.Max(0, float64(r.stepSigns[1])),
		math.Max(0, float64(r.stepSigns[2])),
	}
	currentIndexPoint := IndexToPoint(r.currentIndex)
	startScaledShifted := startScaled.Sub(&currentIndexPoint)

	distanceToBoundaries := correctedStep.Sub(startScaledShifted)

	// Voxblox: (std::abs(ray_scaled.x()) < 0.0) is never true?
	r.tToNextBoundary = Point{
		distanceToBoundaries[0] / rayScaled[0],
		distanceToBoundaries[1] / rayScaled[1],
		distanceToBoundaries[2] / rayScaled[2],
	}

	r.tStepSize = Point{
		float64(r.stepSigns[0]) / rayScaled[0],
		float64(r.stepSigns[1]) / rayScaled[1],
		float64(r.stepSigns[2]) / rayScaled[2],
	}
}

func PointMinCoeff(p *Point) int {
	min := math.Inf(1)
	minIndex := 0
	for i := 0; i < 3; i++ {
		if p[i] < min {
			min = p[i]
			minIndex = i
		}
	}
	return minIndex
}

func (r *RayCaster) nextRayIndex(rayIndex *IndexType) bool {
	r.currentStep++
	if r.currentStep > r.lengthInSteps {
		return false
	}
	*rayIndex = r.currentIndex
	tMinIdx := PointMinCoeff(&r.tToNextBoundary)

	r.currentIndex[tMinIdx] += r.stepSigns[tMinIdx]
	r.tToNextBoundary[tMinIdx] += r.tStepSize[tMinIdx]

	return true
}

func NewRayCaster(
	ray *Ray,
	voxelSizeInv float64,
	truncationDistance float64,
	maxRange float64,
	allowCarving bool,
) *RayCaster {
	rayCaster := &RayCaster{
		TruncationDistance: truncationDistance,
		MaxDistance:        maxRange,
		AllowCarving:       allowCarving,
	}

	unitRay := vec3.Sub(&ray.Point, &ray.Origin)
	unitRay.Normalize()

	var rayStart, rayEnd Point
	if ray.Clearing {
		delta := vec3.Sub(&ray.Point, &ray.Origin)
		ray.Length = delta.Length()
		ray.Length = math.Min(math.Max(ray.Length-truncationDistance, 0), maxRange)
		rayEnd = vec3.Add(&ray.Origin, unitRay.Scale(ray.Length))
		rayStart = ray.Origin
		if !allowCarving {
			rayStart = rayEnd
		}
	} else {
		rayEnd = vec3.Add(&ray.Point, unitRay.Scale(truncationDistance))
		rayStart = ray.Origin
		if !allowCarving {
			rayStart = vec3.Sub(&ray.Point, &unitRay)
		}
	}

	// Scale the ray to the voxel size.
	startScaled := rayStart.Scaled(voxelSizeInv)
	endScaled := rayEnd.Scaled(voxelSizeInv)

	// Set up the ray caster.
	rayCaster.SetUp(startScaled, endScaled)

	return rayCaster
}

func validateRay(
	ray *Ray,
	point Point,
	minLength float64,
	maxLength float64,
	allowClearing bool,
) bool {
	ray.Clearing = false
	// Faster than checking the ray length for 0,0,0 points.
	if point[0] == 0 && point[1] == 0 && point[2] == 0 {
		return false
	}
	ray.Length = point.Length()
	if ray.Length < minLength {
		return false
	} else if ray.Length > maxLength {
		if allowClearing {
			ray.Clearing = true
			return true
		} else {
			return true
		}
	}
	return true
}

func updateTsdfVoxel(
	layer *TsdfLayer,
	config TsdfConfig,
	origin Point,
	pointG Point,
	globalVoxelIndex IndexType,
	color Color,
	weight float64,
	voxel *TsdfVoxel,
) {
	voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, layer.VoxelSize)
	sdf := computeDistance(origin, pointG, voxelCenter)

	updatedWeight := weight

	// Weight drop-off
	dropOffEpsilon := layer.VoxelSize
	if config.ConstWeight && sdf < -dropOffEpsilon {
		updatedWeight = weight * (config.TruncationDistance + sdf) /
			(config.TruncationDistance - dropOffEpsilon)
		updatedWeight = math.Max(updatedWeight, 0.0)
	}

	// TODO: Sparsity compensation

	// Lock the mutex
	voxel.Lock()
	defer voxel.Unlock()

	// Calculate the new weight
	newWeight := voxel.weight + weight
	if newWeight < kEpsilon {
		return
	}
	newWeight = math.Min(newWeight, config.MaxWeight)
	voxel.weight = newWeight

	// Calculate the new distance
	newSdf := (sdf*updatedWeight + voxel.distance*voxel.weight) / newWeight

	// Blend colors
	if math.Abs(sdf) < config.TruncationDistance {
		newColor := blendTwoColors(voxel.color, voxel.weight, color, weight)
		voxel.color = newColor
	}

	var newDistance float64
	if sdf > 0 {
		newDistance = math.Min(config.TruncationDistance, newSdf)
	} else {
		newDistance = math.Max(-config.TruncationDistance, newSdf)
	}
	voxel.distance = newDistance
}

func integratePoints(
	layer *TsdfLayer,
	config TsdfConfig,
	pose Transformation,
	points []Point,
	colors []Color,
	wg *sync.WaitGroup,
) {
	for j, point := range points {
		var ray Ray
		if validateRay(&ray, point, config.MinRange, config.MaxRange, config.AllowClearing) {
			// Transform the point into the global frame.
			ray.Origin = pose.Position
			ray.Point = pose.transformPoint(point)

			// Create a new Ray-caster.
			rayCaster := NewRayCaster(
				&ray,
				layer.VoxelSizeInv,
				config.TruncationDistance,
				config.MaxRange,
				config.AllowCarving,
			)
			var globalVoxelIdx IndexType
			for rayCaster.nextRayIndex(&globalVoxelIdx) {
				block, voxel := getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIdx)
				weight := 1.0
				if !config.ConstWeight {
					weight = calculateWeight(point)
				}
				updateTsdfVoxel(
					layer,
					config,
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

func integratePointsParallel(
	layer *TsdfLayer,
	config TsdfConfig,
	pose Transformation,
	pointCloud PointCloud,
) {
	// Fill color buffer with white if empty
	// TODO: This is a hack. Handle empty color slice downstream.
	if len(pointCloud.Colors) == 0 {
		pointCloud.Colors = make([]Color, len(pointCloud.Points))
		for i := range pointCloud.Colors {
			pointCloud.Colors[i] = ColorWhite
		}
	}

	nThreads := config.Threads
	wg := sync.WaitGroup{}
	nPointsPerThread := len(pointCloud.Points) / nThreads
	for threadIdx := 0; threadIdx < nThreads; threadIdx++ {
		startIdx := threadIdx * nPointsPerThread
		endIdx := (threadIdx + 1) * nPointsPerThread
		if threadIdx == nThreads-1 {
			endIdx = len(pointCloud.Points)
		}
		wg.Add(1)
		go integratePoints(
			layer,
			config,
			pose,
			pointCloud.Points[startIdx:endIdx],
			pointCloud.Colors[startIdx:endIdx],
			&wg)

	}
	wg.Wait()
}
