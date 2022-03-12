package voxblox

import "sync"

type MeshLayer struct {
	VoxelSize        float64
	VoxelSizeInv     float64
	VoxelsPerSide    int
	VoxelsPerSideInv float64
	BlockSize        float64
	BlockSizeInv     float64
	blocks           map[IndexType]*MeshBlock
	mutex            sync.RWMutex
}

func NewMeshLayer(tsdfLayer *TsdfLayer) *MeshLayer {
	meshLayer := MeshLayer{
		VoxelSize:        tsdfLayer.VoxelSize,
		VoxelSizeInv:     tsdfLayer.VoxelSizeInv,
		VoxelsPerSide:    tsdfLayer.VoxelsPerSide,
		VoxelsPerSideInv: tsdfLayer.VoxelsPerSideInv,
		BlockSize:        tsdfLayer.BlockSize,
		BlockSizeInv:     tsdfLayer.BlockSizeInv,
		mutex:            sync.RWMutex{},
		blocks:           make(map[IndexType]*MeshBlock),
	}
	return &meshLayer
}

// getBlockByIndex allocates a new block in the map or returns an existing one
// TODO: Would this be better as an interface shared by TsdfLayer?
func (l *MeshLayer) getBlockByIndex(blockIndex IndexType) *MeshBlock {
	// Test if block already exists
	l.mutex.RLock()
	block, ok := l.blocks[blockIndex]
	l.mutex.RUnlock()
	if ok {
		return block
	}
	newBlock := NewMeshBlock(
		l,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.mutex.Lock()
	l.blocks[blockIndex] = newBlock
	l.mutex.Unlock()
	return newBlock
}

// getBlockByCoordinates returns a pointer to the block by coordinates
func (l *MeshLayer) getBlockByCoordinates(point Point) *MeshBlock {
	return l.getBlockByIndex(getBlockIndexFromCoordinates(point, l.BlockSizeInv))
}
