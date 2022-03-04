package voxblox

import (
	"testing"
)

var blockSize float64
var tsdfMap *TsdfMap

func init() {
	var tsdfVoxelSize = 0.1
	var tsdfVoxelsPerSide = 8
	blockSize = tsdfVoxelSize * float64(tsdfVoxelsPerSide)
	tsdfMap = NewTsdfMap(tsdfVoxelSize, tsdfVoxelsPerSide)
}

func TestTsdfMapBlockAllocation(t *testing.T) {
	// TsdfLayer should have no blocks by default
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 0 {
		t.Errorf("Expected no blocks in Layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{X: 0, Y: 0.15, Z: 0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 1 {
		t.Errorf("Expected one block in Layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{X: 0, Y: 0.13, Z: 0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 1 {
		t.Errorf("Expected one block in Layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{X: -10.0, Y: 13.5, Z: 20.0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 2 {
		t.Errorf("Expected two blocks in Layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
}

func TestTsdfMapIndexLookups(t *testing.T) {
	// BLOCK 0 0 0 (coordinate at origin)
	pointIn000 := Point{X: 0, Y: 0, Z: 0}
	block000 := tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(pointIn000)
	if block000 == nil {
		t.Errorf("Expected block000 to not be nil")
	}
	index000 := IndexType{0, 0, 0}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointIn000) != index000 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if block000.Origin != pointIn000 {
		t.Errorf("Expected block000.Origin to be {0, 0, 0}, got %v", block000.Origin)
	}

	// BLOCK 0 0 0 (coordinate within block)
	pointIn000v2 := Point{X: 0.0, Y: tsdfMap.TsdfVoxelSize, Z: 0.0}
	block000v2 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointIn000v2)
	if block000v2 == nil {
		t.Errorf("Expected block000v2 to not be nil")
	}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointIn000v2) != index000 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if block000 != block000v2 {
		t.Errorf("Expected block000 to be block000v2")
	}

	// BLOCK 1 1 1 (coordinate at origin)
	pointIn111 := Point{X: blockSize, Y: blockSize, Z: blockSize}
	block111 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointIn111)
	if block111 == nil {
		t.Errorf("Expected block111 to not be nil")
	}
	index111 := IndexType{1, 1, 1}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointIn111) != index111 {
		t.Errorf("Expected {1, 1, 1} to be returned by GetBlockPtrByCoordinates")
	}
	if block111.Origin != (Point{X: blockSize, Y: blockSize, Z: blockSize}) {
		t.Errorf("Expected block111.Origin to be {%f, %f, %f}, got %v",
			blockSize, blockSize, blockSize, block111.Origin)
	}
	if block111.BlockIndex != index111 {
		t.Errorf("Expected block111.BlockIndex to be {1, 1, 1}, got %v", block111.BlockIndex)
	}

	// BLOCK 1 1 1 (coordinate within block)
	pointIn111v2 := Point{X: blockSize, Y: blockSize + tsdfMap.TsdfVoxelSize, Z: blockSize}
	block111v2 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointIn111v2)
	if block111v2 == nil {
		t.Errorf("Expected block111v2 to not be nil")
	}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointIn111v2) != index111 {
		t.Errorf("Expected {1, 1, 1} to be returned by GetBlockPtrByCoordinates")
	}
	if block111 != block111v2 {
		t.Errorf("Expected block111 to be block111v2")
	}
	if block111v2.Origin != (Point{X: blockSize, Y: blockSize, Z: blockSize}) {
		t.Errorf("Expected block111.Origin to be {%f, %f, %f}, got %v",
			blockSize, blockSize, blockSize, block111.Origin)
	}
	if block111v2.BlockIndex != index111 {
		t.Errorf("Expected block111.BlockIndex to be {1, 1, 1}, got %v", block111.BlockIndex)
	}
	if block111v2.BlockSize != blockSize {
		t.Errorf("Expected block111.BlockSize to be %f, got %f", blockSize, block111v2.BlockSize)
	}

	// BLOCK -1 -1 -1 (coordinate at origin)
	pointInNeg111 := Point{X: -blockSize, Y: -blockSize, Z: -blockSize}
	blockNeg111 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointInNeg111)
	if blockNeg111 == nil {
		t.Errorf("Expected blockNeg111 to not be nil")
	}
	indexNeg111 := IndexType{-1, -1, -1}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointInNeg111) != indexNeg111 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if blockNeg111.Origin != (Point{X: -blockSize, Y: -blockSize, Z: -blockSize}) {
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
	pointInNeg111v2 := Point{X: -blockSize, Y: -blockSize + tsdfMap.TsdfVoxelSize, Z: -blockSize}
	blockNeg111v2 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointInNeg111v2)
	if blockNeg111v2 == nil {
		t.Errorf("Expected blockNeg111v2 to not be nil")
	}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointInNeg111v2) != indexNeg111 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if blockNeg111 != blockNeg111v2 {
		t.Errorf("Expected blockNeg111 to be blockNeg111v2")
	}
	if blockNeg111v2.Origin != (Point{X: -blockSize, Y: -blockSize, Z: -blockSize}) {
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

	// Block 0 0 0

	// Voxel 0 1 0
	pointIn000 = Point{X: 0.0, Y: 1.0 * tsdfMap.TsdfVoxelSize, Z: 0.0}
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
	if !almostEqual(pointIn000center.X, pointIn000.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.X to be %f, got %f", pointIn000.X, pointIn000center.X)
	}
	if !almostEqual(pointIn000center.Y, pointIn000.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Y to be %f, got %f", pointIn000.Y, pointIn000center.Y)
	}
	if !almostEqual(pointIn000center.Z, pointIn000.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Z to be %f, got %f", pointIn000.Z, pointIn000center.Z)
	}

	// Voxel 0 0 0
	pointIn000 = Point{X: 0.0, Y: 0.0, Z: 0.0}
	if block000.getVoxelPtrByCoordinates(pointIn000) == nil {
		t.Errorf("Expected pointIn000 to not be nil")
	}
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	if voxelIndex != (IndexType{0, 0, 0}) {
		t.Errorf("Expected {0, 0, 0} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center.X, pointIn000.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.X to be %f, got %f", pointIn000.X, pointIn000center.X)
	}
	if !almostEqual(pointIn000center.Y, pointIn000.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Y to be %f, got %f", pointIn000.Y, pointIn000center.Y)
	}
	if !almostEqual(pointIn000center.Z, pointIn000.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Z to be %f, got %f", pointIn000.Z, pointIn000center.Z)
	}

	// Voxel 7 7 7
	pointIn000 = Point{
		X: 7.0 * tsdfMap.TsdfVoxelSize,
		Y: 7.0 * tsdfMap.TsdfVoxelSize,
		Z: 7.0 * tsdfMap.TsdfVoxelSize,
	}
	if block000.getVoxelPtrByCoordinates(pointIn000) == nil {
		t.Errorf("Expected pointIn000 to not be nil")
	}
	voxelIndex = block000.computeTruncatedVoxelIndexFromCoordinates(pointIn000)
	if voxelIndex != (IndexType{7, 7, 7}) {
		t.Errorf("Expected {7, 7, 7} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = block000.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center.X, pointIn000.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.X to be %f, got %f", pointIn000.X, pointIn000center.X)
	}
	if !almostEqual(pointIn000center.Y, pointIn000.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Y to be %f, got %f", pointIn000.Y, pointIn000center.Y)
	}
	if !almostEqual(pointIn000center.Z, pointIn000.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf("Expected pointIn000center.Z to be %f, got %f", pointIn000.Z, pointIn000center.Z)
	}

	// Block -1 -1 -1

	// Voxel 0 0 0
	pointInNeg111 = Point{
		X: -1.0 * blockNeg111.BlockSize,
		Y: -1.0 * blockNeg111.BlockSize,
		Z: -1.0 * blockNeg111.BlockSize,
	}
	if blockNeg111.getVoxelPtrByCoordinates(pointInNeg111) == nil {
		t.Errorf("Expected pointInNeg111 to not be nil")
	}

	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	if voxelIndex != (IndexType{0, 0, 0}) {
		t.Errorf("Expected {0, 0, 0} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn000center = blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn000center.X, pointInNeg111.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.X to be %f, got %f",
			pointInNeg111.X,
			pointIn000center.X,
		)
	}
	if !almostEqual(pointIn000center.Y, pointInNeg111.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.Y to be %f, got %f",
			pointInNeg111.Y,
			pointIn000center.Y,
		)
	}
	if !almostEqual(pointIn000center.Z, pointInNeg111.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.Z to be %f, got %f",
			pointInNeg111.Z,
			pointIn000center.Z,
		)
	}

	// Voxel 7 7 7
	pointInNeg111 = Point{X: -kEpsilon, Y: -kEpsilon, Z: -kEpsilon}
	if blockNeg111.getVoxelPtrByCoordinates(pointInNeg111) == nil {
		t.Errorf("Expected pointInNeg111 to not be nil")
	}
	voxelIndex = blockNeg111.computeTruncatedVoxelIndexFromCoordinates(pointInNeg111)
	if voxelIndex != (IndexType{7, 7, 7}) {
		t.Errorf("Expected {7, 7, 7} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointIn777center := blockNeg111.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointIn777center.X, pointInNeg111.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.X to be %f, got %f",
			pointInNeg111.X,
			pointIn000center.X,
		)
	}
	if !almostEqual(pointIn777center.Y, pointInNeg111.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.Y to be %f, got %f",
			pointInNeg111.Y,
			pointIn000center.Y,
		)
	}
	if !almostEqual(pointIn777center.Z, pointInNeg111.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointIn000center.Z to be %f, got %f",
			pointInNeg111.Z,
			pointIn000center.Z,
		)
	}

	// Block -1 -1 0

	// Voxel 3 6 5
	pointInNeg110 := Point{
		X: -5.0 * tsdfMap.TsdfVoxelSize,
		Y: -2.0 * tsdfMap.TsdfVoxelSize,
		Z: 5.0 * tsdfMap.TsdfVoxelSize,
	}
	blockNeg1Neg1Pos0 := tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(pointInNeg110)
	if blockNeg1Neg1Pos0 == nil {
		t.Errorf("Expected blockNeg1Neg1Pos0 to not be nil")
	}
	if tsdfMap.GetTsdfLayerPtr().
		computeBlockIndexFromCoordinates(pointInNeg110) !=
		(IndexType{-1, -1, 0}) {
		t.Errorf("Expected blockIndex to be {-1, -1, 0}")
	}
	if blockNeg1Neg1Pos0.getVoxelPtrByCoordinates(pointInNeg110) == nil {
		t.Errorf("Expected pointInNeg110 to not be nil")
	}
	voxelIndex = blockNeg1Neg1Pos0.computeTruncatedVoxelIndexFromCoordinates(pointInNeg110)
	if voxelIndex != (IndexType{3, 6, 5}) {
		t.Errorf("Expected {3, 6, 5} to be returned by computeTruncatedVoxelIndexFromCoordinates")
	}
	pointInNeg110center := blockNeg1Neg1Pos0.computeCoordinatesFromVoxelIndex(voxelIndex)
	if !almostEqual(pointInNeg110center.X, pointInNeg110.X, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointInNeg110center.X to be %f, got %f",
			pointInNeg110.X,
			pointInNeg110center.X,
		)
	}
	if !almostEqual(pointInNeg110center.Y, pointInNeg110.Y, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointInNeg110center.Y to be %f, got %f",
			pointInNeg110.Y,
			pointInNeg110center.Y,
		)
	}
	if !almostEqual(pointInNeg110center.Z, pointInNeg110.Z, tsdfMap.TsdfVoxelSize) {
		t.Errorf(
			"Expected pointInNeg111center.Z to be %f, got %f",
			pointInNeg110.Z,
			pointInNeg110center.Z,
		)
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
						"Expected blockIndex (%v) to be equal to blockIndexFromOrigin (%v)",
						blockIndex,
						blockIndexFromOrigin,
					)
				}
			}
		}
	}
}
