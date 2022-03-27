package voxblox

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
		Translation: Point{0.0, 6.0, 2.0},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	pointCloud := world.getPointCloudFromTransform(
		&transform,
		vec2.T{320, 240},
		150.0,
		100.0,
	)
	assert.InEpsilon(t, -2.66666627, pointCloud.Points[0][0], 0.001)
	assert.InEpsilon(t, 5.28546286, pointCloud.Points[0][1], 0.001)
	assert.Equal(t, 0.0, pointCloud.Points[0][2])
}
