package voxblox

import (
	"fmt"
	"os"
)

// convertTsdfLayerToTxtFile writes a tsdf Layer to a text file.
// TODO: This is a temporary function.
// TODO: The output looks incorrect, need to fix.
func convertTsdfLayerToTxtFile(layer *TsdfLayer, fileName string) {
	// Create a new file
	file, _ := os.Create(fileName)

	voxelCenters := layer.getVoxelCenters()
	for _, voxel := range voxelCenters {
		// Write the voxel to the file
		fmt.Fprintf(file, "%f %f %f\n", voxel[0], voxel[1], voxel[2])
	}
}
