package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"testing"
)

func TestDistanceToCylinder(t *testing.T) {
	cylinder := Cylinder{Center: Point{0.0, 0.0, 0.0}, Height: 2.0, Radius: 6.0}
	if cylinder.DistanceToPoint(Point{10.0, 0.0, 0.0}) != 4.0 {
		t.Errorf("Incorrect distance to cylinder")
	}
	if cylinder.DistanceToPoint(Point{0.0, 10.0, 0.0}) != 4.0 {
		t.Errorf("Incorrect distance to cylinder")
	}
	distance := cylinder.DistanceToPoint(Point{10.0, 10.0, 0.0})
	if !almostEqual(distance, 8.14213562, 0.0) {
		t.Errorf("Incorrect distance to cylinder: %f", distance)
	}
	distance = cylinder.DistanceToPoint(Point{0.0, 0.0, 10.0})
	if distance != 9.0 {
		t.Errorf("Incorrect distance to cylinder: %f", distance)
	}
	distance = cylinder.DistanceToPoint(Point{0.0, 0.0, -10.0})
	if distance != 9.0 {
		t.Errorf("Incorrect distance to cylinder: %f", distance)
	}
}

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
	if !almostEqual(intersectPoint.X, 0.5, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Y, 0.5, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Z, 1.0, 0.0) {
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
	if !almostEqual(intersectPoint.X, -0.9, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Y, 0.4, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Z, -1.0, 0.0) {
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
	if !almostEqual(intersectPoint.X, 5.98030172, 0.05) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Y, -0.2, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.Z, 0.0, 0.0) {
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

func TestDistanceToPlane(t *testing.T) {
	plane := Plane{
		Normal: Point{0.0, 0.0, 1.0},
		Center: Point{0.0, 0.0, 0.0},
	}
	if plane.DistanceToPoint(Point{0.0, 0.0, 10.0}) != 10.0 {
		t.Errorf("Incorrect distance to plane")
	}
	if plane.DistanceToPoint(Point{0.0, 0.0, -10.0}) != -10.0 {
		t.Errorf("Incorrect distance to plane")
	}
	if plane.DistanceToPoint(Point{1.0, 1.0, 1.0}) != 1.0 {
		t.Errorf("Incorrect distance to plane")
	}
	if plane.DistanceToPoint(Point{1.0, 1.0, 1.0}) != 1.0 {
		t.Errorf("Incorrect distance to plane")
	}
	plane = Plane{
		Normal: Point{1.2, 8.6, 2.0},
		Center: Point{X: 0.0, Y: 0.0, Z: 0.0},
	}
	if !almostEqual(plane.DistanceToPoint(Point{2.2, 1.0, 10.0}), 3.505910087489654, 0.0) {
		t.Errorf("Incorrect distance to plane")
	}
}
