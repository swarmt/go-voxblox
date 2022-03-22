package voxblox

import (
	"testing"
)

func TestCalculateVertexConfiguration(t *testing.T) {
	vertexSdf := [8]float64{
		0.5, 0.374400645, -0.036557395, 0.188546658,
		0.490951151, 0.121949553, -0.120452709, 0.0246924404,
	}
	result := calculateVertexConfiguration(&vertexSdf)
	if result != 68 {
		t.Errorf("Expected 68, got %d", result)
	}

	vertexSdf = [8]float64{
		-0.939952955, 0.00457197428, 0.111684874, 0.0445192195,
		-0.0905351862, -0.0991786737, 0.110729337, 0.038955044,
	}
	result = calculateVertexConfiguration(&vertexSdf)
	if result != 49 {
		t.Errorf("Expected 0, got %d", result)
	}

	vertexSdf = [8]float64{
		0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5,
	}
	result = calculateVertexConfiguration(&vertexSdf)
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestInterpolateEdge(t *testing.T) {
	vertex1 := Point{-19.6500015, -3.04999995, 2.44999981}
	vertex2 := Point{-19.6500015, -2.95000005, 2.44999981}
	sdf1 := 0.5
	sdf2 := -0.275501609
	point := interpolateEdge(vertex1, vertex2, sdf1, sdf2)
	if !almostEqual(point[0], -19.6500015, kEpsilon) ||
		!almostEqual(point[1], -2.98552561, kEpsilon) ||
		!almostEqual(point[2], 2.44999981, kEpsilon) {
		t.Errorf("Expected (-19.6500015, -2.98552561, 2.44999981), got (%f, %f, %f)",
			point[0], point[1], point[2])
	}
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

	if !almostEqual(edgeCoords[0][0], 0.0, kEpsilon) {
		t.Errorf("Expected 0.0, got %f", edgeCoords[0][0])
	}
	if !almostEqual(edgeCoords[11][2], 1.70743394, kEpsilon) {
		t.Errorf("Expected 1.70743394, got %f", edgeCoords[11][2])
	}
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

	if meshBlock.getVertexCount() != 3 {
		t.Errorf("Expected 3, got %d", meshBlock.getVertexCount())
	}

	if !almostEqual(meshBlock.vertices[0][0], -1.85000002, kEpsilon) ||
		!almostEqual(meshBlock.vertices[0][1], -4.45000029, kEpsilon) ||
		!almostEqual(meshBlock.vertices[0][2], 1.5111171, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}

	if !almostEqual(meshBlock.vertices[1][0], -1.79591823, kEpsilon) ||
		!almostEqual(meshBlock.vertices[1][1], -4.45000029, kEpsilon) ||
		!almostEqual(meshBlock.vertices[1][2], 1.55000007, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}

	if !almostEqual(meshBlock.vertices[2][0], -1.85000002, kEpsilon) ||
		!almostEqual(meshBlock.vertices[2][1], -4.49987888, kEpsilon) ||
		!almostEqual(meshBlock.vertices[2][2], 1.55000007, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}

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

	if meshBlock.getVertexCount() != 6 {
		t.Errorf("Expected 6, got %d", meshBlock.getVertexCount())
	}

	if !almostEqual(meshBlock.vertices[3][0], -15.5500002, kEpsilon) ||
		!almostEqual(meshBlock.vertices[3][1], -0.550000012, kEpsilon) ||
		!almostEqual(meshBlock.vertices[3][2], 1.92335689, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}

	if !almostEqual(meshBlock.vertices[4][0], -15.5233574, kEpsilon) ||
		!almostEqual(meshBlock.vertices[4][1], -0.550000012, kEpsilon) ||
		!almostEqual(meshBlock.vertices[4][2], 1.95000005, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}

	if !almostEqual(meshBlock.vertices[5][0], -15.5500002, kEpsilon) ||
		!almostEqual(meshBlock.vertices[5][1], -0.585834801, kEpsilon) ||
		!almostEqual(meshBlock.vertices[5][2], 1.95000005, kEpsilon) {
		t.Errorf("Incorrect coordinates")
	}
}
