package voxblox

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type Mesh struct {
	Vertices []Point
	Indices  []uint32
}

type MeshIntegrator struct {
	Config                  MeshConfig
	CubeIndexOffsetsInt     []int
	CubeIndexOffsetsFloat64 []float64
	TsdfLayer               *TsdfLayer
	MeshLayer               *MeshLayer
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
	i.CubeIndexOffsetsInt = []int{
		0, 1, 1, 0, 0, 1,
		1, 0, 0, 0, 1, 1,
		0, 0, 1, 1, 0, 0,
		0, 0, 1, 1, 1, 1,
	}
	i.CubeIndexOffsetsFloat64 = []float64{
		0.0, 1.0, 1.0, 0.0, 0.0, 1.0,
		1.0, 0.0, 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 1.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 1.0, 1.0, 1.0,
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

func (i *MeshIntegrator) extractMeshInsideBlock(block *MeshBlock,
	voxelIndex IndexType, coords Point,
	vertexIndex *IndexType, mesh *Mesh) {
	cubeCoordOffsets := mat.NewDense(3, 8, i.CubeIndexOffsetsFloat64)
	cubeCoordOffsets.Scale(
		i.TsdfLayer.VoxelSize,
		cubeCoordOffsets,
	)
	fmt.Println(cubeCoordOffsets)
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
				_ = coords
				//extractMeshInsideBlock(
				//	meshBlock,
				//	voxelIndex,
				//	coords,
				//	&nextMeshIndex,
				//	&mesh,
				//)
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
