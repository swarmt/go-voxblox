package voxblox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGridIndexFromPoint(t *testing.T) {
	assert.Equal(
		t,
		IndexType{0, 105, 0},
		getGridIndexFromPoint(Point{1.31130219e-06, 5.2854619, 1.1920929e-07}, 2.0*10.0),
	)
	assert.Equal(
		t,
		IndexType{-1, 105, 0},
		getGridIndexFromPoint(Point{-0.0166654587, 5.2854619, 1.1920929e-07}, 2.0*10.0),
	)
	assert.Equal(
		t,
		IndexType{-21, 43, 0},
		getGridIndexFromPoint(Point{-2.05839586, 4.35501623, 1.1920929e-07}, 10.0),
	)
}

func TestGetCenterPointFromGridIndex(t *testing.T) {
	centerPoint := getCenterPointFromGridIndex(IndexType{-2, 51, -3}, 0.1)
	testPoint := Point{-0.15, 5.15, -0.25}
	assert.InEpsilon(t, centerPoint[0], testPoint[0], kEpsilon)
	assert.InEpsilon(t, centerPoint[1], testPoint[1], kEpsilon)
	assert.InEpsilon(t, centerPoint[2], testPoint[2], kEpsilon)

	centerPoint = getCenterPointFromGridIndex(IndexType{-2, 56, 9}, 0.1)
	testPoint = Point{-0.15, 5.65, 0.95}
	assert.InEpsilon(t, centerPoint[0], testPoint[0], kEpsilon)
	assert.InEpsilon(t, centerPoint[1], testPoint[1], kEpsilon)
	assert.InEpsilon(t, centerPoint[2], testPoint[2], kEpsilon)
}

func TestBlendTwoColors(t *testing.T) {
	assert.Equal(
		t,
		blendTwoColors(Color{0, 0, 0}, 0, Color{255, 255, 255}, 1.0),
		Color{255, 255, 255},
	)
	assert.Equal(
		t,
		blendTwoColors(Color{255, 255, 255}, 0.500417829, Color{255, 255, 255}, 0.499582082),
		Color{255, 255, 255},
	)
}
