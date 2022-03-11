package voxblox

type Mesh struct {
	Vertices []Point
	Indices  []uint32
}

type MeshIntegrator struct {
	Config           MeshConfig
	VoxelSize        float64
	VoxelSizeInv     float64
	BlockSize        float64
	BlockSizeInv     float64
	VoxelsPerSide    int
	VoxelsPerSideInv float64
	TsdfLayer        *TsdfLayer
	MeshLayer        *MeshLayer
}

func NewMeshIntegrator(
	config MeshConfig,
	tsdfLayer *TsdfLayer,
	meshLayer *MeshLayer,
) *MeshIntegrator {
	i := MeshIntegrator{}
	i.Config = config
	i.VoxelSize = tsdfLayer.VoxelSize
	i.BlockSize = tsdfLayer.BlockSize
	i.VoxelsPerSide = tsdfLayer.VoxelsPerSide
	i.VoxelSizeInv = 1.0 / i.VoxelSize
	i.BlockSizeInv = 1.0 / i.BlockSize
	i.VoxelsPerSideInv = 1.0 / float64(i.VoxelsPerSide)
	i.TsdfLayer = tsdfLayer
	i.MeshLayer = meshLayer
	return &i
}

func (i *MeshIntegrator) generateMesh() {
	updatedBlocks := i.TsdfLayer.getUpdatedBlocks()
	_ = updatedBlocks
}
