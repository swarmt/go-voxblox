package voxblox

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
)

func TestTransformPoint(t *testing.T) {
	transformation := Transformation{
		Translation: Point{0, 6, 2},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	point := transformation.transformPoint(Point{0.714538097, -2.8530097, -1.72378588})
	assert.InEpsilon(t, -2.66666508, point[0], kEpsilon)
	assert.InEpsilon(t, 5.2854619, point[1], kEpsilon)
	assert.InEpsilon(t, 0.0000002384665951371545, point[2], kEpsilon)

	transformationInversed := transformation.inverse()
	point = transformationInversed.transformPoint(point)
	assert.InEpsilon(t, 0.714538097, point[0], kEpsilon)
	assert.InEpsilon(t, -2.8530097, point[1], kEpsilon)
	assert.InEpsilon(t, -1.72378588, point[2], kEpsilon)
}

func TestInverseTransform(t *testing.T) {
	transformation := Transformation{
		Translation: Point{0, 6, 2},
		Rotation:    quaternion.T{0.0353406072, -0.0353406072, -0.706223071, 0.706223071},
	}
	inverse := transformation.inverse()
	assert.Equal(t, -0.0353406072, inverse.Rotation[0])
	assert.Equal(t, 0.0353406072, inverse.Rotation[1])
	assert.Equal(t, 0.706223071, inverse.Rotation[2])
	assert.Equal(t, 0.706223071, inverse.Rotation[3])
	assert.InEpsilon(t, 6, inverse.Translation[0], 0.002)
	assert.InEpsilon(t, -0.2, inverse.Translation[1], 0.002)
	assert.InEpsilon(t, -2, inverse.Translation[2], 0.005)
}
