package voxblox

import (
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
	cylinder := Cylinder{Center: Point{0.0, 0.0, 0.0}, Height: 2.0, Radius: 6.0}
	rayOrigin := Point{0.0, 0.0, -3.0}
	rayDirection := Point{0.0, 0.0, 1.0}
	maxDistance := 100.0
	intersects, intersectPoint, _ := cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if !intersects {
		t.Errorf("Ray should intersect cylinder")
	}
	if !almostEqual(intersectPoint.x, 0.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.y, 0.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.z, -1.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}

	rayOrigin = Point{2.0, 2.0, -3.0}
	rayDirection = Point{0.0, 0.0, 1.0}
	maxDistance = 100.0
	intersects, intersectPoint, _ = cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if !intersects {
		t.Errorf("Ray should intersect cylinder")
	}
	if !almostEqual(intersectPoint.x, 2.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.y, 2.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}
	if !almostEqual(intersectPoint.z, -1.0, 0.0) {
		t.Errorf("Incorrect intersection point: %v", intersectPoint)
	}

	rayOrigin = Point{2.0, 2.0, -3.0}
	rayDirection = Point{0.0, 0.0, -1.0}
	maxDistance = 0.5
	intersects, intersectPoint, _ = cylinder.RayIntersection(rayOrigin, rayDirection, maxDistance)
	if intersects {
		t.Errorf("Ray should not intersect cylinder")
	}

}

func TestDistanceToPlane(t *testing.T) {
	plane := Plane{
		Normal: Point{0.0, 0.0, 1.0},
		Center: Point{x: 0.0, y: 0.0, z: 0.0},
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
		Center: Point{x: 0.0, y: 0.0, z: 0.0},
	}
	if !almostEqual(plane.DistanceToPoint(Point{2.2, 1.0, 10.0}), 3.505910087489654, 0.0) {
		t.Errorf("Incorrect distance to plane")
	}
}
