package voxblox

type Layer struct {
	VoxelSize     float64
	VoxelsPerSide int32
	Blocks        map[IndexType]*Block
	BlockSize     float64
	BlockSizeInv  float64
}

func NewLayer(voxelSize float64, voxelsPerSide int32) *Layer {
	l := new(Layer)
	l.VoxelSize = voxelSize
	l.VoxelsPerSide = voxelsPerSide
	l.Blocks = make(map[IndexType]*Block)
	l.BlockSize = voxelSize * float64(voxelsPerSide)
	l.BlockSizeInv = 1.0 / l.BlockSize
	return l
}

// getNumberOfAllocatedBlocks returns the number of blocks allocated in the map
func (l *Layer) getNumberOfAllocatedBlocks() int {
	return len(l.Blocks)
}

// allocateNewBlock allocates a new block in the map
func (l *Layer) allocateNewBlock(blockIndex IndexType) *Block {
	// Test if block already exists
	if _, ok := l.Blocks[blockIndex]; ok {
		return l.Blocks[blockIndex]
	}
	newBlock := NewBlock(l.VoxelsPerSide, l.VoxelSize, blockIndex, getOriginPointFromGridIndex(blockIndex, l.BlockSize))
	l.Blocks[blockIndex] = newBlock
	return newBlock
}

// allocateNewBlockByCoordinates allocates a new block in the map by coordinates
// TODO: This and getBlockPtrByCoordinates should be merged as they are interchangeable
func (l *Layer) allocateNewBlockByCoordinates(point Point) *Block {
	return l.allocateNewBlock(getGridIndexFromPoint(point, l.BlockSizeInv))
}

// computeBlockIndexFromCoordinates computes the block index from coordinates
func (l *Layer) computeBlockIndexFromCoordinates(point Point) IndexType {
	return getGridIndexFromPoint(point, l.BlockSizeInv)
}

// getBlockPtrByCoordinates returns a pointer to the block in the map by coordinates
func (l *Layer) getBlockPtrByCoordinates(point Point) *Block {
	return l.getBlockPtrByIndex(l.computeBlockIndexFromCoordinates(point))
}

// getBlockPtrByIndex returns a pointer to the block in the map by index
func (l *Layer) getBlockPtrByIndex(index IndexType) *Block {
	block, ok := l.Blocks[index]
	if !ok {
		return l.allocateNewBlock(index)
	}
	return block
}
