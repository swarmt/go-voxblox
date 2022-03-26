package voxblox

import (
	"fmt"
	"sync"
)

type MeshBlock struct {
	Index         IndexType
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	sync.RWMutex
	vertices  []Point
	triangles [][3]int
	colors    []Color
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

// String returns a string representation of the MeshBlock.
func (b *MeshBlock) String() string {
	return fmt.Sprintf("%d_%d_%d", b.Index[0], b.Index[1], b.Index[2])
}

// getVertexCount returns the number of vertices in the block.
// Thread-safe.
func (b *MeshBlock) getVertexCount() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.vertices)
}

// getVertices returns the vertices in the block.
// Thread-safe.
func (b *MeshBlock) getVertices() []Point {
	b.RLock()
	defer b.RUnlock()
	return b.vertices
}
