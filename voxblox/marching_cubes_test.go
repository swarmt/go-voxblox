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
		{-15.6499996, -0.549999952, 1.95000005},
		{-15.5499992, -0.549999952, 1.95000005},
		{-15.5499992, -0.449999958, 1.95000005},
		{-15.6499996, -0.449999958, 1.95000005},
		{-15.6499996, -0.549999952, 2.04999995},
		{-15.5499992, -0.549999952, 2.04999995},
		{-15.5499992, -0.449999958, 2.04999995},
		{-15.6499996, -0.449999958, 2.04999995},
	}
	vertexSdf := [8]float64{
		0.393717766, 0.361241013, -0.0142535316, 0.305908412,
		-0.0258927029, -0.0836802125, -0.177394032, -0.167008847,
	}
	vertexIndex := 489
	mesh := Mesh{}

	meshCube(&vertexCoords, &vertexSdf, &vertexIndex, &mesh)
}
