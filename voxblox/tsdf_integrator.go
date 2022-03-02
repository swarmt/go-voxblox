package voxblox

import (
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

type TsdfIntegrator interface {
	integratePointCloud(
		layer *TsdfLayer,
		transformation *Transformation,
		pointCloud *sensor_msgs.PointCloud2,
		colors []*Color,
		freeSpacePoints bool,
	)
}

type SimpleTsdfIntegrator struct {
	truncationDistance float64
	layer              *TsdfLayer
}

func NewSimpleTsdfIntegrator(truncationDistance float64, layer *TsdfLayer) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		truncationDistance: truncationDistance,
		layer:              layer,
	}
}
