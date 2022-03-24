package voxblox

import (
	"math"

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
	diffIndex := subIndex(r.endIndex, r.currentIndex)

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
		sgn(rayScaled[0]),
		sgn(rayScaled[1]),
		sgn(rayScaled[2]),
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

// NewRayCaster creates a new RayCaster.
func NewRayCaster(
	ray *Ray,
	voxelSizeInv float64,
	truncationDistance float64,
	maxRange float64,
	allowCarving bool,
	castFromOrigin bool,
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
		ray.Length = math.Min(math.Max(delta.Length()-truncationDistance, 0), maxRange)
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
	if castFromOrigin {
		rayCaster.SetUp(startScaled, endScaled)
	} else {
		rayCaster.SetUp(endScaled, startScaled)
	}

	return rayCaster
}

// validateRay checks if the ray is valid.
// Sets the clearing flag if longer than the max range.
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

// updateTsdfVoxel updates the voxel SDF and weight.
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
	if config.WeightDropOff && sdf < -dropOffEpsilon {
		updatedWeight = weight * (config.TruncationDistance + sdf) /
			(config.TruncationDistance - dropOffEpsilon)
		updatedWeight = math.Max(updatedWeight, 0.0)
	}

	// TODO: Sparsity compensation

	// Lock the mutex
	voxel.Lock()
	defer voxel.Unlock()

	// Calculate the new weight
	newWeight := voxel.weight + updatedWeight
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

// splitPointCloud splits a PointCloud in to a slice of smaller PointClouds
// divided by the chunk number.
func splitPointCloud(
	pointCloud *PointCloud,
	chunkCount int,
) []PointCloud {
	chunkSize := len(pointCloud.Points) / chunkCount
	chunks := make([]PointCloud, chunkCount)
	for i := 0; i < chunkCount; i++ {
		chunks[i] = PointCloud{
			Points: pointCloud.Points[i*chunkSize : (i+1)*chunkSize],
			Colors: pointCloud.Colors[i*chunkSize : (i+1)*chunkSize],
		}
	}
	return chunks
}
