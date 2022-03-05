package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"math"
)

// Constants
const kEpsilon = 1e-6 // Used for coordinates

// IndexType is the type used for indexing blocks and voxels.
type IndexType = [3]int
type Color = [3]uint8

type PointCloud struct {
	Width  int
	Height int
	Points []Point
}

// Point is 3x1 vector
// X, Y, Z are the coordinates
type Point = vec3.T

func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func almostEqual(a, b, e float64) bool {
	return math.Abs(a-b) <= e+kEpsilon
}

// getGridIndexFromScaledPoint returns the grid index of a point given the coordinate
func getGridIndexFromScaledPoint(scaledPoint Point) IndexType {
	return IndexType{
		int(math.Floor(scaledPoint[0] + kEpsilon)),
		int(math.Floor(scaledPoint[1] + kEpsilon)),
		int(math.Floor(scaledPoint[2] + kEpsilon)),
	}
}

func getGridIndexFromPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Floor(point[0]*blockSizeInv + kEpsilon)),
		int(math.Floor(point[1]*blockSizeInv + kEpsilon)),
		int(math.Floor(point[2]*blockSizeInv + kEpsilon)),
	}
}

func getGridIndexFromOriginPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Round(point[0] * blockSizeInv)),
		int(math.Round(point[1] * blockSizeInv)),
		int(math.Round(point[2] * blockSizeInv)),
	}
}

func getOriginPointFromGridIndex(index IndexType, gridSize float64) Point {
	return Point{
		float64(index[0]) * gridSize,
		float64(index[1]) * gridSize,
		float64(index[2]) * gridSize,
	}
}

func getCenterPointFromGridIndex(idx IndexType, gridSize float64) Point {
	return Point{
		(float64(idx[0]) + 0.5) * gridSize,
		(float64(idx[1]) + 0.5) * gridSize,
		(float64(idx[2]) + 0.5) * gridSize,
	}
}
