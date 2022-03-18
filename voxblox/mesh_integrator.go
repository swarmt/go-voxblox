package voxblox

import (
	"time"

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
		{0, 0, 0},
		{1, 0, 0},
		{1, 1, 0},
		{0, 1, 0},
		{0, 0, 1},
		{1, 0, 1},
		{1, 1, 1},
		{0, 1, 1},
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
) {
	coords := tsdfBlock.computeCoordinatesFromVoxelIndex(voxelIndex)

	cornerCoords := [8][3]float64{}
	cornerSdf := [8]float64{}

	allNeighborsObserved := true

	for j := 0; j < 8; j++ {
		cornerIndex := AddIndex(voxelIndex, i.CubeIndexOffsets[j])
		voxel := tsdfBlock.getVoxelIfExists(cornerIndex)
		if voxel == nil {
			allNeighborsObserved = false
			break
		}
		if voxel.getWeight() < i.Config.MinWeight {
			allNeighborsObserved = false
			break
		}
		cornerCoords[j] = vec3.Add(&i.CubeCoordOffsets[j], &coords)
		cornerSdf[j] = voxel.getDistance()
	}
	if allNeighborsObserved {
		meshCube(
			&cornerCoords,
			&cornerSdf,
			meshBlock,
		)
	}
}

func (i *MeshIntegrator) extractMeshOnBorder(
	tsdfBlock *TsdfBlock,
	meshBlock *MeshBlock,
	voxelIndex IndexType,
) {
	coords := tsdfBlock.computeCoordinatesFromVoxelIndex(voxelIndex)

	cornerCoords := [8][3]float64{}
	cornerSdf := [8]float64{}

	allNeighborsObserved := true

	for j := 0; j < 8; j++ {
		cornerIndex := AddIndex(voxelIndex, i.CubeIndexOffsets[j])

		if tsdfBlock.isValidVoxelIndex(cornerIndex) {
			voxel := tsdfBlock.getVoxelIfExists(cornerIndex)
			if voxel == nil {
				allNeighborsObserved = false
				break
			}
			if voxel.getWeight() < i.Config.MinWeight {
				allNeighborsObserved = false
				break
			}
			cornerCoords[j] = vec3.Add(&i.CubeCoordOffsets[j], &coords)
			cornerSdf[j] = voxel.getDistance()
		} else {
			// We have to access a different block.
			blockOffset := IndexType{}

			for k := 0; k < 3; k++ {
				if cornerIndex[k] < 0 {
					blockOffset[k] = -1
					cornerIndex[k] = cornerIndex[k] + tsdfBlock.VoxelsPerSide
				} else if cornerIndex[k] >= tsdfBlock.VoxelsPerSide {
					blockOffset[k] = 1
					cornerIndex[k] = cornerIndex[k] - tsdfBlock.VoxelsPerSide
				}
			}

			neighborIndex := AddIndex(tsdfBlock.Index, blockOffset)
			neighborBlock := i.TsdfLayer.getBlockIfExists(neighborIndex)

			if neighborBlock == nil {
				allNeighborsObserved = false
				break
			}

			voxel := neighborBlock.getVoxelIfExists(cornerIndex)
			if voxel == nil {
				allNeighborsObserved = false
				break
			}
			if voxel.getWeight() < i.Config.MinWeight {
				allNeighborsObserved = false
				break
			}
			cornerCoords[j] = vec3.Add(&i.CubeCoordOffsets[j], &coords)
			cornerSdf[j] = voxel.getDistance()
		}
	}
	if allNeighborsObserved {
		meshCube(
			&cornerCoords,
			&cornerSdf,
			meshBlock,
		)
	}
}

func (i *MeshIntegrator) updateMeshForBlock(tsdfBlock *TsdfBlock) {
	meshBlock := i.MeshLayer.getBlockByIndex(tsdfBlock.Index)

	vps := i.TsdfLayer.VoxelsPerSide

	meshBlock.Mesh.VertexCount = 0

	voxelIndex := IndexType{}

	// Inside block
	for voxelIndex[0] = 0; voxelIndex[0] < vps; voxelIndex[0]++ {
		for voxelIndex[1] = 0; voxelIndex[1] < vps; voxelIndex[1]++ {
			for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
				i.extractMeshInsideBlock(
					tsdfBlock,
					meshBlock,
					voxelIndex,
				)
			}
		}
	}

	// Max X plane
	voxelIndex[0] = vps - 1
	for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
		for voxelIndex[1] = 0; voxelIndex[1] < vps; voxelIndex[1]++ {
			i.extractMeshOnBorder(
				tsdfBlock,
				meshBlock,
				voxelIndex,
			)
		}
	}

	// Max Y plane
	voxelIndex[1] = vps - 1
	for voxelIndex[2] = 0; voxelIndex[2] < vps; voxelIndex[2]++ {
		for voxelIndex[0] = 0; voxelIndex[0] < vps-1; voxelIndex[0]++ {
			i.extractMeshOnBorder(
				tsdfBlock,
				meshBlock,
				voxelIndex,
			)
		}
	}

	// Max Z plane
	voxelIndex[2] = vps - 1
	for voxelIndex[1] = 0; voxelIndex[1] < vps-1; voxelIndex[1]++ {
		for voxelIndex[0] = 0; voxelIndex[0] < vps-1; voxelIndex[0]++ {
			i.extractMeshOnBorder(
				tsdfBlock,
				meshBlock,
				voxelIndex,
			)
		}
	}
}

func (i *MeshIntegrator) integrateMesh() {
	defer timeTrack(time.Now(), "integrateMesh")

	// TODO: parallelize
	for _, block := range i.TsdfLayer.getUpdatedBlocks() {
		i.updateMeshForBlock(block)
	}
}
