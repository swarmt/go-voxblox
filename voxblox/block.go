package voxblox

type Status int64

const (
	kMap Status = iota
	kMesh
	kEsdf
	kCount
)

type Block struct {
	HasData       bool
	VoxelsPerSide int32
	VoxelSize     float64
	Origin        Point
	BlockIndex    IndexType
	Updated       bool
	NumVoxels     int32
	VoxelSizeInv  float64
	BlockSize     float64
	BlockSizeInv  float64
}

func NewBlock(voxelsPerSide int32, voxelSize float64, index IndexType, origin Point) *Block {
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
	return b
}

func (b *Block) getVoxelPtrByCoordinates(point Point) {

}

func (b *Block) computeLinearIndexFromCoordinates(point Point) VoxelIndex {
	return b.computeTruncatedVoxelIndexFromCoordinates(point)
}

func (b *Block) computeTruncatedVoxelIndexFromCoordinates(point Point) VoxelIndex {
	return getGridIndexFromPoint(subtractPoints(point, b.Origin), b.BlockSizeInv)
}
