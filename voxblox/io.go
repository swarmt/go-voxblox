package voxblox

import (
	"fmt"
	"os"
)

// WriteTsdfLayerToTxtFile writes a tsdf Layer to a text file.
// TODO: This is a temporary function.
func WriteTsdfLayerToTxtFile(layer *TsdfLayer, fileName string) {
	// Create a new file
	file, _ := os.Create(fileName)

	voxelCenters, voxelColors := layer.getVoxelCenters()
	for i, voxel := range voxelCenters {
		// Write the voxel to the file
		fmt.Fprintf(
			file,
			"%f %f %f %d %d %d\n",
			voxel[0],
			voxel[1],
			voxel[2],
			voxelColors[i][0],
			voxelColors[i][1],
			voxelColors[i][2],
		)
	}
}

// WriteMeshLayerToObjFiles writes a Mesh Layer to an obj file.
// TODO: This is a temporary function.
func WriteMeshLayerToObjFiles(layer *MeshLayer, folderName string) {
	// Create folder if it doesn't exist
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0o755)
	}

	for _, block := range layer.getBlocks() {
		if block.getVertexCount() == 0 {
			continue
		}

		// Create a new file with the block index as the name
		fileName := fmt.Sprintf("%s/%d.obj", folderName, block.Index)
		file, _ := os.Create(fileName)

		block.RLock()
		for i, vertex := range block.vertices {
			// Write the voxel to the file
			fmt.Fprintf(
				file,
				"v %f %f %f %d %d %d\n",
				vertex[0],
				vertex[1],
				vertex[2],
				block.colors[i][0],
				block.colors[i][1],
				block.colors[i][2],
			)
		}
		for _, triangle := range block.triangles {
			// Write the voxel to the file
			fmt.Fprintf(
				file,
				"f %d %d %d\n",
				triangle[0]+1,
				triangle[1]+1,
				triangle[2]+1,
			)
		}
		block.RUnlock()
	}
}
