package voxblox

type TsdfVoxel struct {
	Index    IndexType
	Distance float32
	Weight   float32
}

func NewVoxel(index IndexType) *TsdfVoxel {
	return &TsdfVoxel{
		Index: index,
	}
}
