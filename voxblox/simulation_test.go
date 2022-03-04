package voxblox

import (
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"testing"
)

func TestGetPointCloudFromViewpoint(t *testing.T) {
	world := SimulationWorld{}
	cylinder := Cylinder{Center: Point{0.0, 0.0, 0.0}, Height: 2.0, Radius: 1.0}
	world.AddObject(&cylinder)
	viewOrigin := vec3.T{0.0, 0.0, 10.0}
	viewDirection := vec3.T{0.0, 0.0, -1.0}
	cameraResolution := vec2.T{640, 480}
	fovHorizontal := 60.0
	maxDistance := 100.0
	pointCloud := world.GetPointCloudFromViewpoint(
		viewOrigin,
		viewDirection,
		cameraResolution,
		fovHorizontal,
		maxDistance,
	)
	if len(pointCloud.Points) != 640*480 {
		t.Errorf("Incorrect number of points in point cloud: %d", len(pointCloud.Points))
	}
}
