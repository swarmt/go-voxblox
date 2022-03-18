package voxblox

import (
	"time"
)

type SimpleTsdfIntegrator struct {
	Config TsdfConfig
	Layer  *TsdfLayer
}

// NewSimpleTsdfIntegrator creates a new SimpleTsdfIntegrator.
func NewSimpleTsdfIntegrator(
	config TsdfConfig,
	layer *TsdfLayer,
) *SimpleTsdfIntegrator {
	return &SimpleTsdfIntegrator{
		Config: config,
		Layer:  layer,
	}
}

// integratePointCloud integrates a point cloud into the TSDF Layer.
func (i *SimpleTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer timeTrack(time.Now(), "Integrate Simple")
	integratePointsParallel(
		i.Layer,
		i.Config,
		pose,
		pointCloud,
	)
}

type MergedTsdfIntegrator struct {
	Config TsdfConfig
	Layer  *TsdfLayer
}

// NewMergedTsdfIntegrator creates a new MergedTsdfIntegrator.
func NewMergedTsdfIntegrator(
	config TsdfConfig,
	layer *TsdfLayer,
) *MergedTsdfIntegrator {
	return &MergedTsdfIntegrator{
		Config: config,
		Layer:  layer,
	}
}

// bundleRays decimates the point cloud by the voxel size.
func bundleRays(voxelSizeInv float64, pointCloud PointCloud) map[IndexType]int {
	voxelMap := make(map[IndexType]int)
	for j, point := range pointCloud.Points {
		voxelIndex := getGridIndexFromPoint(point, voxelSizeInv)
		voxelMap[voxelIndex] = j
	}
	return voxelMap
}

// integratePointCloud integrates a point cloud into the TSDF Layer.
func (i *MergedTsdfIntegrator) integratePointCloud(
	pose Transformation,
	pointCloud PointCloud,
) {
	defer timeTrack(time.Now(), "Integrate Merged")

	voxelMap := bundleRays(i.Layer.VoxelSizeInv, pointCloud)

	// Filter the point cloud by the voxel map
	filteredPoints := make([]Point, 0, len(pointCloud.Points))
	filteredColors := make([]Color, 0, len(pointCloud.Colors))
	for _, pointIndex := range voxelMap {
		filteredPoints = append(
			filteredPoints,
			pointCloud.Points[pointIndex],
		)
		filteredColors = append(
			filteredColors,
			pointCloud.Colors[pointIndex],
		)
	}
	pointCloud.Points = filteredPoints
	pointCloud.Colors = filteredColors

	integratePointsParallel(
		i.Layer,
		i.Config,
		pose,
		pointCloud,
	)
}
