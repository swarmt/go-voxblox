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
	i.TsdfLayer = tsdfLayer
	i.MeshLayer = meshLayer
	return &i
}

func (i *MeshIntegrator) generateMeshBlock(block *TsdfBlock) {
}

func (i *MeshIntegrator) generateMesh() {
	updatedBlocks := i.TsdfLayer.getUpdatedBlocks()

	// TODO: parallelize
	for _, block := range updatedBlocks {
		i.generateMeshBlock(block)
	}
}
