package voxblox

import (
	"time"
)

type TsdfIntegrator interface {
	IntegratePointCloud(pose Transformation, cloud PointCloud)
}

type SimpleTsdfIntegrator struct {
	Config TsdfConfig
	Layer  *TsdfLayer
}

// IntegratePointCloud integrates a point cloud into the TSDF Layer.
func (i *SimpleTsdfIntegrator) IntegratePointCloud(
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

// bundleRays decimates the point cloud by the voxel size.
func bundleRays(voxelSizeInv float64, pointCloud PointCloud) map[IndexType]int {
	voxelMap := make(map[IndexType]int)
	for j, point := range pointCloud.Points {
		voxelIndex := getGridIndexFromPoint(point, voxelSizeInv)
		voxelMap[voxelIndex] = j
	}
	return voxelMap
}

// IntegratePointCloud integrates a point cloud into the TSDF Layer.
func (i *MergedTsdfIntegrator) IntegratePointCloud(
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
