package voxblox

import (
	"log"
	"math"
	"time"

	"github.com/ungerik/go3d/float64/vec3"
)

// Constants
const kEpsilon = 1e-6 // Used for coordinates

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

func SubIndex(a, b IndexType) IndexType {
	return IndexType{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
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

func almostEqual(a, b, e float64) bool {
	return math.Abs(a-b) <= e+kEpsilon
}

func Sgn(a float64) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

// timeTrack is a helper function for timing the execution of a function.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// IndexToPoint converts an index to a point.
func IndexToPoint(index IndexType) Point {
	return Point{float64(index[0]), float64(index[1]), float64(index[2])}
}

// getGridIndexFromScaledPoint returns the grid Index of a point given the coordinate
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

func getBlockIndexFromGlobalVoxelIndex(
	globalVoxelIndex IndexType,
	voxelsPerSideInv float64,
) IndexType {
	return IndexType{
		int(math.Floor(float64(globalVoxelIndex[0]) * voxelsPerSideInv)),
		int(math.Floor(float64(globalVoxelIndex[1]) * voxelsPerSideInv)),
		int(math.Floor(float64(globalVoxelIndex[2]) * voxelsPerSideInv)),
	}
}

func getLocalFromGlobalVoxelIndex(
	globalVoxelIndex IndexType,
	blockIndex IndexType,
	voxelsPerSide int,
) IndexType {
	return IndexType{
		globalVoxelIndex[0] - blockIndex[0]*voxelsPerSide,
		globalVoxelIndex[1] - blockIndex[1]*voxelsPerSide,
		globalVoxelIndex[2] - blockIndex[2]*voxelsPerSide,
	}
}

func getGlobalVoxelIndexFromBlockAndVoxelIndex(
	blockIndex IndexType,
	voxelIndex IndexType,
	voxelsPerSide int,
) IndexType {
	return IndexType{
		blockIndex[0]*voxelsPerSide + voxelIndex[0],
		blockIndex[1]*voxelsPerSide + voxelIndex[1],
		blockIndex[2]*voxelsPerSide + voxelIndex[2],
	}
}
