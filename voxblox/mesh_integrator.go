package voxblox

type Mesh struct {
	Vertices []Point
	Indices  []uint32
}

type MeshIntegrator struct {
	Config           MeshConfig
	CubeIndexOffsets []int
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
	i.CubeIndexOffsets = []int{
		0,
		1,
		1,
		0,
		0,
		1,
		1,
		0,
		0,
		0,
		1,
		1,
		0,
		0,
		1,
		1,
		0,
		0,
		0,
		0,
		1,
		1,
		1,
		1,
	}
	return &i
}

func extractMeshInsideBlock(
	block MeshBlock,
	voxelIndex IndexType,
	point Point,
	nextMeshIndex *IndexType,
	mesh *Mesh,
) {
}

func (i *MeshIntegrator) generateMeshBlock(block *TsdfBlock) {
	meshBlock := i.MeshLayer.getBlockByIndex(block.Index)
	_ = meshBlock

	vps := i.TsdfLayer.VoxelsPerSide
	// nextMeshIndex := 0

	voxelIndex := IndexType{}
	for voxelIndex[0] = 0; voxelIndex[0] < vps; voxelIndex[0]++ {
		for voxelIndex[1] = 0; voxelIndex[1] < vps; voxelIndex[1]++ {
			for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
				coords := block.computeCoordinatesFromVoxelIndex(voxelIndex)
				// TODO
				_ = coords
			}
		}
	}
}

func (i *MeshIntegrator) generateMesh() {
	updatedBlocks := i.TsdfLayer.getUpdatedBlocks()

	// TODO: parallelize
	for _, block := range updatedBlocks {
		i.generateMeshBlock(block)
	}
}
