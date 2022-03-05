package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
)

type Transformation struct {
	Position vec3.T
	Rotation quaternion.T
}

// TransformPoint by rotation and translation.
func (t Transformation) TransformPoint(point Point) vec3.T {
	rotatedPoint := t.Rotation.RotatedVec3(&point)
	return vec3.Add(&rotatedPoint, &t.Position)
}
