package main

import (
	"go-voxblox/voxblox"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
		assert.Equal(t, test.output, float32ToRGB(test.input))
	}
}

func TestPointCloud2ToPointCloud(t *testing.T) {
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

	pointCloud2.Fields = []sensor_msgs.PointField{
		{
			Name:     "x",
			Offset:   0,
			Datatype: 7,
			Count:    1,
		},
		{
			Name:     "y",
			Offset:   4,
			Datatype: 7,
			Count:    1,
		},
		{
			Name:     "z",
			Offset:   8,
			Datatype: 7,
			Count:    1,
		},
		{
			Name:     "rgb",
			Offset:   16,
			Datatype: 7,
			Count:    1,
		},
	}

	// Read the PointCloud2 data from the test file
	data, err := os.ReadFile("testdata/PointCloud2.bin")
	assert.NoError(t, err)
	assert.Len(t, data, int(pointCloud2.RowStep*pointCloud2.Height))

	pointCloud2.Data = data
	pointCloud := PointCloud2ToPointCloud(&pointCloud2)
	assert.Len(t, pointCloud.Points, 148284)
	assert.Equal(
		t,
		voxblox.Point{-0.843722939491272, -0.6124486327171326, 1.499000072479248},
		pointCloud.Points[0],
	)
	assert.Equal(t, voxblox.Color{65, 69, 69}, pointCloud.Colors[0])
}
