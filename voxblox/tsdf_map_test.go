package voxblox

import (
	"testing"
)

var blockSize float64
var tsdfMap *TsdfMap

func init() {
	var tsdfVoxelSize = 0.1
	var tsdfVoxelsPerSide int32 = 8
	blockSize = tsdfVoxelSize * float64(tsdfVoxelsPerSide)
	tsdfMap = NewTsdfMap(tsdfVoxelSize, tsdfVoxelsPerSide)
}

func TestTsdfMapBlockAllocation(t *testing.T) {
	// Layer should have no blocks by default
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 0 {
		t.Errorf("Expected no blocks in layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{x: 0, y: 0.15, z: 0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 1 {
		t.Errorf("Expected one block in layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{x: 0, y: 0.13, z: 0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 1 {
		t.Errorf("Expected one block in layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
	tsdfMap.GetTsdfLayerPtr().allocateNewBlockByCoordinates(Point{x: -10.0, y: 13.5, z: 20.0})
	if tsdfMap.GetTsdfLayerPtr().getNumberOfAllocatedBlocks() != 2 {
		t.Errorf("Expected two blocks in layer, got %d", len(tsdfMap.TsdfLayer.Blocks))
	}
}

func TestTsdfMapIndexLookups(t *testing.T) {
	// TODO: Break this into multiple tests

	// BLOCK 0 0 0 (coordinate at origin)
	pointIn000 := Point{x: 0, y: 0, z: 0}
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
	pointIn000v2 := Point{x: 0.0, y: tsdfMap.TsdfVoxelSize, z: 0.0}
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
	pointIn111 := Point{x: blockSize, y: blockSize, z: blockSize}
	block111 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointIn111)
	if block111 == nil {
		t.Errorf("Expected block111 to not be nil")
	}
	index111 := IndexType{1, 1, 1}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointIn111) != index111 {
		t.Errorf("Expected {1, 1, 1} to be returned by GetBlockPtrByCoordinates")
	}
	if block111.Origin != (Point{x: blockSize, y: blockSize, z: blockSize}) {
		t.Errorf("Expected block111.Origin to be {%f, %f, %f}, got %v",
			blockSize, blockSize, blockSize, block111.Origin)
	}
	if block111.BlockIndex != index111 {
		t.Errorf("Expected block111.BlockIndex to be {1, 1, 1}, got %v", block111.BlockIndex)
	}

	// BLOCK 1 1 1 (coordinate within block)
	pointIn111v2 := Point{x: blockSize, y: blockSize + tsdfMap.TsdfVoxelSize, z: blockSize}
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
	if block111v2.Origin != (Point{x: blockSize, y: blockSize, z: blockSize}) {
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
	pointInNeg111 := Point{x: -blockSize, y: -blockSize, z: -blockSize}
	blockNeg111 := tsdfMap.GetTsdfLayerPtr().getBlockPtrByCoordinates(pointInNeg111)
	if blockNeg111 == nil {
		t.Errorf("Expected blockNeg111 to not be nil")
	}
	indexNeg111 := IndexType{-1, -1, -1}
	if tsdfMap.GetTsdfLayerPtr().computeBlockIndexFromCoordinates(pointInNeg111) != indexNeg111 {
		t.Errorf("Expected {0, 0, 0} to be returned by GetBlockPtrByCoordinates")
	}
	if blockNeg111.Origin != (Point{x: -blockSize, y: -blockSize, z: -blockSize}) {
		t.Errorf("Expected blockNeg111.Origin to be {%f, %f, %f}, got %v",
			-blockSize, -blockSize, -blockSize, blockNeg111.Origin)
	}
	if blockNeg111.BlockSize != blockSize {
		t.Errorf("Expected blockNeg111.BlockSize to be %f, got %f", blockSize, blockNeg111.BlockSize)
	}

	// BLOCK -1 -1 -1 (coordinate within block)
	pointInNeg111v2 := Point{x: -blockSize, y: -blockSize + tsdfMap.TsdfVoxelSize, z: -blockSize}
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
	if blockNeg111v2.Origin != (Point{x: -blockSize, y: -blockSize, z: -blockSize}) {
		t.Errorf("Expected blockNeg111.Origin to be {%f, %f, %f}, got %v",
			-blockSize, -blockSize, -blockSize, blockNeg111.Origin)
	}
	if blockNeg111v2.BlockSize != blockSize {
		t.Errorf("Expected blockNeg111.BlockSize to be %f, got %f", blockSize, blockNeg111v2.BlockSize)
	}
}

func TestVoxelIndexing(t *testing.T) {
	// Test Voxel 0 1 0
	
}
