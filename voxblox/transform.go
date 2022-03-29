package voxblox

import (
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec3"
)

type Transform struct {
	Translation Point
	Rotation    quaternion.T
}

// transformPoint by rotation and translation.
func (t Transform) transformPoint(point Point) vec3.T {
	rotatedPoint := t.Rotation.RotatedVec3(&point)
	return vec3.Add(&rotatedPoint, &t.Translation)
}

// inverse returns the inverse Transform.
func (t Transform) inverse() Transform {
	rotationInverted := t.Rotation.Inverted()
	pointRotated := rotationInverted.RotatedVec3(&t.Translation)
	return Transform{
		Translation: pointRotated.Inverted(),
		Rotation:    rotationInverted,
	}
}

// transformPointCloud by rotation and translation.
func transformPointCloud(transformation Transform, pointCloud PointCloud) PointCloud {
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

func ApplyTransform(t1, t2 *Transform) Transform {
	rotation := quaternion.Mul(&t1.Rotation, &t2.Rotation)
	translation := vec3.Add(&t1.Translation, &t2.Translation)
	return Transform{
		Translation: translation,
		Rotation:    rotation,
	}
}

// interpolatePoints interpolates between two Points
func interpolatePoints(p1, p2 Point, f float64) Point {
	return Point{
		p1[0] + (p2[0]-p1[0])*f,
		p1[1] + (p2[1]-p1[1])*f,
		p1[2] + (p2[2]-p1[2])*f,
	}
}

// InterpolateTransform interpolates between two Transformations
func InterpolateTransform(t1, t2 Transform, alpha float64) Transform {
	return Transform{
		Translation: interpolatePoints(t1.Translation, t2.Translation, alpha),
		Rotation:    quaternion.Slerp(&t1.Rotation, &t2.Rotation, alpha),
	}
}
