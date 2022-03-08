package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"math"
	"testing"
)

var (
	world              *SimulationWorld
	poses              []Transformation
	voxelCarving       bool
	voxelSize          float64
	voxelsPerSide      int
	truncationDistance float64
	cameraResolution   vec2.T
	fovHorizontal      float64
	minDistance        float64
	maxDistance        float64
	maxWeight          float64
)

func init() {
	// Create a test environment.
	// It consists of a 10x10x7 m environment with a cylinder in the middle.
	voxelSize = 0.1
	voxelsPerSide = 16
	minBound := Point{-5.0, -5.0, -1.0}
	maxBound := Point{5.0, 5.0, 6.0}
	world = NewSimulationWorld(voxelSize, minBound, maxBound)
	cylinder := Cylinder{
		Center: Point{0.0, 0.0, 2.0},
		Radius: 2.0,
		Height: 4.0,
	}
	world.Objects = append(world.Objects, &cylinder)
	//world.AddGroundLevel(0.0) // TODO: Add ground level.

	// Generate poses around the cylinder.
	radius := 6.0
	height := 2.0
	numPoses := 50
	maxAngle := 2.0 * math.Pi
	angleIncrement := maxAngle / float64(numPoses)
	poses = []Transformation{}
	for angle := 0.0; angle < maxAngle; angle += angleIncrement {
		position := Point{
			radius * math.Sin(angle),
			radius * math.Cos(angle),
			height,
		}
		facingDirection := vec3.Sub(&cylinder.Center, &position)
		desiredYaw := 0.0
		if facingDirection[0] > 1e-4 || facingDirection[1] > 1e-4 {
			desiredYaw = math.Atan2(facingDirection[1], facingDirection[0])
		}
		qY := quaternion.FromYAxisAngle(-0.1)
		qZ := quaternion.FromZAxisAngle(desiredYaw)
		q := quaternion.Mul(&qY, &qZ)
		transform := Transformation{
			Position: position,
			Rotation: q,
		}
		poses = append(poses, transform)
	}

	truncationDistance = voxelSize * 4.0
	cameraResolution = vec2.T{320, 240}
	fovHorizontal = 150.0
	minDistance = 0.1
	maxDistance = 10.0
	maxWeight = 1000.0
	voxelCarving = true

}

func TestTsdfIntegrators(t *testing.T) {
	// Simple integrator
	simpleLayer := NewTsdfLayer(voxelSize, voxelsPerSide)
	simpleTsdfIntegrator := NewSimpleTsdfIntegrator(
		voxelCarving,
		truncationDistance,
		minDistance,
		maxDistance,
		maxWeight,
		simpleLayer,
	)

	// TODO: Merged integrator

	// TODO: Fast integrator

	// Iterate over all poses and integrate.
	for _, pose := range poses {
		pointCloud := world.getPointCloudFromTransform(
			&pose,
			cameraResolution,
			fovHorizontal,
			maxDistance,
		)
		simpleTsdfIntegrator.integratePointCloud(pose, pointCloud)
	}

	// Check the number of blocks in the layers
	if len(simpleLayer.blocks) == 0 {
		t.Errorf("No blocks in simple layer")
	}

}

func TestGetVoxelWeight(t *testing.T) {
	pointC := Point{0.714538097, -2.8530097, -1.72378588}
	weight := getVoxelWeight(pointC, false)
	if !almostEqual(weight, 0.336537421, kEpsilon) {
		t.Errorf("Expected weight to be 0.336537421, got %f", weight)
	}
	weight = getVoxelWeight(pointC, true)
	if weight != 1.0 {
		t.Errorf("Expected weight to be 1.0, got %f", weight)
	}

}

func TestUpdateTsdfVoxel(t *testing.T) {
	layer := NewTsdfLayer(voxelSize, voxelsPerSide)
	origin := Point{0.0, 6.0, 2.0}
	pointC := Point{0.714538097, -2.8530097, -1.72378588}
	pointG := Point{-2.66666508, 5.2854619, 1.1920929e-07}
	globalVoxelIndex := IndexType{0, 60, 20}
	truncationDistance := 0.4
	maxWeight := 1000.0
	voxel := allocateStorageAndGetVoxelPtr(layer, globalVoxelIndex)
	weight := getVoxelWeight(pointC, false)

	updateTsdfVoxel(
		layer,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		truncationDistance,
		maxWeight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.336537421, kEpsilon) {
		t.Errorf("Expected weight to be 0.336537421, got %f", voxel.getWeight())
	}
	if len(layer.blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(layer.blocks))
	}

	// Update the voxel again.
	updateTsdfVoxel(
		layer,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		truncationDistance,
		maxWeight,
		voxel,
	)
	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.336537421*2, kEpsilon) {
		t.Errorf("Expected weight to be 0.336537421 * 2, got %f", voxel.getWeight())
	}
}
