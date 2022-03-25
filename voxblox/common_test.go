package voxblox

import "testing"

func TestGetGridIndexFromPoint(t *testing.T) {
	point := Point{1.31130219e-06, 5.2854619, 1.1920929e-07}
	globalVoxelIndex := getGridIndexFromPoint(point, 2.0*10.0)
	testIndex := IndexType{0, 105, 0}
	if globalVoxelIndex != testIndex {
		t.Errorf("Expected grid index to be {0, 105, 0}, got %v", globalVoxelIndex)
	}

	point = Point{-0.0166654587, 5.2854619, 1.1920929e-07}
	globalVoxelIndex = getGridIndexFromPoint(point, 2.0*10.0)
	testIndex = IndexType{-1, 105, 0}
	if globalVoxelIndex != testIndex {
		t.Errorf("Expected grid index to be {-1, 105, 0}, got %v", globalVoxelIndex)
	}
}

func TestGetCenterPointFromGridIndex(t *testing.T) {
	globalVoxelIndex := IndexType{-2, 51, -3}
	centerPoint := getCenterPointFromGridIndex(globalVoxelIndex, 0.1)
	testPoint := Point{-0.15, 5.15, -0.25}
	if !almostEqual(centerPoint[0], testPoint[0], kEpsilon) {
		t.Errorf("Expected center point to be -0.15, got %v", centerPoint)
	}
	if !almostEqual(centerPoint[1], testPoint[1], kEpsilon) {
		t.Errorf("Expected center point to be 5.15, got %v", centerPoint)
	}
	if !almostEqual(centerPoint[2], testPoint[2], kEpsilon) {
		t.Errorf("Expected center point to be -0.25, got %v", centerPoint)
	}

	globalVoxelIndex = IndexType{-2, 56, 9}
	centerPoint = getCenterPointFromGridIndex(globalVoxelIndex, 0.1)
	testPoint = Point{-0.15, 5.65, 0.95}
	if !almostEqual(centerPoint[0], testPoint[0], kEpsilon) {
		t.Errorf("Expected center point to be -0.15, got %v", centerPoint)
	}
	if !almostEqual(centerPoint[1], testPoint[1], kEpsilon) {
		t.Errorf("Expected center point to be 5.65, got %v", centerPoint)
	}
	if !almostEqual(centerPoint[2], testPoint[2], kEpsilon) {
		t.Errorf("Expected center point to be 0.95, got %v", centerPoint)
	}
}

func TestBlendTwoColors(t *testing.T) {
	color1 := Color{0, 0, 0}
	color2 := Color{255, 255, 255}
	blendedColor := blendTwoColors(color1, 0, color2, 1.0)
	testColor := Color{255, 255, 255}
	if blendedColor != testColor {
		t.Errorf("Expected blended color to be {255, 255, 255}, got %v", blendedColor)
	}

	color1 = Color{255, 255, 255}
	color2 = Color{255, 255, 255}
	blendedColor = blendTwoColors(color1, 0.500417829, color2, 0.499582082)
	testColor = Color{255, 255, 255}
	if blendedColor != testColor {
		t.Errorf("Expected blended color to be {255, 255, 255}, got %v", blendedColor)
	}

}
