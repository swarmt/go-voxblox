package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
)

type Transformation struct {
	Position Point
	Rotation quaternion.T
}

// transformPoint by rotation and translation.
func (t Transformation) transformPoint(point Point) vec3.T {
	rotatedPoint := t.Rotation.RotatedVec3(&point)
	return vec3.Add(&rotatedPoint, &t.Position)
}

// Inverse returns the inverse transformation.
func (t Transformation) Inverse() Transformation {
	rotationInverted := t.Rotation.Inverted()
	pointRotated := rotationInverted.RotatedVec3(&t.Position)
	return Transformation{
		Position: pointRotated.Inverted(),
		Rotation: rotationInverted,
	}
}

// transformPointCloud by rotation and translation.
func transformPointCloud(transformation Transformation, pointCloud PointCloud) PointCloud {
	transformedPoints := make([]Point, len(pointCloud.Points))
	for i, point := range pointCloud.Points {
		transformedPoints[i] = transformation.transformPoint(point)
	}
	return PointCloud{
		Width:  pointCloud.Width,
		Height: pointCloud.Height,
		Points: transformedPoints,
	}
}
