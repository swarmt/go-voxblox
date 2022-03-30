package voxblox

import "sync"

type MeshLayer struct {
	VoxelSize        float64
	VoxelSizeInv     float64
	VoxelsPerSide    int
	VoxelsPerSideInv float64
	BlockSize        float64
	BlockSizeInv     float64
	sync.RWMutex
	blocks map[IndexType]*MeshBlock
}

func NewMeshLayer(tsdfLayer *TsdfLayer) *MeshLayer {
	meshLayer := MeshLayer{
		VoxelSize:        tsdfLayer.VoxelSize,
		VoxelSizeInv:     tsdfLayer.VoxelSizeInv,
		VoxelsPerSide:    tsdfLayer.VoxelsPerSide,
		VoxelsPerSideInv: tsdfLayer.VoxelsPerSideInv,
		BlockSize:        tsdfLayer.BlockSize,
		BlockSizeInv:     tsdfLayer.BlockSizeInv,
		blocks:           make(map[IndexType]*MeshBlock),
	}
	return &meshLayer
}

// getBlockCount returns the number of blocks allocated in the map
func (l *MeshLayer) getBlockCount() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.blocks)
}

// GetBlocks returns the blocks in the map
// Thread-safe.
func (l *MeshLayer) GetBlocks() map[IndexType]*MeshBlock {
	l.RLock()
	defer l.RUnlock()
	return l.blocks
}

// getBlockByIndex allocates a new block in the map or returns an existing one
func (l *MeshLayer) getBlockByIndex(blockIndex IndexType) *MeshBlock {
	// Test if block already exists
	l.RLock()
	block, ok := l.blocks[blockIndex]
	l.RUnlock()
	if ok {
		return block
	}
	newBlock := NewMeshBlock(
		l,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.Lock()
	l.blocks[blockIndex] = newBlock
	l.Unlock()
	return newBlock
}

// getNewBlockByIndex allocates a new block in the map and returns it
// Overwrites any existing block
func (l *MeshLayer) getNewBlockByIndex(blockIndex IndexType) *MeshBlock {
	l.Lock()
	defer l.Unlock()
	newBlock := NewMeshBlock(
		l,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.blocks[blockIndex] = newBlock
	return newBlock
}

// getBlockIfExists returns a pointer to the block if it exists
// Thread-safe.
func (l *MeshLayer) getBlockIfExists(index IndexType) *MeshBlock {
	l.RLock()
	defer l.RUnlock()
	block, ok := l.blocks[index]
	if ok {
		return block
	}
	return nil
}

// getBlockByCoordinates returns a pointer to the block by coordinates
func (l *MeshLayer) getBlockByCoordinates(point Point) *MeshBlock {
	return l.getBlockByIndex(getBlockIndexFromCoordinates(point, l.BlockSizeInv))
}
