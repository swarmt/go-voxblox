package voxblox

import (
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"math"
)

// Constants
const kEpsilon = 1e-6 // Used for coordinates

// IndexType is the type used for indexing blocks and voxels.
type IndexType = [3]int

type Color = [3]uint8

// Point is a matrix of 3x1
// TODO: Is there a way to alias a vec3 matrix to x y z?
type Point struct {
	x, y, z float64
}

func (p Point) asVec2() *vec2.T {
	return &vec2.T{p.x, p.y}
}

func (p Point) asVec3() *vec3.T {
	return &vec3.T{p.x, p.y, p.z}
}

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

func subtractPoints(a, b Point) Point {
	return Point{a.x - b.x, a.y - b.y, a.z - b.z}
}

func addPoints(a, b Point) Point {
	return Point{a.x + b.x, a.y + b.y, a.z + b.z}
}

func almostEqual(a, b, e float64) bool {
	return math.Abs(a-b) <= e+kEpsilon
}

// getGridIndexFromScaledPoint returns the grid index of a point given the coordinate
func getGridIndexFromScaledPoint(scaledPoint Point) IndexType {
	return IndexType{
		int(math.Floor(scaledPoint.x + kEpsilon)),
		int(math.Floor(scaledPoint.y + kEpsilon)),
		int(math.Floor(scaledPoint.z + kEpsilon)),
	}
}

func getGridIndexFromPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Floor(point.x*blockSizeInv + kEpsilon)),
		int(math.Floor(point.y*blockSizeInv + kEpsilon)),
		int(math.Floor(point.z*blockSizeInv + kEpsilon)),
	}
}

func getGridIndexFromOriginPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Round(point.x * blockSizeInv)),
		int(math.Round(point.y * blockSizeInv)),
		int(math.Round(point.z * blockSizeInv)),
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
