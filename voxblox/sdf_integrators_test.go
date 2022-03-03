package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"math"
	"testing"
)

var world *SimulationWorld
var poses []*Transformation
var voxelSize float64
var voxelsPerSide int
var truncationDistance float64

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
	world.AddGroundLevel(0.0)

	// Generate poses around the cylinder.
	radius := 6.0
	height := 2.0
	numPoses := 50
	maxAngle := 2.0 * math.Pi
	angleIncrement := maxAngle / float64(numPoses)
	poses = []*Transformation{}
	for angle := 0.0; angle < maxAngle; angle += angleIncrement {
		position := Point{
			x: radius * math.Sin(angle),
			y: radius * math.Cos(angle),
			z: height,
		}
		facingDirection := subtractPoints(cylinder.Center, position)
		desiredYaw := 0.0
		if facingDirection.x > 1e-4 || facingDirection.y > 1e-4 {
			desiredYaw = math.Atan2(facingDirection.y, facingDirection.x)
		}
		qY := quaternion.FromYAxisAngle(-0.1)
		qZ := quaternion.FromZAxisAngle(desiredYaw)
		q := quaternion.Mul(&qY, &qZ)
		transform := Transformation{
			Position: *position.asVec3(),
			Rotation: q,
		}
		poses = append(poses, &transform)
	}

	truncationDistance = voxelSize * 4.0

}

func TestTsdfIntegrators(t *testing.T) {
	// Simple integrator
	//simpleLayer := NewTsdfLayer(voxelSize, voxelsPerSide)
	//simpleTsdfIntegrator := NewSimpleTsdfIntegrator(truncationDistance, simpleLayer)

	// TODO: Merged integrator

	// TODO: Fast integrator

}
