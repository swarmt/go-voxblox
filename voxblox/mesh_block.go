package voxblox

import "sync"

type MeshBlock struct {
	Index         IndexType
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	mutex         sync.RWMutex
}

// NewMeshBlock creates a new MeshBlock.
func NewMeshBlock(voxelsPerSide int, voxelSize float64, index IndexType, origin Point) *MeshBlock {
	b := new(MeshBlock)
	b.VoxelsPerSide = voxelsPerSide
	b.VoxelSize = voxelSize
	b.Origin = origin
	b.Index = index
	b.VoxelSizeInv = 1.0 / voxelSize
	b.BlockSize = float64(voxelsPerSide) * voxelSize
	b.BlockSizeInv = 1.0 / b.BlockSize
	return b
}
