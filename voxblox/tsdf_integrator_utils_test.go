package voxblox

import (
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
)

func TestGetVoxelWeight(t *testing.T) {
	pointC := Point{0.714538097, -2.8530097, -1.72378588}
	weight := calculateWeight(pointC)
	if !almostEqual(weight, 0.336537421, kEpsilon) {
		t.Errorf("Expected weight to be 0.336537421, got %f", weight)
	}
	pointC = Point{1.42907524, -5.14151907, -1.49416912}
	weight = calculateWeight(pointC)
	if !almostEqual(weight, 0.447920054, kEpsilon) {
		t.Errorf("Expected weight to be 0.447920054, got %f", weight)
	}
}

func TestValidateRay(t *testing.T) {
	var ray Ray
	valid := validateRay(&ray, Point{0, 0, 0}, 1, 15, true)
	if valid == true {
		t.Errorf("Ray should not be valid")
	}

	if validateRay(&ray, Point{0, 0, 10}, 1, 15, true) == false {
		t.Errorf("Ray should be valid")
	}
	if ray.Length != 10.0 {
		t.Errorf("Ray length should be 10.0, got %f", ray.Length)
	}
	if ray.Clearing != false {
		t.Errorf("Ray should not be clearing")
	}

	if validateRay(&ray, Point{0, 0, 10}, 1, 8, true) == false {
		t.Errorf("Ray should be valid")
	}
	if ray.Clearing != true {
		t.Errorf("Ray should be clearing")
	}

	if validateRay(&ray, Point{0, 0, 10}, 1, 8, false) == false {
		t.Errorf("Ray should be valid")
	}
	if ray.Clearing == true {
		t.Errorf("Ray should not be clearing")
	}
	if validateRay(&ray, Point{0.714538097, -2.8530097, -1.72378588}, 0.1, 5, true) == false {
		t.Errorf("Ray should be valid")
	}
	if ray.Clearing == true {
		t.Errorf("Ray should not be clearing")
	}
}

func TestRayCaster(t *testing.T) {
	point := Point{0.714538097, -2.8530097, -1.72378588}
	var ray Ray
	validateRay(&ray, point, 0.1, 5, true)
	pose := Transformation{
		Position: Point{0, 6, 2},
		Rotation: quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	// Transform the point into the global frame.
	ray.Origin = pose.Position
	ray.Point = pose.transformPoint(point)
	voxelSizeInv := 10.0
	truncationDistance := 0.4

	// Create a ray caster.
	rayCaster := NewRayCaster(&ray, voxelSizeInv, truncationDistance, 5, true)
	if rayCaster.currentIndex[0] != 0 ||
		rayCaster.currentIndex[1] != 60 ||
		rayCaster.currentIndex[2] != 20 {
		t.Errorf(
			"RayCaster should start at (0,60,20), got (%d,%d,%d)",
			rayCaster.currentIndex[0],
			rayCaster.currentIndex[1],
			rayCaster.currentIndex[2],
		)
	}

	// Check the start scaled value
	if rayCaster.startScaled[0] != 0.0 ||
		rayCaster.startScaled[1] != 60.0 ||
		rayCaster.startScaled[2] != 20.0 {
		t.Errorf(
			"RayCaster startScaled should be (0,60,20), got (%f,%f,%f)",
			rayCaster.startScaled[0],
			rayCaster.startScaled[1],
			rayCaster.startScaled[2],
		)
	}

	// Check the end scaled value
	if !almostEqual(rayCaster.endScaled[0], -29.7955704, kEpsilon) ||
		!almostEqual(rayCaster.endScaled[1], 52.0162201, kEpsilon) ||
		!almostEqual(rayCaster.endScaled[2], -2.34668899, kEpsilon) {
		t.Errorf(
			"RayCaster endScaled should be (-29.7955704,52.0162201,-2.34668899), got (%f,%f,%f)",
			rayCaster.endScaled[0],
			rayCaster.endScaled[1],
			rayCaster.endScaled[2],
		)
	}

	// Step signs
	if rayCaster.stepSigns[0] != -1 ||
		rayCaster.stepSigns[1] != -1 ||
		rayCaster.stepSigns[2] != -1 {
		t.Errorf(
			"RayCaster stepSign should be (-1,-1,-1), got (%d,%d,%d)",
			rayCaster.stepSigns[0],
			rayCaster.stepSigns[1],
			rayCaster.stepSigns[2],
		)
	}

	// Step size
	if !almostEqual(rayCaster.tStepSize[0], 0.0335620344, kEpsilon) ||
		!almostEqual(rayCaster.tStepSize[1], 0.12525396, kEpsilon) ||
		!almostEqual(rayCaster.tStepSize[2], 0.0447493568, kEpsilon) {
		t.Errorf(
			"RayCaster stepSize should be (0.0335620344,0.12525396,0.0447493568), got (%f,%f,%f)",
			rayCaster.tStepSize[0],
			rayCaster.tStepSize[1],
			rayCaster.tStepSize[2],
		)
	}

	// Ray length in steps
	if rayCaster.lengthInSteps != 61 {
		t.Errorf(
			"RayCaster rayLengthInSteps should be 61, got %d",
			rayCaster.lengthInSteps,
		)
	}

	var globalVoxelIdx IndexType
	for rayCaster.nextRayIndex(&globalVoxelIdx) {
	}

	if rayCaster.currentStep != 62 {
		t.Errorf("Raycaster current step should be 62")
	}

	ray = Ray{
		Origin:   Point{0.0, 6.0, 2.0},
		Point:    Point{3.04000235, 2.57022285, 2.38418579e-07},
		Length:   4.60049868,
		Clearing: true,
	}

	rayCaster = NewRayCaster(
		&ray,
		10.0,
		0.4,
		5.0,
		true,
	)
	if !almostEqual(rayCaster.startScaled[0], 0.0, kEpsilon) ||
		!almostEqual(rayCaster.startScaled[1], 60.0, kEpsilon) ||
		!almostEqual(rayCaster.startScaled[2], 20.0, kEpsilon) {
		t.Errorf("Raycaster start scaled incorrect")
	}
	if !almostEqual(rayCaster.endScaled[0], 27.9682636, kEpsilon) ||
		!almostEqual(rayCaster.endScaled[1], 28.4457779, kEpsilon) ||
		!almostEqual(rayCaster.endScaled[2], 1.5998435, kEpsilon) {
		t.Errorf("Raycaster start scaled incorrect")
	}
}
