package voxblox

import (
	"testing"
)

var (
	blockSize float64
	tsdfLayer *TsdfLayer
)

func init() {
	tsdfVoxelSize := 0.1
	tsdfVoxelsPerSide := 8
	blockSize = tsdfVoxelSize * float64(tsdfVoxelsPerSide)
	tsdfLayer = NewTsdfLayer(tsdfVoxelSize, tsdfVoxelsPerSide)
}

func TestTsdfMapBlockAllocation(t *testing.T) {
	// Layer should have no blocks by default
	if tsdfLayer.getBlockCount() != 0 {
		t.Errorf("Expected no blocks in Layer, got %d", tsdfLayer.getBlockCount())
	}
	tsdfLayer.getBlockByCoordinates(Point{0, 0.15, 0})
	if tsdfLayer.getBlockCount() != 1 {
		t.Errorf("Expected one block in Layer, got %d", tsdfLayer.getBlockCount())
	}
	tsdfLayer.getBlockByCoordinates(Point{0, 0.13, 0})
	if tsdfLayer.getBlockCount() != 1 {
		t.Errorf("Expected one block in Layer, got %d", tsdfLayer.getBlockCount())
	}
	tsdfLayer.getBlockByCoordinates(Point{-10.0, 13.5, 20.0})
	if tsdfLayer.getBlockCount() != 2 {
		t.Errorf("Expected two blocks in Layer, got %d", tsdfLayer.getBlockCount())
	}
}

func TestTsdfMapIndexLookups(t *testing.T) {
	// BLOCK 0 0 0 (coordinate at Origin)
	pointIn000 := Point{0, 0, 0}
	block000 := tsdfLayer.getBlockByCoordinates(pointIn000)
	if block000 == nil {
		t.Errorf("Expected block000 to not be nil")
	}
	index000 := IndexType{0, 0, 0}
	if getBlockIndexFromCoordinates(pointIn000, tsdfLayer.BlockSizeInv) != index000 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if block000.Origin != pointIn000 {
		t.Errorf("Expected block000.Origin to be {0, 0, 0}, got %v", block000.Origin)
	}

	// BLOCK 0 0 0 (coordinate within block)
	pointIn000v2 := Point{0.0, tsdfLayer.VoxelSize, 0.0}
	block000v2 := tsdfLayer.getBlockByCoordinates(pointIn000v2)
	if block000v2 == nil {
		t.Errorf("Expected block000v2 to not be nil")
	}
	if getBlockIndexFromCoordinates(pointIn000v2, tsdfLayer.BlockSizeInv) != index000 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if block000 != block000v2 {
		t.Errorf("Expected block000 to be block000v2")
	}

	// BLOCK 1 1 1 (coordinate at Origin)
	pointIn111 := Point{blockSize, blockSize, blockSize}
	block111 := tsdfLayer.getBlockByCoordinates(pointIn111)
	if block111 == nil {
		t.Errorf("Expected block111 to not be nil")
	}
	index111 := IndexType{1, 1, 1}
	if getBlockIndexFromCoordinates(pointIn111, tsdfLayer.BlockSizeInv) != index111 {
		t.Errorf("Expected {1, 1, 1} to be returned by GetBlockPtrByCoordinates")
	}
	if block111.Origin != (Point{blockSize, blockSize, blockSize}) {
		t.Errorf("Expected block111.Origin to be {%f, %f, %f}, got %v",
			blockSize, blockSize, blockSize, block111.Origin)
	}
	if block111.Index != index111 {
		t.Errorf("Expected block111.Index to be {1, 1, 1}, got %v", block111.Index)
	}

	// BLOCK 1 1 1 (coordinate within block)
	pointIn111v2 := Point{blockSize, blockSize + tsdfLayer.VoxelSize, blockSize}
	block111v2 := tsdfLayer.getBlockByCoordinates(pointIn111v2)
	if block111v2 == nil {
		t.Errorf("Expected block111v2 to not be nil")
	}
	if getBlockIndexFromCoordinates(pointIn111v2, tsdfLayer.BlockSizeInv) != index111 {
		t.Errorf("Expected {1, 1, 1} to be returned by GetBlockPtrByCoordinates")
	}
	if block111 != block111v2 {
		t.Errorf("Expected block111 to be block111v2")
	}
	if block111v2.Origin != (Point{blockSize, blockSize, blockSize}) {
		t.Errorf("Expected block111.Origin to be {%f, %f, %f}, got %v",
			blockSize, blockSize, blockSize, block111.Origin)
	}
	if block111v2.Index != index111 {
		t.Errorf("Expected block111.Index to be {1, 1, 1}, got %v", block111.Index)
	}
	if block111v2.BlockSize != blockSize {
		t.Errorf("Expected block111.BlockSize to be %f, got %f", blockSize, block111v2.BlockSize)
	}

	// BLOCK -1 -1 -1 (coordinate at Origin)
	pointInNeg111 := Point{-blockSize, -blockSize, -blockSize}
	blockNeg111 := tsdfLayer.getBlockByCoordinates(pointInNeg111)
	if blockNeg111 == nil {
		t.Errorf("Expected blockNeg111 to not be nil")
	}
	indexNeg111 := IndexType{-1, -1, -1}
	if getBlockIndexFromCoordinates(pointInNeg111, tsdfLayer.BlockSizeInv) != indexNeg111 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if blockNeg111.Origin != (Point{-blockSize, -blockSize, -blockSize}) {
		t.Errorf("Expected blockNeg111.Origin to be {%f, %f, %f}, got %v",
			-blockSize, -blockSize, -blockSize, blockNeg111.Origin)
	}
	if blockNeg111.BlockSize != blockSize {
		t.Errorf(
			"Expected blockNeg111.BlockSize to be %f, got %f",
			blockSize,
			blockNeg111.BlockSize,
		)
	}

	// BLOCK -1 -1 -1 (coordinate within block)
	pointInNeg111v2 := Point{-blockSize, -blockSize + tsdfLayer.VoxelSize, -blockSize}
	blockNeg111v2 := tsdfLayer.getBlockByCoordinates(pointInNeg111v2)
	if blockNeg111v2 == nil {
		t.Errorf("Expected blockNeg111v2 to not be nil")
	}
	if getBlockIndexFromCoordinates(pointInNeg111v2, tsdfLayer.BlockSizeInv) != indexNeg111 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if blockNeg111 != blockNeg111v2 {
		t.Errorf("Expected blockNeg111 to be blockNeg111v2")
	}
	if blockNeg111v2.Origin != (Point{-blockSize, -blockSize, -blockSize}) {
		t.Errorf("Expected blockNeg111.Origin to be {%f, %f, %f}, got %v",
			-blockSize, -blockSize, -blockSize, blockNeg111.Origin)
	}
	if blockNeg111v2.BlockSize != blockSize {
		t.Errorf(
			"Expected blockNeg111.BlockSize to be %f, got %f",
			blockSize,
			blockNeg111v2.BlockSize,
		)
	}

	// TsdfBlock 0 0 0

	// Voxel 0 1 0
	pointIn000 = Point{0.0, 1.0 * tsdfLayer.VoxelSize, 0.0}
	if block000.getVoxelPtrByCoordinates(pointIn000) == nil {
		t.Errorf("Expected pointIn000 to not be nil")
	}

	// TODO: Ignoring linear indexing for now.
	// TODO: Don't think it is worth the additional code complexity but I could be wrong.
	voxelIndex := block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	if voxelIndex != (IndexType{0, 1, 0}) {
		t.Errorf("Expected {0, 1, 0} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}

	pointIn000center := block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center[0], pointIn000[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[1], pointIn000[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[2], pointIn000[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointIn000center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointIn000[0], pointIn000[1], pointIn000[2],
			pointIn000center[0], pointIn000center[1], pointIn000center[2])
	}

	// Voxel 0 0 0
	pointIn000 = Point{0.0, 0.0, 0.0}
	if block000.getVoxelPtrByCoordinates(pointIn000) == nil {
		t.Errorf("Expected pointIn000 to not be nil")
	}
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	if voxelIndex != (IndexType{0, 0, 0}) {
		t.Errorf("Expected {0, 0, 0} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center[0], pointIn000[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[1], pointIn000[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[2], pointIn000[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointIn000center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointIn000[0], pointIn000[1], pointIn000[2],
			pointIn000center[0], pointIn000center[1], pointIn000center[2])
	}

	// Voxel 7 7 7
	pointIn000 = Point{
		7.0 * tsdfLayer.VoxelSize,
		7.0 * tsdfLayer.VoxelSize,
		7.0 * tsdfLayer.VoxelSize,
	}
	if block000.getVoxelPtrByCoordinates(pointIn000) == nil {
		t.Errorf("Expected pointIn000 to not be nil")
	}
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	if voxelIndex != (IndexType{7, 7, 7}) {
		t.Errorf("Expected {7, 7, 7} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center[0], pointIn000[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[1], pointIn000[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[2], pointIn000[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointIn000center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointIn000[0], pointIn000[1], pointIn000[2],
			pointIn000center[0], pointIn000center[1], pointIn000center[2])
	}

	// TsdfBlock -1 -1 -1

	// Voxel 0 0 0
	pointInNeg111 = Point{
		-1.0 * blockNeg111.BlockSize,
		-1.0 * blockNeg111.BlockSize,
		-1.0 * blockNeg111.BlockSize,
	}
	if blockNeg111.getVoxelPtrByCoordinates(pointInNeg111) == nil {
		t.Errorf("Expected pointInNeg111 to not be nil")
	}

	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	if voxelIndex != (IndexType{0, 0, 0}) {
		t.Errorf("Expected {0, 0, 0} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center[0], pointInNeg111[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[1], pointInNeg111[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn000center[2], pointInNeg111[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointIn000center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointIn000[0], pointIn000[1], pointIn000[2],
			pointIn000center[0], pointIn000center[1], pointIn000center[2])
	}

	// Voxel 7 7 7
	pointInNeg111 = Point{-kEpsilon, -kEpsilon, -kEpsilon}
	if blockNeg111.getVoxelPtrByCoordinates(pointInNeg111) == nil {
		t.Errorf("Expected pointInNeg111 to not be nil")
	}
	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	if voxelIndex != (IndexType{7, 7, 7}) {
		t.Errorf("Expected {7, 7, 7} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn777center := blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn777center[0], pointInNeg111[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn777center[1], pointInNeg111[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointIn777center[2], pointInNeg111[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointIn000center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointInNeg111[0], pointInNeg111[1], pointInNeg111[2],
			pointIn777center[0], pointIn777center[1], pointIn777center[2])
	}

	// TsdfBlock -1 -1 0

	// Voxel 3 6 5
	pointInNeg110 := Point{
		-5.0 * tsdfLayer.VoxelSize,
		-2.0 * tsdfLayer.VoxelSize,
		5.0 * tsdfLayer.VoxelSize,
	}
	blockNeg1Neg1Pos0 := tsdfLayer.getBlockByCoordinates(pointInNeg110)
	if blockNeg1Neg1Pos0 == nil {
		t.Errorf("Expected blockNeg1Neg1Pos0 to not be nil")
	}
	if getBlockIndexFromCoordinates(pointInNeg110, tsdfLayer.BlockSizeInv) !=
		(IndexType{-1, -1, 0}) {
		t.Errorf("Expected Index to be {-1, -1, 0}")
	}
	if blockNeg1Neg1Pos0.getVoxelPtrByCoordinates(pointInNeg110) == nil {
		t.Errorf("Expected pointInNeg110 to not be nil")
	}
	voxelIndex = blockNeg1Neg1Pos0.computeTruncatedVoxelIndexFromCoordinates(pointInNeg110)
	if voxelIndex != (IndexType{3, 6, 5}) {
		t.Errorf("Expected {3, 6, 5} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointInNeg110center := blockNeg1Neg1Pos0.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointInNeg110center[0], pointInNeg110[0], tsdfLayer.VoxelSize) ||
		!almostEqual(pointInNeg110center[1], pointInNeg110[1], tsdfLayer.VoxelSize) ||
		!almostEqual(pointInNeg110center[2], pointInNeg110[2], tsdfLayer.VoxelSize) {
		t.Errorf("Expected pointInNeg110center to be {%f, %f, %f}, got {%f, %f, %f}",
			pointInNeg110[0], pointInNeg110[1], pointInNeg110[2],
			pointInNeg110center[0], pointInNeg110center[1], pointInNeg110center[2])
	}
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

				if blockIndex != blockIndexFromOrigin {
					t.Errorf(
						"Expected Index (%v) to be equal to blockIndexFromOrigin (%v)",
						blockIndex,
						blockIndexFromOrigin,
					)
				}
			}
		}
	}
}
