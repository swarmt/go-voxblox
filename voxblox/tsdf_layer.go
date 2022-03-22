package voxblox

import (
	"math"
	"sync"
)

type TsdfLayer struct {
	VoxelSize        float64
	VoxelSizeInv     float64
	VoxelsPerSide    int
	VoxelsPerSideInv float64
	BlockSize        float64
	BlockSizeInv     float64
	sync.RWMutex
	blocks map[IndexType]*TsdfBlock
}

// NewTsdfLayer creates a new TsdfLayer.
// Computes inverse variables for faster access.
func NewTsdfLayer(voxelSize float64, voxelsPerSide int) *TsdfLayer {
	l := new(TsdfLayer)
	l.VoxelSize = voxelSize
	l.VoxelsPerSide = voxelsPerSide
	l.VoxelSizeInv = 1.0 / voxelSize
	l.VoxelsPerSideInv = 1.0 / float64(voxelsPerSide)
	l.BlockSize = voxelSize * float64(voxelsPerSide)
	l.BlockSizeInv = 1.0 / l.BlockSize
	l.blocks = make(map[IndexType]*TsdfBlock)
	return l
}

// getBlocks returns a copy of the map of blocks
// Thread-safe.
func (l *TsdfLayer) getBlocks() map[IndexType]*TsdfBlock {
	l.RLock()
	defer l.RUnlock()
	return l.blocks
}

// getUpdatedBlocks returns a map of references to TsdfBlocks that have been updated
// Thread-safe.
func (l *TsdfLayer) getUpdatedBlocks() map[IndexType]*TsdfBlock {
	l.RLock()
	defer l.RUnlock()
	updatedBlocks := make(map[IndexType]*TsdfBlock)
	for index, block := range l.blocks {
		if block.getUpdated() {
			updatedBlocks[index] = block
		}
	}
	return updatedBlocks
}

// getVoxelCenters returns all voxel centers (global coordinates) in the Layer close to the surface.
// Thread-safe.
func (l *TsdfLayer) getVoxelCenters() ([]Point, []Color) {
	var voxelCenters []Point
	var voxelColors []Color
	l.RLock()
	defer l.RUnlock()
	for _, block := range l.getBlocks() {
		for _, voxel := range block.getVoxels() {
			if math.Abs(voxel.getDistance()) < block.VoxelSize && voxel.getWeight() > 2 {
				coordinates := block.computeCoordinatesFromVoxelIndex(voxel.Index)
				voxelCenters = append(voxelCenters, coordinates)
				color := voxel.getColor()
				voxelColors = append(voxelColors, color)
			}
		}
	}
	return voxelCenters, voxelColors
}

// getBlockCount returns the number of blocks allocated in the map
// Thread-safe.
func (l *TsdfLayer) getBlockCount() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.blocks)
}

// getBlockByIndex allocates a new block in the map or returns an existing one
// Thread-safe.
func (l *TsdfLayer) getBlockByIndex(blockIndex IndexType) *TsdfBlock {
	// Test if block already exists
	l.RLock()
	block, ok := l.blocks[blockIndex]
	l.RUnlock()
	if ok {
		return block
	}
	newBlock := NewTsdfBlock(
		l,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.Lock()
	l.blocks[blockIndex] = newBlock
	l.Unlock()
	return newBlock
}

// getBlockByCoordinates returns a pointer to the block by coordinates
func (l *TsdfLayer) getBlockByCoordinates(point Point) *TsdfBlock {
	return l.getBlockByIndex(getBlockIndexFromCoordinates(point, l.BlockSizeInv))
}

// getBlockIfExists returns a pointer to the block if it exists
// Thread-safe.
func (l *TsdfLayer) getBlockIfExists(index IndexType) *TsdfBlock {
	l.RLock()
	defer l.RUnlock()
	block, ok := l.blocks[index]
	if ok {
		return block
	}
	return nil
}

// getBlockAndVoxelFromGlobalVoxelIndex allocates a new block in the map and returns the block and voxel
func getBlockAndVoxelFromGlobalVoxelIndex(
	layer *TsdfLayer,
	globalVoxelIndex IndexType,
) (*TsdfBlock, *TsdfVoxel) {
	blockIndex := getBlockIndexFromGlobalVoxelIndex(globalVoxelIndex, layer.VoxelsPerSideInv)
	block := layer.getBlockByIndex(blockIndex)
	voxelIndex := getLocalFromGlobalVoxelIndex(globalVoxelIndex, blockIndex, layer.VoxelsPerSide)
	voxel := block.getVoxel(voxelIndex)
	return block, voxel
}
