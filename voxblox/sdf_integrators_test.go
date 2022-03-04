package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"math"
	"os"
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
			X: radius * math.Sin(angle),
			Y: radius * math.Cos(angle),
			Z: height,
		}
		facingDirection := subtractPoints(cylinder.Center, position)
		desiredYaw := 0.0
		if facingDirection.X > 1e-4 || facingDirection.Y > 1e-4 {
			desiredYaw = math.Atan2(facingDirection.Y, facingDirection.X)
		}
		qY := quaternion.FromYAxisAngle(-0.1)
		qZ := quaternion.FromZAxisAngle(desiredYaw)
		q := quaternion.Mul(&qY, &qZ)
		transform := Transformation{
			Position: *position.asVec3(),
			Rotation: q,
		}
		poses = append(poses, transform)
	}

	truncationDistance = voxelSize * 4.0
	cameraResolution = vec2.T{320, 240}
	fovHorizontal = 150.0
	minDistance = 0.1
	maxDistance = 10.0
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
		simpleLayer,
	)

	// TODO: Merged integrator

	// TODO: Fast integrator

	// Create a text file to store the results.
	file, _ := os.Create("pointcloud.txt")
	defer file.Close()

	// Iterate over all poses and integrate.
	for _, pose := range poses {
		pointCloud := world.getPointCloudFromTransform(
			&pose,
			cameraResolution,
			fovHorizontal,
			maxDistance,
		)
		for _, point := range pointCloud.Points {
			if point.X != 0.0 && point.Y != 0.0 && point.Z != 0.0 {
				simpleTsdfIntegrator.integratePointCloud(pose, pointCloud, false)
			}
		}
	}
}
