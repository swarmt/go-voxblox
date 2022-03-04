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

type PointCloud struct {
	Width  int
	Height int
	Points []Point
}

// Point is a matrix of 3x1
// TODO: Is there a way to alias a vec3 matrix to X Y Z?
type Point struct {
	X, Y, Z float64
}

func (p Point) asVec2() *vec2.T {
	return &vec2.T{p.X, p.Y}
}

func (p Point) asVec3() *vec3.T {
	return &vec3.T{p.X, p.Y, p.Z}
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
	return Point{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func addPoints(a, b Point) Point {
	return Point{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func almostEqual(a, b, e float64) bool {
	return math.Abs(a-b) <= e+kEpsilon
}

// getGridIndexFromScaledPoint returns the grid index of a point given the coordinate
func getGridIndexFromScaledPoint(scaledPoint Point) IndexType {
	return IndexType{
		int(math.Floor(scaledPoint.X + kEpsilon)),
		int(math.Floor(scaledPoint.Y + kEpsilon)),
		int(math.Floor(scaledPoint.Z + kEpsilon)),
	}
}

func getGridIndexFromPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Floor(point.X*blockSizeInv + kEpsilon)),
		int(math.Floor(point.Y*blockSizeInv + kEpsilon)),
		int(math.Floor(point.Z*blockSizeInv + kEpsilon)),
	}
}

func getGridIndexFromOriginPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int(math.Round(point.X * blockSizeInv)),
		int(math.Round(point.Y * blockSizeInv)),
		int(math.Round(point.Z * blockSizeInv)),
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
