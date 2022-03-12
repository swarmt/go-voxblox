package voxblox

import (
	"sync"

	"github.com/ungerik/go3d/float64/vec3"
)

// TsdfBlock contains a map of voxels.
type TsdfBlock struct {
	Index         IndexType
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	mutex         sync.RWMutex
	updated       bool
	voxels        map[IndexType]*TsdfVoxel
}

// NewTsdfBlock creates a new TsdfBlock.
func NewTsdfBlock(layer *TsdfLayer, index IndexType, origin Point) *TsdfBlock {
	b := new(TsdfBlock)
	b.Origin = origin
	b.Index = index
	b.VoxelsPerSide = layer.VoxelsPerSide
	b.VoxelSize = layer.VoxelSize
	b.VoxelSizeInv = layer.VoxelSizeInv
	b.BlockSize = layer.BlockSize
	b.BlockSizeInv = layer.BlockSizeInv
	b.mutex = sync.RWMutex{}
	b.updated = true
	b.voxels = make(map[IndexType]*TsdfVoxel)
	return b
}

// getVoxels returns a copy of the map of voxels.
// Thread-safe.
func (b *TsdfBlock) getVoxels() map[IndexType]*TsdfVoxel {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.voxels
}

// addVoxel adds a voxel to the block.
// Thread-safe.
func (b *TsdfBlock) addVoxel(voxel *TsdfVoxel) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.voxels[voxel.Index] = voxel
}

// getUpdated gets the updated flag.
// Thread-safe.
func (b *TsdfBlock) getUpdated() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.updated
}

// setUpdated sets the updated flag.
// Thread-safe.
func (b *TsdfBlock) setUpdated() {
	// Avoid getting a mutex write lock if we don't need to.
	if !b.getUpdated() {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		b.updated = true
	}
}

// getVoxel returns a reference to a voxel at the given Index .
// Creates a new voxel if it doesn't exist.
// Thread-safe.
func (b *TsdfBlock) getVoxel(voxelIndex IndexType) *TsdfVoxel {
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
func (b *TsdfBlock) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	return b.getVoxel(getGridIndexFromPoint(point, b.VoxelSize))
}

// computeTruncatedVoxelIndexFromCoordinates
// Computes the truncated voxel Index from the given coordinates.
func (b *TsdfBlock) computeTruncatedVoxelIndexFromCoordinates(point Point) IndexType {
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
// Computes the coordinates (Voxel center) from the given truncated voxel Index.
func (b *TsdfBlock) computeCoordinatesFromVoxelIndex(index IndexType) Point {
	centerPoint := getCenterPointFromGridIndex(index, b.VoxelSize)
	return vec3.Add(&b.Origin, &centerPoint)
}
