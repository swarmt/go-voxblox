package voxblox

import "sync"

type Mesh struct {
	Vertices []Point
	Indices  []uint32
}

type MeshBlock struct {
	Index         IndexType
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	sync.RWMutex
	mesh Mesh
}

// NewMeshBlock creates a new MeshBlock.
func NewMeshBlock(layer *MeshLayer, index IndexType, origin Point) *MeshBlock {
	b := new(MeshBlock)
	b.Origin = origin
	b.Index = index
	b.VoxelsPerSide = layer.VoxelsPerSide
	b.VoxelSize = layer.VoxelSize
	b.VoxelSizeInv = layer.VoxelSizeInv
	b.BlockSize = layer.BlockSize
	b.BlockSizeInv = layer.BlockSizeInv
	return b
}
