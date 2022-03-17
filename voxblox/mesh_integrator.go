package voxblox

import (
	"github.com/ungerik/go3d/float64/vec3"
)

type MeshIntegrator struct {
	Config           MeshConfig
	CubeIndexOffsets []IndexType
	CubeCoordOffsets []Point
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
	i.CubeIndexOffsets = []IndexType{
		{0, 1, 1},
		{0, 0, 1},
		{1, 0, 0},
		{0, 1, 1},
		{0, 0, 1},
		{1, 0, 0},
		{0, 0, 1},
		{1, 1, 1},
	}
	i.CubeCoordOffsets = make([]Point, 8)
	for j := 0; j < 8; j++ {
		offset := IndexToPoint(i.CubeIndexOffsets[j])
		i.CubeCoordOffsets[j] = offset.Scaled(i.TsdfLayer.VoxelSize)
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

	cornerCoords := [8][3]float64{}
	cornerSdf := [8]float64{}

	allNeighborsObserved := true

	for j := 0; j < 8; j++ {
		indexOffset := i.CubeIndexOffsets[j]
		cornerIndex := AddIndex(voxelIndex, indexOffset)
		voxel := tsdfBlock.getVoxelIfExists(cornerIndex)
		if voxel == nil {
			allNeighborsObserved = false
			continue
		} else if voxel.getWeight() < i.Config.MinWeight {
			allNeighborsObserved = false
			continue
		}
		cornerCoords[j] = vec3.Add(&i.CubeCoordOffsets[j], &coords)
		cornerSdf[j] = voxel.getWeight()
	}
	if allNeighborsObserved {
		meshCube(
			&cornerCoords,
			&cornerSdf,
			vertexIndex,
			&meshBlock.mesh,
		)
	}
}

func (i *MeshIntegrator) extractMeshOnBorder(
	tsdfBlock *TsdfBlock,
	meshBlock *MeshBlock,
	voxelIndex IndexType,
	vertexIndex *int,
) {
	// TODO
}

func (i *MeshIntegrator) updateMeshForBlock(tsdfBlock *TsdfBlock) {
	meshBlock := i.MeshLayer.getBlockByIndex(tsdfBlock.Index)

	vps := i.TsdfLayer.VoxelsPerSide
	nextMeshIndex := 0

	voxelIndex := IndexType{}

	// Inside block
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

	// Max X plane
	// takes care of edge (x_max, y_max, z),
	// takes care of edge (x_max, y, z_max).
	voxelIndex[0] = vps - 1
	for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
		for voxelIndex[1] = 0; voxelIndex[1] < vps; voxelIndex[1]++ {
			i.extractMeshOnBorder(
				tsdfBlock,
				meshBlock,
				voxelIndex,
				&nextMeshIndex,
			)
		}
	}
}

func (i *MeshIntegrator) generateMesh() {
	// TODO: parallelize
	for _, block := range i.TsdfLayer.getUpdatedBlocks() {
		i.updateMeshForBlock(block)
	}
}
