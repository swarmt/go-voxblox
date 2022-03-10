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
	blocks           map[IndexType]*Block
	BlockSize        float64
	BlockSizeInv     float64
	mutex            sync.RWMutex
}

func NewTsdfLayer(voxelSize float64, voxelsPerSide int) *TsdfLayer {
	l := new(TsdfLayer)
	l.VoxelSize = voxelSize
	l.VoxelsPerSide = voxelsPerSide
	l.VoxelSizeInv = 1.0 / voxelSize
	l.VoxelsPerSideInv = 1.0 / float64(voxelsPerSide)
	l.blocks = make(map[IndexType]*Block)
	l.BlockSize = voxelSize * float64(voxelsPerSide)
	l.BlockSizeInv = 1.0 / l.BlockSize
	return l
}

// getBlocks returns a copy of the map of blocks
func (l *TsdfLayer) getBlocks() map[IndexType]*Block {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.blocks
}

// getVoxelCenters returns all voxel centers (global coordinates) in the layer close to the surface.
// Thread-safe.
func (l *TsdfLayer) getVoxelCenters() []Point {
	var voxelCenters []Point
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for _, block := range l.getBlocks() {
		for _, voxel := range block.getVoxels() {
			if math.Abs(voxel.getDistance()) < block.VoxelSize {
				coordinates := block.computeCoordinatesFromVoxelIndex(voxel.getIndex())
				voxelCenters = append(voxelCenters, coordinates)
			}
		}
	}
	return voxelCenters
}

// getNumberOfAllocatedBlocks returns the number of blocks allocated in the map
func (l *TsdfLayer) getNumberOfAllocatedBlocks() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return len(l.blocks)
}

// getBlock allocates a new block in the map or returns an existing one
func (l *TsdfLayer) getBlock(blockIndex IndexType) *Block {
	// Test if block already exists
	l.mutex.RLock()
	block, ok := l.blocks[blockIndex]
	l.mutex.RUnlock()
	if ok {
		return block
	}
	newBlock := NewBlock(
		l.VoxelsPerSide,
		l.VoxelSize,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.mutex.Lock()
	l.blocks[blockIndex] = newBlock
	l.mutex.Unlock()
	return newBlock
}

// computeBlockIndexFromCoordinates computes the block index from coordinates
func (l *TsdfLayer) computeBlockIndexFromCoordinates(point Point) IndexType {
	return getGridIndexFromPoint(point, l.BlockSizeInv)
}

// getBlockPtrByCoordinates returns a pointer to the block in the map by coordinates
func (l *TsdfLayer) getBlockPtrByCoordinates(point Point) *Block {
	return l.getBlock(l.computeBlockIndexFromCoordinates(point))
}

// getVoxelPtrByCoordinates returns a pointer to the voxel in the block in the map by coordinates
func (l *TsdfLayer) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	block := l.getBlock(l.computeBlockIndexFromCoordinates(point))
	if block == nil {
		return nil
	}
	return l.getVoxelPtrByCoordinates(point)
}

// allocateStorageAndGetVoxelPtr allocates a new block in the map and returns a pointer to the voxel
func allocateStorageAndGetVoxelPtr(layer *TsdfLayer, globalVoxelIndex IndexType) *TsdfVoxel {
	blockIndex := getBlockIndexFromGlobalVoxelIndex(globalVoxelIndex, layer.VoxelsPerSideInv)
	block := layer.getBlock(blockIndex)
	voxelIndex := getLocalFromGlobalVoxelIndex(globalVoxelIndex, blockIndex, layer.VoxelsPerSide)
	voxel := block.getVoxel(voxelIndex)
	return voxel
}
