package voxblox

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateVertexConfiguration(t *testing.T) {
	vertexSdf := [8]float64{
		0.5, 0.374400645, -0.036557395, 0.188546658,
		0.490951151, 0.121949553, -0.120452709, 0.0246924404,
	}
	assert.Equal(t, 68, calculateVertexConfiguration(&vertexSdf))

	vertexSdf = [8]float64{
		-0.939952955, 0.00457197428, 0.111684874, 0.0445192195,
		-0.0905351862, -0.0991786737, 0.110729337, 0.038955044,
	}
	assert.Equal(t, 49, calculateVertexConfiguration(&vertexSdf))

	vertexSdf = [8]float64{
		0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5,
	}
	assert.Equal(t, 0, calculateVertexConfiguration(&vertexSdf))
}

func TestInterpolateEdge(t *testing.T) {
	point := interpolateEdge(
		Point{-19.6500015, -3.04999995, 2.44999981},
		Point{-19.6500015, -2.95000005, 2.44999981},
		0.5,
		-0.275501609,
	)
	assert.InEpsilon(t, -19.6500015, point[0], kEpsilon)
	assert.InEpsilon(t, -2.98552561, point[1], kEpsilon)
	assert.InEpsilon(t, 2.44999981, point[2], kEpsilon)
}

func TestInterpolateEdgeVertices(t *testing.T) {
	vertexCoords := [8][3]float64{
		{-0.549999952, -3.45000029, 1.64999998},
		{-0.449999958, -3.45000029, 1.64999998},
		{-0.449999958, -3.35000038, 1.64999998},
		{-0.549999952, -3.35000038, 1.64999998},
		{-0.549999952, -3.45000029, 1.75},
		{-0.449999958, -3.45000029, 1.75},
		{-0.449999958, -3.35000038, 1.75},
		{-0.549999952, -3.35000038, 1.75},
	}
	vertexSdf := [8]float64{
		-0.325170636,
		-0.40814352,
		-0.387115151,
		-0.30503884,
		0.054089319,
		-0.0170794651,
		0.126751885,
		0.226074144,
	}
	edgeCoords := [12][3]float64{}
	interpolateEdgeVertices(&vertexCoords, &vertexSdf, &edgeCoords)
	assert.Equal(t, 0.0, edgeCoords[0][0])
	assert.InEpsilon(t, 1.70743394, edgeCoords[11][2], kEpsilon)
}

func TestMeshCube(t *testing.T) {
	vertexCoords := [8][3]float64{
		{-1.85000002, -4.55000019, 1.45000005},
		{-1.75, -4.55000019, 1.45000005},
		{-1.75, -4.45000029, 1.45000005},
		{-1.85000002, -4.45000029, 1.45000005},
		{-1.85000002, -4.55000019, 1.55000007},
		{-1.75, -4.55000019, 1.55000007},
		{-1.75, -4.45000029, 1.55000007},
		{-1.85000002, -4.45000029, 1.55000007},
	}
	vertexSdf := [8]float64{
		0.360097408,
		0.311582834,
		0.21046859,
		0.279508501,
		0.17869179,
		0.320180953,
		0.150982603,
		-0.177825138,
	}
	tsdfLayer := NewTsdfLayer(0.1, 16)
	meshLayer := NewMeshLayer(tsdfLayer)
	meshBlock := NewMeshBlock(meshLayer, IndexType{0, 0, 0}, Point{0, 0, 0})

	meshCube(&vertexCoords, &vertexSdf, meshBlock)

	assert.Equal(t, 3, meshBlock.getVertexCount())

	assert.InEpsilon(t, -1.85000002, meshBlock.vertices[0][0], kEpsilon)
	assert.InEpsilon(t, -4.45000029, meshBlock.vertices[0][1], kEpsilon)
	assert.InEpsilon(t, 1.5111171, meshBlock.vertices[0][2], kEpsilon)

	assert.InEpsilon(t, -1.79591823, meshBlock.vertices[1][0], kEpsilon)
	assert.InEpsilon(t, -4.45000029, meshBlock.vertices[1][1], kEpsilon)
	assert.InEpsilon(t, 1.55000007, meshBlock.vertices[1][2], kEpsilon)

	assert.InEpsilon(t, -1.85000002, meshBlock.vertices[2][0], kEpsilon)
	assert.InEpsilon(t, -4.49987888, meshBlock.vertices[2][1], kEpsilon)
	assert.InEpsilon(t, 1.55000007, meshBlock.vertices[2][2], kEpsilon)

	vertexCoords = [8][3]float64{
		{-15.5500002, -0.650000036, 1.85000002},
		{-15.4499998, -0.650000036, 1.85000002},
		{-15.4499998, -0.550000012, 1.85000002},
		{-15.5500002, -0.550000012, 1.85000002},
		{-15.5500002, -0.650000036, 1.95000005},
		{-15.4499998, -0.650000036, 1.95000005},
		{-15.4499998, -0.550000012, 1.95000005},
		{-15.5500002, -0.550000012, 1.95000005},
	}
	vertexSdf = [8]float64{
		0.5,
		0.5,
		0.5,
		0.5,
		0.325169295,
		0.5,
		0.5,
		-0.181599542,
	}

	meshCube(&vertexCoords, &vertexSdf, meshBlock)

	assert.Equal(t, 6, meshBlock.getVertexCount())

	assert.InEpsilon(t, -15.5500002, meshBlock.vertices[3][0], kEpsilon)
	assert.InEpsilon(t, -0.550000012, meshBlock.vertices[3][1], kEpsilon)
	assert.InEpsilon(t, 1.92335689, meshBlock.vertices[3][2], kEpsilon)

	assert.InEpsilon(t, -15.5233574, meshBlock.vertices[4][0], kEpsilon)
	assert.InEpsilon(t, -0.550000012, meshBlock.vertices[4][1], kEpsilon)
	assert.InEpsilon(t, 1.95000005, meshBlock.vertices[4][2], kEpsilon)

	assert.InEpsilon(t, -15.5500002, meshBlock.vertices[5][0], kEpsilon)
	assert.InEpsilon(t, -0.585834801, meshBlock.vertices[5][1], kEpsilon)
	assert.InEpsilon(t, 1.95000005, meshBlock.vertices[5][2], kEpsilon)
}

func TestVertexIndex(t *testing.T) {
	meshBlock := new(MeshBlock)
	assert.Equal(t, 0, vertexIndex(meshBlock, Point{0, 0, 0}))
	assert.Equal(t, 0, vertexIndex(meshBlock, Point{0, 0, 0}))
	assert.Equal(t, 1, vertexIndex(meshBlock, Point{1, 1, 1}))
	assert.Equal(t, 0, vertexIndex(meshBlock, Point{0, 0, 0}))
	assert.Len(t, meshBlock.vertices, 2)
}
