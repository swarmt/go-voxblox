package voxblox

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
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

// Clear clears the block.
func (b *MeshBlock) Clear() {
	b.Lock()
	defer b.Unlock()
	b.vertices = nil
	b.triangles = nil
	b.colors = nil
}

// getVertexCount returns the number of vertices in the block.
// Thread-safe.
func (b *MeshBlock) getVertexCount() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.vertices)
}

// HasData returns whether the block has data.
// Thread-safe.
func (b *MeshBlock) HasData() bool {
	b.RLock()
	defer b.RUnlock()
	return len(b.vertices) > 0 && len(b.triangles) > 0
}

// getVertices returns the vertices in the block.
// Thread-safe.
func (b *MeshBlock) getVertices() []Point {
	b.RLock()
	defer b.RUnlock()
	return b.vertices
}

// verticesAsFloat32 returns the vertices in the block as float32.
// Thread-safe.
func (b *MeshBlock) verticesAsFloat32() [][3]float32 {
	b.RLock()
	defer b.RUnlock()
	vertices := make([][3]float32, len(b.vertices))
	for i, v := range b.vertices {
		vertices[i][0] = float32(v[0])
		vertices[i][1] = float32(v[1])
		vertices[i][2] = float32(v[2])
	}
	return vertices
}

// indicesAsInt32 returns the indices in the block as int32.
// Thread-safe.
func (b *MeshBlock) indicesAsUint16() []uint16 {
	b.RLock()
	defer b.RUnlock()
	indices := make([]uint16, len(b.triangles)*3)
	for i, t := range b.triangles {
		indices[i*3+0] = uint16(t[0])
		indices[i*3+1] = uint16(t[1])
		indices[i*3+2] = uint16(t[2])
	}
	return indices
}

// colorsAsUint32 returns the colors in the block as uint32.
// Thread-safe.
func (b *MeshBlock) colorsAsUint8() [][3]uint8 {
	b.RLock()
	defer b.RUnlock()
	colors := make([][3]uint8, len(b.colors))
	for i, c := range b.colors {
		colors[i][0] = c[0]
		colors[i][1] = c[1]
		colors[i][2] = c[2]
	}
	return colors
}

// Gltf returns the vertices and triangles in the block as glTF bytes.
// Thread-safe.
func (b *MeshBlock) Gltf() (bytes.Buffer, error) {
	b.RLock()
	defer b.RUnlock()

	doc := gltf.NewDocument()
	positionAccessor := modeler.WritePosition(doc, b.verticesAsFloat32())
	indicesAccessor := modeler.WriteIndices(doc, b.indicesAsUint16())
	colorIndices := modeler.WriteColor(doc, b.colorsAsUint8())
	doc.Meshes = []*gltf.Mesh{{
		Primitives: []*gltf.Primitive{
			{
				Indices: gltf.Index(indicesAccessor),
				Attributes: map[string]uint32{
					gltf.POSITION: positionAccessor,
					gltf.COLOR_0:  colorIndices,
				},
				Mode: gltf.PrimitiveTriangles,
			},
		},
	}}
	doc.Nodes = []*gltf.Node{{Name: fmt.Sprintf(b.String()), Mesh: gltf.Index(0)}}
	doc.Scenes[0].Nodes = append(doc.Scenes[0].Nodes, 0)

	var buf bytes.Buffer
	err := gltf.NewEncoder(&buf).Encode(doc)
	if err != nil {
		return buf, err
	}

	// TODO: GZIP the buffer

	return buf, nil
}
