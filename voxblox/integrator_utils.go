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
