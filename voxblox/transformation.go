package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
)

type Transformation struct {
	Position vec3.T
	Rotation quaternion.T
}

func (t Transformation) Inverted() Transformation {
	return Transformation{
		Position: t.Position.Inverted(),
		Rotation: t.Rotation.Inverted(),
	}
}
