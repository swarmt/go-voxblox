package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
	"testing"
)

func TestTransformPoint(t *testing.T) {
	transformation := Transformation{
		Position: vec3.T{0, 6, 2},
		Rotation: quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	point := transformation.transformPoint(Point{0.714538097, -2.8530097, -1.72378588})
	if !almostEqual(point[0], -2.66666508, 0.0) {
		t.Errorf("Expected %f, got %f", -2.66666508, point[0])
	}
	if !almostEqual(point[1], 5.2854619, 0.0) {
		t.Errorf("Expected %f, got %f", 5.2854619, point[1])
	}
	if !almostEqual(point[2], 1.1920929e-07, 0.0) {
		t.Errorf("Expected %f, got %f", 1.1920929e-07, point[2])
	}
}
