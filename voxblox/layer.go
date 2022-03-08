package voxblox

import "sync"

type TsdfLayer struct {
	voxelSize        float64
	voxelSizeInv     float64
	voxelsPerSide    int
	voxelsPerSideInv float64
	blocks           map[IndexType]*Block
	blockSize        float64
	blockSizeInv     float64
	mutex            sync.RWMutex
}

func NewTsdfLayer(voxelSize float64, voxelsPerSide int) *TsdfLayer {
	l := new(TsdfLayer)
	l.voxelSize = voxelSize
	l.voxelsPerSide = voxelsPerSide
	l.voxelSizeInv = 1.0 / voxelSize
	l.voxelsPerSideInv = 1.0 / float64(voxelsPerSide)
	l.blocks = make(map[IndexType]*Block)
	l.blockSize = voxelSize * float64(voxelsPerSide)
	l.blockSizeInv = 1.0 / l.blockSize
	return l
}

func (l *TsdfLayer) GetVoxelsPerSide() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelsPerSide
}

func (l *TsdfLayer) GetVoxelsPerSideInv() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelsPerSideInv
}

func (l *TsdfLayer) GetVoxelSize() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelSize
}

func (l *TsdfLayer) GetVoxelSizeInv() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelSizeInv
}

func (l *TsdfLayer) GetBlockSize() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.blockSize
}

// getNumberOfAllocatedBlocks returns the number of blocks allocated in the map
func (l *TsdfLayer) getNumberOfAllocatedBlocks() int {
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
		l.GetVoxelsPerSide(),
		l.GetVoxelSize(),
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.GetBlockSize()),
	)
	l.mutex.Lock()
	l.blocks[blockIndex] = newBlock
	l.mutex.Unlock()
	return newBlock
}

// allocateNewBlockByCoordinates allocates a new block in the map by coordinates
// TODO: This and getBlockPtrByCoordinates should be merged as they are interchangeable
func (l *TsdfLayer) allocateNewBlockByCoordinates(point Point) *Block {
	return l.getBlock(getGridIndexFromPoint(point, l.blockSizeInv))
}

// computeBlockIndexFromCoordinates computes the block index from coordinates
func (l *TsdfLayer) computeBlockIndexFromCoordinates(point Point) IndexType {
	return getGridIndexFromPoint(point, l.blockSizeInv)
}

// getBlockPtrByCoordinates returns a pointer to the block in the map by coordinates
func (l *TsdfLayer) getBlockPtrByCoordinates(point Point) *Block {
	return l.getBlock(l.computeBlockIndexFromCoordinates(point))
}

func (l *TsdfLayer) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	block := l.getBlock(l.computeBlockIndexFromCoordinates(point))
	if block == nil {
		return nil
	}
	return l.getVoxelPtrByCoordinates(point)
}

func allocateStorageAndGetVoxelPtr(layer *TsdfLayer, globalVoxelIndex IndexType) *TsdfVoxel {
	blockIndex := getBlockIndexFromGlobalVoxelIndex(globalVoxelIndex, layer.GetVoxelsPerSideInv())
	block := layer.getBlock(blockIndex)
	voxelIndex := getLocalFromGlobalVoxelIndex(globalVoxelIndex, blockIndex, layer.GetVoxelsPerSide())
	voxel := block.getVoxel(voxelIndex)
	return voxel
}
