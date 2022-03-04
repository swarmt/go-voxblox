package voxblox

type TsdfLayer struct {
	VoxelSize     float64
	VoxelsPerSide int
	Blocks        map[IndexType]*Block
	BlockSize     float64
	BlockSizeInv  float64
}

func NewTsdfLayer(voxelSize float64, voxelsPerSide int) *TsdfLayer {
	l := new(TsdfLayer)
	l.VoxelSize = voxelSize
	l.VoxelsPerSide = voxelsPerSide
	l.Blocks = make(map[IndexType]*Block)
	l.BlockSize = voxelSize * float64(voxelsPerSide)
	l.BlockSizeInv = 1.0 / l.BlockSize
	return l
}

// getNumberOfAllocatedBlocks returns the number of blocks allocated in the map
func (l *TsdfLayer) getNumberOfAllocatedBlocks() int {
	return len(l.Blocks)
}

// getBlock allocates a new block in the map
func (l *TsdfLayer) getBlock(blockIndex IndexType) *Block {
	// Test if block already exists
	if block, ok := l.Blocks[blockIndex]; ok {
		return block
	}
	newBlock := NewBlock(
		l.VoxelsPerSide,
		l.VoxelSize,
		blockIndex,
		getOriginPointFromGridIndex(blockIndex, l.BlockSize),
	)
	l.Blocks[blockIndex] = newBlock
	return newBlock
}

// allocateNewBlockByCoordinates allocates a new block in the map by coordinates
// TODO: This and getBlockPtrByCoordinates should be merged as they are interchangeable
func (l *TsdfLayer) allocateNewBlockByCoordinates(point Point) *Block {
	return l.getBlock(getGridIndexFromPoint(point, l.BlockSizeInv))
}

// computeBlockIndexFromCoordinates computes the block index from coordinates
func (l *TsdfLayer) computeBlockIndexFromCoordinates(point Point) IndexType {
	return getGridIndexFromPoint(point, l.BlockSizeInv)
}

// getBlockPtrByCoordinates returns a pointer to the block in the map by coordinates
func (l *TsdfLayer) getBlockPtrByCoordinates(point Point) *Block {
	return l.getBlockPtrByIndex(l.computeBlockIndexFromCoordinates(point))
}

// getBlockPtrByIndex returns a pointer to the block in the map by index
func (l *TsdfLayer) getBlockPtrByIndex(index IndexType) *Block {
	block, ok := l.Blocks[index]
	if !ok {
		return l.getBlock(index)
	}
	return block
}

func (l *TsdfLayer) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	block := l.getBlockPtrByIndex(l.computeBlockIndexFromCoordinates(point))
	if block == nil {
		return nil
	}
	return l.getVoxelPtrByCoordinates(point)
}
