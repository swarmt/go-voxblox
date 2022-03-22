package voxblox

import (
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
)

func TestGetPointCloudFromTransform(t *testing.T) {
	world := SimulationWorld{}
	cylinder := Cylinder{Center: Point{0.0, 0.0, 2.0}, Height: 4.0, Radius: 2.0}
	world.AddObject(&cylinder)
	plane := Plane{Center: Point{0.0, 0.0, 0.0}, Normal: vec3.T{0.0, 0.0, 1.0}}
	world.AddObject(&plane)
	transform := Transformation{
		Translation: vec3.T{0.0, 6.0, 2.0},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	cameraResolution := vec2.T{320, 240}
	fovHorizontal := 150.0
	maxDistance := 100.0
	pointCloud := world.getPointCloudFromTransform(
		&transform,
		cameraResolution,
		fovHorizontal,
		maxDistance,
	)
	if !almostEqual(pointCloud.Points[0][0], -2.66666627, 0.001) ||
		!almostEqual(pointCloud.Points[0][1], 5.28546286, 0.001) ||
		!almostEqual(pointCloud.Points[0][2], 0.0, 0.001) {
		t.Errorf("Incorrect point in point cloud: %v", pointCloud.Points[0])
	}
}
