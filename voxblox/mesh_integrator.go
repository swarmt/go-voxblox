package voxblox

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

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

func (i *MeshIntegrator) extractMeshInsideBlock(
	tsdfBlock *TsdfBlock,
	meshBlock *MeshBlock,
	voxelIndex IndexType,
	vertexIndex *int,
) {
	coords := tsdfBlock.computeCoordinatesFromVoxelIndex(voxelIndex)
	coordsMat := mat.NewDense(3, 1, coords.Slice())

	cubeIndexOffsets := mat.NewDense(3, 8, i.CubeIndexOffsetsFloat64) // TODO: init in constructor
	cubeCoordOffsets := mat.NewDense(3, 8, i.CubeIndexOffsetsFloat64) // TODO: init in constructor
	cubeCoordOffsets.Scale(
		i.TsdfLayer.VoxelSize,
		cubeCoordOffsets,
	)
	cornerCoords := mat.NewDense(3, 8, nil)
	cornerSdfs := mat.NewDense(8, 1, nil)

	allNeighborsObserved := true

	for j := 0; j < 8; j++ {
		indexOffset := cubeIndexOffsets.ColView(j)
		var cornerIndex mat.Dense
		cornerIndex.Add(indexOffset, IndexToMatrix(voxelIndex))
		voxel := tsdfBlock.getVoxelIfExists(MatrixToIndex(&cornerIndex))
		if voxel == nil {
			allNeighborsObserved = false
			continue
		} else if voxel.getWeight() < i.Config.MinWeight {
			allNeighborsObserved = false
			continue
		}
		var cornerCoord mat.Dense
		cornerCoord.Add(
			cubeCoordOffsets.ColView(j),
			coordsMat,
		)
		cornerCoords.SetCol(j, cornerCoord.RawMatrix().Data)
		cornerSdfs.Set(j, 0, voxel.getWeight())
	}
	if allNeighborsObserved {
		// TODO: Marching Cubes
		fmt.Println("Marching Cubes")
	}
}

func (i *MeshIntegrator) updateMeshForBlock(tsdfBlock *TsdfBlock) {
	meshBlock := i.MeshLayer.getBlockByIndex(tsdfBlock.Index)

	vps := i.TsdfLayer.VoxelsPerSide
	nextMeshIndex := 0

	voxelIndex := IndexType{}
	for voxelIndex[0] = 0; voxelIndex[0] < vps; voxelIndex[0]++ {
		for voxelIndex[1] = 0; voxelIndex[1] < vps; voxelIndex[1]++ {
			for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
				i.extractMeshInsideBlock(
					tsdfBlock,
					meshBlock,
					voxelIndex,
					&nextMeshIndex,
				)
			}
		}
	}
}

func (i *MeshIntegrator) generateMesh() {
	updatedBlocks := i.TsdfLayer.getUpdatedBlocks()

	// TODO: parallelize
	for _, block := range updatedBlocks {
		i.updateMeshForBlock(block)
	}
}
