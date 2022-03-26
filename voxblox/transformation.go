package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
)

type Transformation struct {
	Translation Point
	Rotation    quaternion.T
}

// transformPoint by rotation and translation.
func (t Transformation) transformPoint(point Point) vec3.T {
	rotatedPoint := t.Rotation.RotatedVec3(&point)
	return vec3.Add(&rotatedPoint, &t.Translation)
}

// inverse returns the inverse Transformation.
func (t Transformation) inverse() Transformation {
	rotationInverted := t.Rotation.Inverted()
	pointRotated := rotationInverted.RotatedVec3(&t.Translation)
	return Transformation{
		Translation: pointRotated.Inverted(),
		Rotation:    rotationInverted,
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
		Colors: pointCloud.Colors,
	}
}

func CombineTransformations(t1, t2 *Transformation) Transformation {
	rotation := quaternion.Mul(&t1.Rotation, &t2.Rotation)
	translation := vec3.Add(&t1.Translation, &t2.Translation)
	return Transformation{
		Translation: translation,
		Rotation:    rotation,
	}
}
