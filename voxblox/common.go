package voxblox

import (
	"math"
)

// Constants
const kEpsilon = 1e-6 // Used for coordinates

// IndexType is the type used for indexing blocks and voxels.
type IndexType = [3]int32

// Point is a matrix of 3x1
type Point struct {
	x, y, z float64
}

// MaxInt32
func MaxInt32(x, y int32) int32 {
	if x < y {
		return y
	}
	return x
}

func MinInt32(x, y int32) int32 {
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

// getGridIndexFromScaledPoint returns the grid index of a point given the coordinate
func getGridIndexFromScaledPoint(scaledPoint Point) IndexType {
	return IndexType{
		int32(math.Floor(scaledPoint.x + kEpsilon)),
		int32(math.Floor(scaledPoint.y + kEpsilon)),
		int32(math.Floor(scaledPoint.z + kEpsilon)),
	}
}

func getGridIndexFromPoint(point Point, blockSizeInv float64) IndexType {
	return IndexType{
		int32(math.Floor(point.x*blockSizeInv + kEpsilon)),
		int32(math.Floor(point.y*blockSizeInv + kEpsilon)),
		int32(math.Floor(point.z*blockSizeInv + kEpsilon)),
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
