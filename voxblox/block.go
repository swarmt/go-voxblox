package voxblox

import "github.com/ungerik/go3d/float64/vec3"

type Status int

const (
	kMap Status = iota
	kMesh
	kEsdf
	kCount
)

type Block struct {
	HasData       bool
	VoxelsPerSide int
	VoxelSize     float64
	Origin        Point
	BlockIndex    IndexType
	Updated       bool
	NumVoxels     int
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
	Voxels        map[IndexType]*TsdfVoxel
}

func NewBlock(voxelsPerSide int, voxelSize float64, index IndexType, origin Point) *Block {
	b := new(Block)
	b.HasData = false
	b.VoxelsPerSide = voxelsPerSide
	b.VoxelSize = voxelSize
	b.Origin = origin
	b.BlockIndex = index
	b.Updated = false
	b.NumVoxels = voxelsPerSide * voxelsPerSide * voxelsPerSide
	b.VoxelSizeInv = 1.0 / voxelSize
	b.BlockSize = float64(voxelsPerSide) * voxelSize
	b.BlockSizeInv = 1.0 / b.BlockSize
	b.Voxels = make(map[IndexType]*TsdfVoxel)
	return b
}

// getBlock allocates a new block in the map
func (b *Block) getVoxel(voxelIndex IndexType) *TsdfVoxel {
	// Test if block already exists
	if voxel, ok := b.Voxels[voxelIndex]; ok {
		return voxel
	}
	newVoxel := NewVoxel(voxelIndex)
	b.Voxels[voxelIndex] = newVoxel
	return newVoxel
}

func (b *Block) getVoxelPtrByCoordinates(point Point) *TsdfVoxel {
	return b.getVoxel(getGridIndexFromPoint(point, b.VoxelSize))
}

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

func (b *Block) computeCoordinatesFromVoxelIndex(index IndexType) Point {
	centerPoint := getCenterPointFromGridIndex(index, b.VoxelSize)
	return vec3.Add(&b.Origin, &centerPoint)
}
