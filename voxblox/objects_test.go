package voxblox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRayIntersectionCylinder(t *testing.T) {
	// Top-down view of cylinder
	cylinder := Cylinder{Center: Point{0.0, 0.0, 0.0}, Height: 2.0, Radius: 6.0}
	intersects, intersectPoint, _ := cylinder.RayIntersection(
		Point{0.5, 0.5, 10.0},
		Point{0.0, 0.0, -1.0},
		100.0,
	)
	assert.True(t, intersects)
	assert.Equal(t, Point{0.5, 0.5, 1.0}, intersectPoint)

	// Bottom-up view of cylinder
	intersects, intersectPoint, _ = cylinder.RayIntersection(
		Point{-1.2, -0.2, -4.0},
		Point{0.1, 0.2, 1.0},
		100.0,
	)
	assert.True(t, intersects)
	assert.InEpsilon(t, -0.9, intersectPoint[0], kEpsilon)
	assert.InEpsilon(t, 0.4, intersectPoint[1], kEpsilon)
	assert.InEpsilon(t, -1.0, intersectPoint[2], kEpsilon)

	// Side view of cylinder
	intersects, intersectPoint, _ = cylinder.RayIntersection(
		Point{10, -0.2, 0.0},
		Point{-1.0, 0.0, 0.0},
		100.0,
	)
	assert.True(t, intersects)
	assert.InEpsilon(t, 5.98030172, intersectPoint[0], 0.05)
	assert.InEpsilon(t, -0.2, intersectPoint[1], kEpsilon)
	assert.InDelta(t, 0.0, intersectPoint[2], kEpsilon)

	cylinder = Cylinder{Center: Point{0.0, 0.0, 1.0}, Height: 2.0, Radius: 2.0}
	intersects, intersectPoint, _ = cylinder.RayIntersection(
		Point{0.0, 0.0, 10.0},
		Point{-0.35112344158839154, -0.4681645887845222, -0.8108848540793833},
		100.0,
	)
	assert.False(t, intersects)
}

func TestRayIntersectionPlane(t *testing.T) {
	plane := Plane{
		Center: Point{0.0, 0.0, 0.0},
		Normal: Point{0.0, 0.0, 1.0},
		Color:  Color{},
	}
	intersects, intersectPoint, intersectDistance := plane.RayIntersection(
		Point{0, 6, 2},
		Point{-0.782229722, -0.209599614, -0.586672366},
		10.0,
	)
	assert.True(t, intersects)
	assert.InEpsilon(t, -2.66666627, intersectPoint[0], kEpsilon)
	assert.InEpsilon(t, 5.28546286, intersectPoint[1], kEpsilon)
	assert.InDelta(t, 0.0, intersectPoint[2], kEpsilon)
	assert.InEpsilon(t, 3.40905786, intersectDistance, kEpsilon)
}
