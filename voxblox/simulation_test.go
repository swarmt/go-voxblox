package voxblox

import (
	"testing"
)

func TestDistanceToCylinder(t *testing.T) {
	cylinder := NewCylinder(Point{0.0, 0.0, 0.0}, 6.0, 2.0)
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
