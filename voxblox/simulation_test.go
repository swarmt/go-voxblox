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
}
