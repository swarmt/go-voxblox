package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
	"sync"
)

// Block contains a map of voxels.
type Block struct {
	hasData       bool
	voxelsPerSide int
	voxelSize     float64
	origin        Point
	index         IndexType
	updated       bool
	numVoxels     int
	voxelSizeInv  float64
	blockSize     float64
	blockSizeInv  float64
	voxels        map[IndexType]*TsdfVoxel
	mutex         sync.RWMutex
}

// NewBlock creates a new Block.
func NewBlock(voxelsPerSide int, voxelSize float64, index IndexType, origin Point) *Block {
	b := new(Block)
	b.hasData = false
	b.voxelsPerSide = voxelsPerSide
	b.voxelSize = voxelSize
	b.origin = origin
	b.index = index
	b.updated = false
	b.numVoxels = voxelsPerSide * voxelsPerSide * voxelsPerSide
	b.voxelSizeInv = 1.0 / voxelSize
	b.blockSize = float64(voxelsPerSide) * voxelSize
	b.blockSizeInv = 1.0 / b.blockSize
	b.voxels = make(map[IndexType]*TsdfVoxel)
	return b
}

// getOrigin returns a reference the origin of the block.
// Thread-safe.
func (b *Block) getOrigin() Point {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.origin
}

// getVoxelSize returns the voxel size for the Block.
// Thread-safe.
func (b *Block) getVoxelSize() float64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxelSize
}

// getVoxelSizeInv returns the Inverse of the voxel size for the Block.
// Thread-safe.
func (b *Block) getVoxelSizeInv() float64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxelSizeInv
}

// getVoxelsPerSide returns the number of voxels per side for the Block.
// Thread-safe.
func (b *Block) getVoxelsPerSide() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxelsPerSide
}

// getVoxels returns a copy of the map of voxels.
// Thread-safe.
func (b *Block) getVoxels() map[IndexType]*TsdfVoxel {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxels
}

// getVoxel returns a reference to a voxel at the given index .
// Creates a new voxel if it doesn't exist.
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
	b.mutex.Lock()
	b.voxels[voxelIndex] = newVoxel
	b.mutex.Unlock()
	return newVoxel
}

// getVoxelPtrByCoordinates returns a reference to a voxel at the given coordinates.
// Creates a new voxel if it does not exist.
func (b *Block) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	return b.getVoxel(getGridIndexFromPoint(point, b.getVoxelSize()))
}

// computeTruncatedVoxelIndexFromCoordinates
// Computes the truncated voxel index from the given coordinates.
func (b *Block) computeTruncatedVoxelIndexFromCoordinates(point Point) IndexType {
	maxValue := b.getVoxelsPerSide() - 1
	origin := b.getOrigin()
	voxelIndex := getGridIndexFromPoint(vec3.Sub(&point, &origin), b.getVoxelSizeInv())
	index := IndexType{
		MaxInt(MinInt(voxelIndex[0], maxValue), 0.0),
		MaxInt(MinInt(voxelIndex[1], maxValue), 0.0),
		MaxInt(MinInt(voxelIndex[2], maxValue), 0.0),
	}
	return b.getVoxel(index).getIndex()
}

// computeCoordinatesFromVoxelIndex
// Computes the coordinates (Voxel center) from the given truncated voxel index.
func (b *Block) computeCoordinatesFromVoxelIndex(index IndexType) Point {
	centerPoint := getCenterPointFromGridIndex(index, b.getVoxelSize())
	origin := b.getOrigin()
	return vec3.Add(&origin, &centerPoint)
}
