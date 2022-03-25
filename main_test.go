package main

import (
	"go-voxblox/voxblox"
	"os"
	"testing"

	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

func TestFloat32ToRGB(t *testing.T) {
	tests := []struct {
		input  float32
		output voxblox.Color
	}{
		{-2.9685543604723502e+38, voxblox.Color{95, 84, 71}},
		{-3.1686210629232505e+38, voxblox.Color{110, 97, 108}},
		{-3.1949936716585907e+38, voxblox.Color{112, 93, 87}},
	}

	for _, test := range tests {
		output := float32ToRGB(test.input)
		if output != test.output {
			t.Errorf("Expected %v, got %v", test.output, output)
		}
	}
}

func TestPointCloud2ToPointCloudFast(t *testing.T) {
	pointCloud2 := sensor_msgs.PointCloud2{
		Header:      std_msgs.Header{},
		Height:      480,
		Width:       640,
		Fields:      nil,
		IsBigendian: false,
		PointStep:   32,
		RowStep:     20480,
		Data:        nil,
		IsDense:     false,
	}

	// Read the PointCloud2 data from the test file
	data, err := os.ReadFile("test_data/PointCloud2.bin")
	if err != nil {
		t.Errorf("Error reading test file: %v", err)
	}
	if len(data) != int(pointCloud2.RowStep*pointCloud2.Height) {
		t.Errorf("Incorrect data length: %v", len(data))
	}
	pointCloud2.Data = data

	// Convert the PointCloud2 to a PointCloud
	pointCloud := pointCloud2ToPointCloud(&pointCloud2)

	if len(pointCloud.Points) != 148284 {
		t.Errorf("Incorrect number of points: %v", len(pointCloud.Points))
	}
	testPoint := voxblox.Point{-0.843722939491272, -0.6124486327171326, 1.499000072479248}
	if pointCloud.Points[0] != testPoint {
		t.Errorf("Incorrect point: %v", pointCloud.Points[0])
	}
	testColor := voxblox.Color{65, 69, 69}
	if pointCloud.Colors[0] != testColor {
		t.Errorf("Incorrect color: %v", pointCloud.Colors[0])
	}
}
