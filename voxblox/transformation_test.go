package voxblox

import (
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
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

	transformationInversed := transformation.Inverse()
	point = transformationInversed.transformPoint(point)
	if !almostEqual(point[0], 0.714538097, kEpsilon) {
		t.Errorf("Expected %f, got %f", 0.714538097, point[0])
	}
	if !almostEqual(point[1], -2.8530097, kEpsilon) {
		t.Errorf("Expected %f, got %f", -2.8530097, point[1])
	}
	if !almostEqual(point[2], -1.72378588, kEpsilon) {
		t.Errorf("Expected %f, got %f", -1.72378588, point[2])
	}
}

func TestInverseTransform(t *testing.T) {
	transformation := Transformation{
		Position: vec3.T{0, 6, 2},
		Rotation: quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	inverse := transformation.Inverse()
	if inverse.Rotation[0] != -0.0353406072 {
		t.Errorf("Expected %f, got %f", -0.0353406072, inverse.Rotation[0])
	}
	if inverse.Rotation[1] != 0.0353406072 {
		t.Errorf("Expected %f, got %f", 0.0353406072, inverse.Rotation[1])
	}
	if inverse.Rotation[2] != 0.706223071 {
		t.Errorf("Expected %f, got %f", 0.706223071, inverse.Rotation[2])
	}
	if inverse.Rotation[3] != 0.706223071 {
		t.Errorf("Expected %f, got %f", -0.706223071, inverse.Rotation[3])
	}
	if !almostEqual(inverse.Position[0], 6.0, 0.001) {
		t.Errorf("Expected %f, got %f", 6.0, inverse.Position[0])
	}
	if !almostEqual(inverse.Position[1], -0.2, 0.001) {
		t.Errorf("Expected %f, got %f", -0.2, inverse.Position[1])
	}
	if !almostEqual(inverse.Position[2], -2.0, 0.01) {
		t.Errorf("Expected %f, got %f", -2.0, inverse.Position[2])
	}
}
