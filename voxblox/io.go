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

	for _, block := range layer.GetBlocks() {
		if block.getVertexCount() == 0 {
			continue
		}

		fileName := fmt.Sprintf("%s/%s.obj", folderName, block)
		file, _ := os.Create(fileName)

		block.RLock()
		for i, vertex := range block.vertices {
			r := float64(block.colors[i][0]) / 255.0
			g := float64(block.colors[i][1]) / 255.0
			b := float64(block.colors[i][2]) / 255.0

			fmt.Fprintf(
				file,
				"v %f %f %f %f %f %f\n",
				vertex[0],
				vertex[1],
				vertex[2],
				r,
				g,
				b,
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
