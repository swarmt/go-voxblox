package voxblox

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
)

func TestGetVoxelWeight(t *testing.T) {
	weight := calculateWeight(Point{0.714538097, -2.8530097, -1.72378588})
	assert.InEpsilon(t, 0.336537421, weight, kEpsilon)

	weight = calculateWeight(Point{1.42907524, -5.14151907, -1.49416912})
	assert.InEpsilon(t, 0.447920054, weight, kEpsilon)

	weight = calculateWeight(Point{0.7455977528518183, 2.0326033809822617, -2.2139824305115448})
	assert.InEpsilon(t, 0.20401009577963028, weight, kEpsilon)
}

func TestComputeDistance(t *testing.T) {
	distance := computeDistance(
		Point{0, 6, 2},
		Point{2.2434782608695647, 5.254402247148182, 0},
		Point{2.55, 5.15, -0.25},
	)
	assert.InEpsilon(t, -0.4086756820479831, distance, kEpsilon)
}

func TestValidateRay(t *testing.T) {
	var ray Ray
	valid := validateRay(&ray, Point{3.42977858, -3.22447705, -1.68651485}, 0.01, 5.0, true)
	assert.True(t, valid)
	assert.True(t, ray.Clearing)
}

func TestRayCaster(t *testing.T) {
	point := Point{0.714538097, -2.8530097, -1.72378588}
	var ray Ray
	validateRay(&ray, point, 0.1, 5, true)
	pose := Transformation{
		Translation: Point{0, 6, 2},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}

	ray.Origin = pose.Translation
	ray.Point = pose.transformPoint(point)
	voxelSizeInv := 10.0
	truncationDistance := 0.4

	rayCaster := NewRayCaster(
		&ray,
		voxelSizeInv,
		truncationDistance,
		5,
		true,
		true)

	assert.Equal(t, IndexType{0, 60, 20}, rayCaster.currentIndex)
	assert.Equal(t, Point{0, 60, 20}, rayCaster.startScaled)
	assert.InEpsilon(t, -29.7955704, rayCaster.endScaled[0], kEpsilon)
	assert.InEpsilon(t, 52.0162201, rayCaster.endScaled[1], kEpsilon)
	assert.InEpsilon(t, -2.34668899, rayCaster.endScaled[2], kEpsilon)
	assert.Equal(t, IndexType{-1, -1, -1}, rayCaster.stepSigns)
	assert.InEpsilon(t, 0.0335620344, rayCaster.tStepSize[0], kEpsilon)
	assert.InEpsilon(t, 0.12525396, rayCaster.tStepSize[1], kEpsilon)
	assert.InEpsilon(t, 0.0447493568, rayCaster.tStepSize[2], kEpsilon)
	assert.Equal(t, 61, rayCaster.lengthInSteps)

	var globalVoxelIdx IndexType
	for rayCaster.nextRayIndex(&globalVoxelIdx) {
	}
	assert.Equal(t, 62, rayCaster.currentStep)

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
	assert.Equal(t, 0.0, rayCaster.startScaled[0])
	assert.Equal(t, 60.0, rayCaster.startScaled[1])
	assert.Equal(t, 20.0, rayCaster.startScaled[2])
	assert.InEpsilon(t, 27.9682636, rayCaster.endScaled[0], kEpsilon)
	assert.InEpsilon(t, 28.4457779, rayCaster.endScaled[1], kEpsilon)
	assert.InEpsilon(t, 1.5998435, rayCaster.endScaled[2], kEpsilon)

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
	assert.Equal(t, 0.0, rayCaster.startScaled[0])
	assert.Equal(t, 60.0, rayCaster.startScaled[1])
	assert.Equal(t, 20.0, rayCaster.startScaled[2])
	assert.InEpsilon(t, -40.6885185, rayCaster.endScaled[0], kEpsilon)
	assert.InEpsilon(t, 48.8891029, rayCaster.endScaled[1], kEpsilon)
	assert.InEpsilon(t, 1.59945011, rayCaster.endScaled[2], kEpsilon)
}

func TestUpdateTsdfVoxel(t *testing.T) {
	layer := NewTsdfLayer(config.VoxelSize, config.VoxelsPerSide)

	origin := Point{0.0, 6.0, 2.0}
	pointG := Point{1.31130219e-06, 5.2854619, 1.1920929e-07}
	weight := 0.252516776
	globalVoxelIndex := IndexType{0, 60, 20}
	_, voxel := getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, 0.4, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.252516776, voxel.getWeight(), kEpsilon)
	assert.Len(t, layer.blocks, 1)

	globalVoxelIndex = IndexType{0, 59, 20}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, 0.4, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.252516776, voxel.getWeight(), kEpsilon)
	assert.Len(t, layer.blocks, 1)

	globalVoxelIndex = IndexType{0, 54, 3}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, 0.384953856, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.252516776, voxel.getWeight(), kEpsilon)

	globalVoxelIndex = IndexType{0, 52, -1}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, -0.05901622175430088, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.252516776, voxel.getWeight(), kEpsilon)

	pointG = Point{-0.0166654587, 5.2854619, 1.1920929e-07}
	weight = 0.252939552
	globalVoxelIndex = IndexType{0, 60, 20}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, 0.4, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.505456328, voxel.getWeight(), kEpsilon)

	weight = 0.252939552
	globalVoxelIndex = IndexType{-1, 52, -3}
	_, voxel = getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)

	updateTsdfVoxel(
		layer,
		&config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)
	assert.InEpsilon(t, -0.247611046, voxel.getDistance(), kEpsilon)
	assert.InEpsilon(t, 0.128483981, voxel.getWeight(), kEpsilon)
}
