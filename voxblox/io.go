package voxblox

import (
	"fmt"
	"os"
)

// writeTsdfLayerToTxtFile writes a tsdf Layer to a text file.
// TODO: This is a temporary function.
func writeTsdfLayerToTxtFile(layer *TsdfLayer, fileName string) {
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

// writeMeshLayerToObjFile writes a Mesh Layer to an obj file.
// TODO: This is a temporary function.
func writeMeshLayerToObjFile(layer *MeshLayer, folderName string) {
	// Create folder if it doesn't exist
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0755)
	}

	for _, block := range layer.blocks {
		if len(block.Mesh.Vertices) == 0 {
			continue
		}

		// Create a new file with the block index as the name
		fileName := fmt.Sprintf("%s/%d.obj", folderName, block.Index)
		file, _ := os.Create(fileName)

		for _, vertex := range block.Mesh.Vertices {
			// Write the voxel to the file
			fmt.Fprintf(
				file,
				"v %f %f %f\n",
				vertex[0],
				vertex[1],
				vertex[2],
			)
		}
		for _, triangle := range block.Mesh.Triangles {
			// Write the voxel to the file
			fmt.Fprintf(
				file,
				"f %d %d %d\n",
				triangle[0]+1,
				triangle[1]+1,
				triangle[2]+1,
			)
		}
	}
}
