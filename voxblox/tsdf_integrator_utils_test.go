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
	pointC = Point{0.7455977528518183, 2.0326033809822617, -2.2139824305115448}
	weight = calculateWeight(pointC)
	if !almostEqual(weight, 0.20401009577963028, kEpsilon) {
		t.Errorf("Expected weight to be 0.20401009577963028, got %f", weight)
	}
}

func TestComputeDistance(t *testing.T) {
	origin := Point{0, 6, 2}
	point := Point{2.2434782608695647, 5.254402247148182, 0}
	voxelCenter := Point{2.55, 5.15, -0.25}
	distance := computeDistance(origin, point, voxelCenter)
	if !almostEqual(distance, -0.4086756820479831, kEpsilon) {
		t.Errorf("Expected distance to be -0.4086756820479831, got %f", distance)
	}
}

func TestValidateRay(t *testing.T) {
	var ray Ray
	valid := validateRay(&ray, Point{3.42977858, -3.22447705, -1.68651485}, 0.01, 5.0, true)
	if valid != true {
		t.Errorf("Ray should be valid")
	}
	if ray.Clearing != true {
		t.Errorf("Ray should be clearing")
	}
}

func TestRayCaster(t *testing.T) {
	point := Point{0.714538097, -2.8530097, -1.72378588}
	var ray Ray
	validateRay(&ray, point, 0.1, 5, true)
	pose := Transformation{
		Translation: Point{0, 6, 2},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	// Transform the point into the global frame.
	ray.Origin = pose.Translation
	ray.Point = pose.transformPoint(point)
	voxelSizeInv := 10.0
	truncationDistance := 0.4

	// Create a ray caster.
	rayCaster := NewRayCaster(
		&ray,
		voxelSizeInv,
		truncationDistance,
		5,
		true,
		true)
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

	ray = Ray{
		Origin:   Point{0.0, 6.0, 2.0},
		Point:    Point{-4.42253256, 4.79232979, 1.1920929e-07},
		Length:   5.00172567,
		Clearing: true,
	}

	rayCaster = NewRayCaster(
		&ray,
		10.0,
		0.4,
		5.0,
		true,
		true,
	)
	if !almostEqual(rayCaster.startScaled[0], 0.0, kEpsilon) ||
		!almostEqual(rayCaster.startScaled[1], 60.0, kEpsilon) ||
		!almostEqual(rayCaster.startScaled[2], 20.0, kEpsilon) {
		t.Errorf("Raycaster start scaled incorrect")
	}
	if !almostEqual(rayCaster.endScaled[0], -40.6885185, kEpsilon*2) ||
		!almostEqual(rayCaster.endScaled[1], 48.8891029, kEpsilon*2) ||
		!almostEqual(rayCaster.endScaled[2], 1.59945011, kEpsilon*2) {
		t.Errorf("Raycaster end scaled incorrect")
	}
}

func TestUpdateTsdfVoxel(t *testing.T) {
	layer := NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)

	origin := Point{0.0, 6.0, 2.0}
	pointG := Point{1.31130219e-06, 5.2854619, 1.1920929e-07}
	weight := 0.252516776
	globalVoxelIndex := IndexType{0, 60, 20}
	_, voxel := getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.252516776, kEpsilon) {
		t.Errorf("Expected weight to be 0.252516776, got %f", voxel.getWeight())
	}
	if len(layer.blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(layer.blocks))
	}

	globalVoxelIndex = IndexType{0, 59, 20}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.252516776, kEpsilon) {
		t.Errorf("Expected weight to be 0.252516776, got %f", voxel.getWeight())
	}
	if len(layer.blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(layer.blocks))
	}

	globalVoxelIndex = IndexType{0, 54, 3}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.384953856, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.252516776, kEpsilon) {
		t.Errorf("Expected weight to be 0.252516776, got %f", voxel.getWeight())
	}

	globalVoxelIndex = IndexType{0, 52, -1}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), -0.0590159893, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.252516776, kEpsilon) {
		t.Errorf("Expected weight to be 0.252516776, got %f", voxel.getWeight())
	}

	pointG = Point{-0.0166654587, 5.2854619, 1.1920929e-07}
	weight = 0.252939552
	globalVoxelIndex = IndexType{0, 60, 20}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.505456328, kEpsilon) {
		t.Errorf("Expected weight to be 0.505456328, got %f", voxel.getWeight())
	}

	weight = 0.252939552
	globalVoxelIndex = IndexType{-1, 52, -3}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		tsdfConfig,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), -0.247611046, kEpsilon) {
		t.Errorf("Expected Tsdf to be -0.247611046, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.128483981, kEpsilon) {
		t.Errorf("Expected weight to be 0.128483981, got %f", voxel.getWeight())
	}
}
