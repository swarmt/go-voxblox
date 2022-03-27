package voxblox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTsdfMapBlockAllocation(t *testing.T) {
	tsdfLayer := NewTsdfLayer(0.1, 8)
	assert.Equal(t, 0, tsdfLayer.getBlockCount())

	tsdfLayer.getBlockByCoordinates(Point{0, 0.15, 0})
	assert.Equal(t, 1, tsdfLayer.getBlockCount())

	tsdfLayer.getBlockByCoordinates(Point{0, 0.13, 0})
	assert.Equal(t, 1, tsdfLayer.getBlockCount())

	tsdfLayer.getBlockByCoordinates(Point{-10.0, 13.5, 20.0})
	assert.Equal(t, 2, tsdfLayer.getBlockCount())
}

func TestTsdfMapIndexLookups(t *testing.T) {
	tsdfVoxelSize := 0.1
	tsdfVoxelsPerSide := 8
	blockSize := tsdfVoxelSize * float64(tsdfVoxelsPerSide)
	tsdfLayer := NewTsdfLayer(tsdfVoxelSize, tsdfVoxelsPerSide)

	// BLOCK 0 0 0 (coordinate at Origin)
	pointIn000 := Point{0, 0, 0}
	block000 := tsdfLayer.getBlockByCoordinates(pointIn000)
	assert.NotNil(t, block000)

	index000 := IndexType{0, 0, 0}
	assert.Equal(t, index000, getBlockIndexFromCoordinates(pointIn000, tsdfLayer.BlockSizeInv))
	assert.Equal(t, pointIn000, block000.Origin)

	// BLOCK 0 0 0 (coordinate within block)
	pointIn000v2 := Point{0.0, tsdfLayer.VoxelSize, 0.0}
	block000v2 := tsdfLayer.getBlockByCoordinates(pointIn000v2)
	assert.NotNil(t, block000v2)
	assert.Equal(t, index000, getBlockIndexFromCoordinates(pointIn000v2, tsdfLayer.BlockSizeInv))
	assert.Equal(t, block000, block000v2)

	// BLOCK 1 1 1 (coordinate at Origin)
	pointIn111 := Point{blockSize, blockSize, blockSize}
	block111 := tsdfLayer.getBlockByCoordinates(pointIn111)
	assert.NotNil(t, block111)
	index111 := IndexType{1, 1, 1}
	assert.Equal(t, index111, getBlockIndexFromCoordinates(pointIn111, tsdfLayer.BlockSizeInv))
	assert.Equal(t, Point{blockSize, blockSize, blockSize}, block111.Origin)
	assert.Equal(t, index111, block111.Index)

	// BLOCK 1 1 1 (coordinate within block)
	pointIn111v2 := Point{blockSize, blockSize + tsdfLayer.VoxelSize, blockSize}
	block111v2 := tsdfLayer.getBlockByCoordinates(pointIn111v2)
	assert.NotNil(t, block111v2)
	assert.Equal(t, index111, getBlockIndexFromCoordinates(pointIn111v2, tsdfLayer.BlockSizeInv))
	assert.Equal(t, block111, block111v2)
	assert.Equal(t, Point{blockSize, blockSize, blockSize}, block111v2.Origin)
	assert.Equal(t, index111, block111v2.Index)
	assert.Equal(t, blockSize, block111v2.BlockSize)

	// BLOCK -1 -1 -1 (coordinate at Origin)
	pointInNeg111 := Point{-blockSize, -blockSize, -blockSize}
	blockNeg111 := tsdfLayer.getBlockByCoordinates(pointInNeg111)
	assert.NotNil(t, blockNeg111)
	indexNeg111 := IndexType{-1, -1, -1}
	assert.Equal(t, indexNeg111, getBlockIndexFromCoordinates(pointInNeg111, tsdfLayer.BlockSizeInv))
	assert.Equal(t, Point{-blockSize, -blockSize, -blockSize}, blockNeg111.Origin)
	assert.Equal(t, blockSize, blockNeg111.BlockSize)

	// BLOCK -1 -1 -1 (coordinate within block)
	pointInNeg111v2 := Point{-blockSize, -blockSize + tsdfLayer.VoxelSize, -blockSize}
	blockNeg111v2 := tsdfLayer.getBlockByCoordinates(pointInNeg111v2)
	assert.NotNil(t, blockNeg111v2)
	assert.Equal(t, indexNeg111, getBlockIndexFromCoordinates(pointInNeg111v2, tsdfLayer.BlockSizeInv))
	assert.Equal(t, blockNeg111, blockNeg111v2)
	assert.Equal(t, Point{-blockSize, -blockSize, -blockSize}, blockNeg111v2.Origin)
	assert.Equal(t, blockSize, blockNeg111v2.BlockSize)

	// TsdfBlock 0 0 0

	// Voxel 0 1 0
	pointIn000 = Point{0.0, 1.0 * tsdfLayer.VoxelSize, 0.0}
	assert.NotNil(t, block000.getVoxelPtrByCoordinates(pointIn000))
	voxelIndex := block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	assert.Equal(t, IndexType{0, 1, 0}, voxelIndex)

	pointIn000center := block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointIn000[0], pointIn000center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[1], pointIn000center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[2], pointIn000center[2], tsdfLayer.VoxelSize+kEpsilon)

	// Voxel 0 0 0
	pointIn000 = Point{0.0, 0.0, 0.0}
	assert.NotNil(t, block000.getVoxelPtrByCoordinates(pointIn000))
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	assert.Equal(t, IndexType{0, 0, 0}, voxelIndex)
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointIn000[0], pointIn000center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[1], pointIn000center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[2], pointIn000center[2], tsdfLayer.VoxelSize+kEpsilon)

	// Voxel 7 7 7
	pointIn000 = Point{
		7.0 * tsdfLayer.VoxelSize,
		7.0 * tsdfLayer.VoxelSize,
		7.0 * tsdfLayer.VoxelSize,
	}
	assert.NotNil(t, block000.getVoxelPtrByCoordinates(pointIn000))
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	assert.Equal(t, IndexType{7, 7, 7}, voxelIndex)
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointIn000[0], pointIn000center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[1], pointIn000center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointIn000[2], pointIn000center[2], tsdfLayer.VoxelSize+kEpsilon)

	// TsdfBlock -1 -1 -1

	// Voxel 0 0 0
	pointInNeg111 = Point{
		-1.0 * blockNeg111.BlockSize,
		-1.0 * blockNeg111.BlockSize,
		-1.0 * blockNeg111.BlockSize,
	}
	assert.NotNil(t, blockNeg111.getVoxelPtrByCoordinates(pointInNeg111))
	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	assert.Equal(t, IndexType{0, 0, 0}, voxelIndex)
	pointIn000center = blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointInNeg111[0], pointIn000center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg111[1], pointIn000center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg111[2], pointIn000center[2], tsdfLayer.VoxelSize+kEpsilon)

	// Voxel 7 7 7
	pointInNeg111 = Point{-kEpsilon, -kEpsilon, -kEpsilon}
	assert.NotNil(t, blockNeg111.getVoxelPtrByCoordinates(pointInNeg111))
	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	assert.Equal(t, IndexType{7, 7, 7}, voxelIndex)
	pointIn777center := blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointInNeg111[0], pointIn777center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg111[1], pointIn777center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg111[2], pointIn777center[2], tsdfLayer.VoxelSize+kEpsilon)

	// TsdfBlock -1 -1 0

	// Voxel 3 6 5
	pointInNeg110 := Point{
		-5.0 * tsdfLayer.VoxelSize,
		-2.0 * tsdfLayer.VoxelSize,
		5.0 * tsdfLayer.VoxelSize,
	}
	blockNeg1Neg1Pos0 := tsdfLayer.getBlockByCoordinates(pointInNeg110)
	assert.NotNil(t, blockNeg1Neg1Pos0)
	assert.Equal(t, IndexType{-1, -1, 0}, getBlockIndexFromCoordinates(pointInNeg110, tsdfLayer.BlockSizeInv))
	assert.NotNil(t, blockNeg1Neg1Pos0.getVoxelPtrByCoordinates(pointInNeg110))
	voxelIndex = blockNeg1Neg1Pos0.computeTruncatedVoxelIndexFromCoordinates(pointInNeg110)
	assert.Equal(t, IndexType{3, 6, 5}, voxelIndex)
	pointInNeg110center := blockNeg1Neg1Pos0.computeCoordinatesFromVoxelIndex(voxelIndex)
	assert.InDelta(t, pointInNeg110[0], pointInNeg110center[0], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg110[1], pointInNeg110center[1], tsdfLayer.VoxelSize+kEpsilon)
	assert.InDelta(t, pointInNeg110[2], pointInNeg110center[2], tsdfLayer.VoxelSize+kEpsilon)
}

func TestComputeBlockIndexFromOriginFromBlockIndex(t *testing.T) {
	const kBlockVolumeDiameter = 100
	const kBlockSize = 0.32
	const kBlockSizeInv = 1.0 / kBlockSize
	const halfIndexRange = kBlockVolumeDiameter / 2

	for x := -halfIndexRange; x <= halfIndexRange; x++ {
		for y := -halfIndexRange; y <= halfIndexRange; y++ {
			for z := -halfIndexRange; z <= halfIndexRange; z++ {
				blockIndex := IndexType{x, y, z}
				blockOrigin := getOriginPointFromGridIndex(blockIndex, kBlockSize)
				blockIndexFromOrigin := getGridIndexFromOriginPoint(blockOrigin, kBlockSizeInv)
				assert.Equal(t, blockIndex, blockIndexFromOrigin)
			}
		}
	}
}
