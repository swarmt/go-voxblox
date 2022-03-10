package voxblox

import (
	"testing"

	"github.com/ungerik/go3d/float64/vec3"
)

func TestRayIntersectionCylinder(t *testing.T) {
	// Top-down view of cylinder
	cylinder := Cylinder{Center: Point{0.0, 0.0, 0.0}, Height: 2.0, Radius: 6.0}
	rayOrigin := vec3.T{0.5, 0.5, 10.0}
	rayDirection := vec3.T{0.0, 0.0, -1.0}
	maxDistance := 100.0
	intersects, intersectPoint, _ := cylinder.RayIntersection(
		rayOrigin,
		rayDirection,
		maxDistance,
	)
	if !intersects {
		t.Errorf("Ray should intersect cylinder")
	}
	if !almostEqual(intersectPoint[0], 0.5, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[1], 0.5, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[2], 1.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}

	// Bottom-up view of cylinder
	rayOrigin = vec3.T{-1.2, -0.2, -4.0}
	rayDirection = vec3.T{0.1, 0.2, 1.0}
	maxDistance = 100.0
	intersects, intersectPoint, _ = cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if !intersects {
		t.Errorf("Ray should intersect cylinder")
	}
	if !almostEqual(intersectPoint[0], -0.9, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[1], 0.4, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[2], -1.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}

	// Side view of cylinder
	rayOrigin = vec3.T{10, -0.2, 0.0}
	rayDirection = vec3.T{-1.0, 0.0, 0.0}
	maxDistance = 100.0
	intersects, intersectPoint, _ = cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if !intersects {
		t.Errorf("Ray should intersect cylinder")
	}
	if !almostEqual(intersectPoint[0], 5.98030172, 0.05) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[1], -0.2, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint[2], 0.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}

	cylinder = Cylinder{Center: Point{0.0, 0.0, 1.0}, Height: 2.0, Radius: 2.0}
	rayOrigin = vec3.T{0.0, 0.0, 10.0}
	rayDirection = vec3.T{-0.35112344158839154, -0.4681645887845222, -0.8108848540793833}
	maxDistance = 100.0
	intersects, intersectPoint, _ = cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if intersects {
		t.Errorf("Ray should not intersect cylinder")
	}
}

func TestRayIntersectionPlane(t *testing.T) {
	plane := Plane{
		Center: Point{0.0, 0.0, 0.0},
		Normal: vec3.T{0.0, 0.0, 1.0},
		Color:  Color{},
	}
	intersects, intersectPoint, intersectDistance := plane.RayIntersection(
		Point{0, 6, 2},
		Point{-0.782229722, -0.209599614, -0.586672366},
		10.0,
	)
	if !intersects {
		t.Errorf("Ray should intersect plane")
	}
	if !almostEqual(intersectPoint[0], -2.66666627, kEpsilon) ||
		!almostEqual(intersectPoint[1], 5.28546286, kEpsilon) ||
		!almostEqual(intersectPoint[2], 0.0, kEpsilon) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectDistance, 3.40905786, kEpsilon) {
		t.Errorf("Incorrect intersection distance: %v", intersectDistance)
	}
}
