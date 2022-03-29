package voxblox

import (
	"fmt"
	"os"
)

// WriteMeshLayerToObjFiles writes a Mesh Layer to an obj file.
func WriteMeshLayerToObjFiles(layer *MeshLayer, folderName string) {
	// Create folder if it doesn't exist
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		os.Mkdir(folderName, 0o755)
	}

	for _, block := range layer.getBlocks() {
		if block.getVertexCount() == 0 {
			continue
		}

		fileName := fmt.Sprintf("%s/%s.obj", folderName, block)
		file, _ := os.Create(fileName)

		block.RLock()
		for i, vertex := range block.vertices {
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
