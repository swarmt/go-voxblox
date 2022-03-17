package voxblox

import (
	"math"
)

var kEdgeIndexPairs = [][2]int{
	{0, 1},
	{1, 2},
	{2, 3},
	{3, 0},
	{4, 5},
	{5, 6},
	{6, 7},
	{7, 4},
	{0, 4},
	{1, 5},
	{2, 6},
	{3, 7},
}

func calculateVertexConfiguration(vertexSdf *[8]float64) int {
	var index int
	for i := 0; i < 8; i++ {
		if vertexSdf[i] < 0 {
			index |= 1 << i
		}
	}
	return index
}

// interpolateEdge performs linear interpolation on two cube corners to find the approximate
// zero crossing (surface) value.
func interpolateEdge(vertex1 Point, vertex2 Point, sdf1 float64, sdf2 float64) *Point {
	kMinSdfDifference := 1e-6
	sdfDifference := sdf1 - sdf2

	if math.Abs(sdfDifference) >= kMinSdfDifference {
		t := sdf1 / sdfDifference
		point := vertex2.Sub(&vertex1).Scale(t).Add(&vertex1)
		return point
	}
	return vertex1.Add(&vertex2).Scale(0.5)
}

func interpolateEdgeVertices(
	vertexCoords *[8][3]float64,
	vertexSdf *[8]float64,
	edgeCoords *[12][3]float64,
) {
	for i := 0; i < 12; i++ {
		pairs := kEdgeIndexPairs[i]
		edge0 := pairs[0]
		edge1 := pairs[1]
		// Only interpolate along edges where there is a zero crossing.
		if vertexSdf[edge0] < 0 && vertexSdf[edge1] >= 0 ||
			vertexSdf[edge0] >= 0 && vertexSdf[edge1] < 0 {
			edge := interpolateEdge(
				vertexCoords[edge0], vertexCoords[edge1],
				vertexSdf[edge0], vertexSdf[edge1],
			)
			edgeCoords[i] = *edge
		}
	}
}

func meshCube(
	vertexCoords *[8][3]float64,
	vertexSdf *[8]float64,
	vertexIndex *int,
	mesh *Mesh,
) {
	index := calculateVertexConfiguration(vertexSdf)
	if index == 0 {
		return
	}

	// edgeVertexCoordinates := mat.NewDense(3, 12, nil)
}
