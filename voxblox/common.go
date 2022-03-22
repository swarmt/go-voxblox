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

// Color RGBA
type Color = [4]uint8

// colors
var (
	ColorWhite = Color{255, 255, 255, 255}
	ColorRed   = Color{255, 0, 0, 255}
)

// PointCloud is a collection of points
type PointCloud struct {
	Width  int
	Height int
	Points []Point
	Colors []Color
}

// Point is 3x1 vector
// X, Y, Z are the coordinates
type Point = vec3.T

func subIndex(a, b IndexType) IndexType {
	return IndexType{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func addIndex(a, b IndexType) IndexType {
	return IndexType{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func minInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func almostEqual(a, b, e float64) bool {
	return math.Abs(a-b) <= e+kEpsilon
}

func sgn(a float64) int {
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
	log.Printf("%s: %s", name, elapsed)
}

// IndexToPoint converts an Index to a point.
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

func getBlockIndexFromCoordinates(point Point, blockSizeInv float64) IndexType {
	return getGridIndexFromPoint(point, blockSizeInv)
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

func blendTwoColors(
	firstColor Color,
	firstWeight float64,
	secondColor Color,
	secondWeight float64,
) Color {
	totalWeight := firstWeight + secondWeight
	firstWeight /= totalWeight
	secondWeight /= totalWeight

	newR := uint8(float64(firstColor[0])*firstWeight + float64(secondColor[0])*secondWeight)
	newG := uint8(float64(firstColor[1])*firstWeight + float64(secondColor[1])*secondWeight)
	newB := uint8(float64(firstColor[2])*firstWeight + float64(secondColor[2])*secondWeight)
	newA := uint8(float64(firstColor[3])*firstWeight + float64(secondColor[3])*secondWeight)

	return Color{newR, newG, newB, newA}
}
