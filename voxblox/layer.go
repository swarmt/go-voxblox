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

func (l *TsdfLayer) getBlockSizeInv() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.blockSizeInv
}

func (l *TsdfLayer) getVoxelsPerSide() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelsPerSide
}

func (l *TsdfLayer) getVoxelsPerSideInv() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelsPerSideInv
}

func (l *TsdfLayer) getVoxelSize() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelSize
}

func (l *TsdfLayer) getVoxelSizeInv() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.voxelSizeInv
}

func (l *TsdfLayer) getBlockSize() float64 {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.blockSize
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
		l.getVoxelsPerSide(),
		l.getVoxelSize(),
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.getBlockSize()),
	)
	l.mutex.Lock()
	l.blocks[blockIndex] = newBlock
	l.mutex.Unlock()
	return newBlock
}

// computeBlockIndexFromCoordinates computes the block index from coordinates
func (l *TsdfLayer) computeBlockIndexFromCoordinates(point Point) IndexType {
	return getGridIndexFromPoint(point, l.getBlockSizeInv())
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
	blockIndex := getBlockIndexFromGlobalVoxelIndex(globalVoxelIndex, layer.getVoxelsPerSideInv())
	block := layer.getBlock(blockIndex)
	voxelIndex := getLocalFromGlobalVoxelIndex(globalVoxelIndex, blockIndex, layer.getVoxelsPerSide())
	voxel := block.getVoxel(voxelIndex)
	return voxel
}
