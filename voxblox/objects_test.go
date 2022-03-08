package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"testing"
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
