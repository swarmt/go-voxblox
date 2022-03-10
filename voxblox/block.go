package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"sync"
)

// Block contains a map of voxels.
type Block struct {
	hasData       bool
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	index         IndexType
	updated       bool
	numVoxels     int
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	voxels        map[IndexType]*TsdfVoxel
	mutex         sync.RWMutex
}

// NewBlock creates a new Block.
func NewBlock(voxelsPerSide int, voxelSize float64, index IndexType, origin Point) *Block {
	b := new(Block)
	b.hasData = false
	b.VoxelsPerSide = voxelsPerSide
	b.VoxelSize = voxelSize
	b.Origin = origin
	b.index = index
	b.updated = false
	b.numVoxels = voxelsPerSide * voxelsPerSide * voxelsPerSide
	b.VoxelSizeInv = 1.0 / voxelSize
	b.BlockSize = float64(voxelsPerSide) * voxelSize
	b.BlockSizeInv = 1.0 / b.BlockSize
	b.voxels = make(map[IndexType]*TsdfVoxel)
	return b
}

// getVoxels returns a copy of the map of voxels.
// Thread-safe.
func (b *Block) getVoxels() map[IndexType]*TsdfVoxel {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxels
}

// addVoxel adds a voxel to the block.
// Thread-safe.
func (b *Block) addVoxel(voxel *TsdfVoxel) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.voxels[voxel.Index] = voxel
}

// getVoxel returns a reference to a voxel at the given index .
// Creates a new voxel if it doesn't exist.
// Thread-safe.
func (b *Block) getVoxel(voxelIndex IndexType) *TsdfVoxel {
	// Test if voxel already exists
	b.mutex.RLock()
	voxel, ok := b.voxels[voxelIndex]
	b.mutex.RUnlock()
	if ok {
		return voxel
	}
	// Create a new voxel
	newVoxel := NewVoxel(voxelIndex)
	b.addVoxel(newVoxel)
	return newVoxel
}

// getVoxelPtrByCoordinates returns a reference to a voxel at the given coordinates.
// Creates a new voxel if it does not exist.
// Thread-safe.
func (b *Block) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	return b.getVoxel(getGridIndexFromPoint(point, b.VoxelSize))
}

// computeTruncatedVoxelIndexFromCoordinates
// Computes the truncated voxel index from the given coordinates.
func (b *Block) computeTruncatedVoxelIndexFromCoordinates(point Point) IndexType {
	maxValue := b.VoxelsPerSide - 1
	voxelIndex := getGridIndexFromPoint(vec3.Sub(&point, &b.Origin), b.VoxelSizeInv)
	index := IndexType{
		MaxInt(MinInt(voxelIndex[0], maxValue), 0.0),
		MaxInt(MinInt(voxelIndex[1], maxValue), 0.0),
		MaxInt(MinInt(voxelIndex[2], maxValue), 0.0),
	}
	return b.getVoxel(index).Index
}

// computeCoordinatesFromVoxelIndex
// Computes the coordinates (Voxel center) from the given truncated voxel index.
func (b *Block) computeCoordinatesFromVoxelIndex(index IndexType) Point {
	centerPoint := getCenterPointFromGridIndex(index, b.VoxelSize)
	return vec3.Add(&b.Origin, &centerPoint)
}
