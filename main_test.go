package main

import (
	"go-voxblox/voxblox"
	"testing"
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
