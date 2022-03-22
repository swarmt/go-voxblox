package voxblox

import (
	"math"
	"runtime"
	"testing"

	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
)

var (
	world            *SimulationWorld
	poses            []Transformation
	cameraResolution vec2.T
	fovHorizontal    float64
	maxDistance      float64
	tsdfConfig       TsdfConfig
	meshConfig       MeshConfig
)

func init() {
	// Configuration
	cameraResolution = [2]float64{320, 240}
	fovHorizontal = 150.0
	maxDistance = 10.0

	tsdfConfig = TsdfConfig{
		VoxelSize:          0.1,
		VoxelsPerSide:      16,
		MinRange:           0.1,
		MaxRange:           5.0,
		TruncationDistance: 0.1 * 4.0,
		AllowClearing:      true,
		AllowCarving:       true,
		ConstWeight:        false,
		MaxWeight:          10000.0,
		Threads:            runtime.NumCPU(),
	}

	meshConfig = MeshConfig{
		UseColor:  true,
		MinWeight: 2.0,
	}

	// Create a test environment.
	// It consists of a 10x10x7m environment with a cylinder in the middle.
	minBound := Point{-5.0, -5.0, -1.0}
	maxBound := Point{5.0, 5.0, 6.0}
	world = NewSimulationWorld(tsdfConfig.VoxelSize, minBound, maxBound)
	cylinder := Cylinder{
		Center: Point{0.0, 0.0, 2.0},
		Radius: 2.0,
		Height: 4.0,
		Color:  ColorRed,
	}
	world.Objects = append(world.Objects, &cylinder)
	plane := Plane{
		Center: Point{0.0, 0.0, 0.0},
		Normal: vec3.T{0.0, 0.0, 1.0},
		Color:  ColorWhite,
	}
	world.AddObject(&plane)

	// Generate poses around the cylinder.
	radius := 6.0
	height := 2.0
	numPoses := 50
	maxAngle := 2.0 * math.Pi
	angleIncrement := maxAngle / float64(numPoses)
	poses = []Transformation{}
	for angle := 0.0; angle < maxAngle; angle += angleIncrement {
		position := Point{
			radius * math.Sin(angle),
			radius * math.Cos(angle),
			height,
		}
		facingDirection := vec3.Sub(&cylinder.Center, &position)
		desiredYaw := -math.Pi / 2.0
		if facingDirection[0] > 1e-4 || facingDirection[1] > 1e-4 {
			desiredYaw = math.Atan2(facingDirection[1], facingDirection[0])
		}
		qY := quaternion.FromYAxisAngle(-0.1)
		qZ := quaternion.FromZAxisAngle(desiredYaw)
		q := quaternion.Mul(&qY, &qZ)
		transform := Transformation{
			Translation: position,
			Rotation:    q,
		}
		poses = append(poses, transform)
	}
}

func TestSimpleIntegratorSingleCloud(t *testing.T) {
	// Simple integrator
	tsdfLayer := NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	simpleTsdfIntegrator := SimpleTsdfIntegrator{tsdfConfig, tsdfLayer}

	pointCloud := world.getPointCloudFromTransform(
		&poses[0],
		cameraResolution,
		fovHorizontal,
		maxDistance,
	)

	poseInverse := poses[0].inverse()
	transformedPointCloud := transformPointCloud(poseInverse, pointCloud)

	// Check transformed point cloud.
	if !almostEqual(pointCloud.Points[0][0], -2.66666627, 0.001) ||
		!almostEqual(pointCloud.Points[0][1], 5.28546286, 0.001) ||
		!almostEqual(pointCloud.Points[0][2], 0.0, 0.001) {
		t.Errorf("Pointcloud is not correct")
	}
	if !almostEqual(transformedPointCloud.Points[0][0], 0.714538097, 0.001) ||
		!almostEqual(transformedPointCloud.Points[0][1], -2.8530097, 0.001) ||
		!almostEqual(transformedPointCloud.Points[0][2], -1.72378588, 0.001) {
		t.Errorf("Transformed pointcloud is not correct")
	}

	simpleTsdfIntegrator.IntegratePointCloud(poses[0], transformedPointCloud)

	if tsdfLayer.getBlockCount() != 62 {
		t.Errorf("Number of allocated blocks is not correct")
	}

	_, voxel := getBlockAndVoxelFromGlobalVoxelIndex(
		simpleTsdfIntegrator.Layer,
		IndexType{0, 60, 20},
	)
	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Wrong distance: %v", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 10000.0, kEpsilon) {
		t.Errorf("Wrong weight: %v", voxel.getWeight())
	}

	voxel = tsdfLayer.getBlockByIndex(IndexType{-1, 0, 2}).getVoxel(IndexType{4, 15, 0})
	if !almostEqual(voxel.getDistance(), -0.122520447, 0.001) {
		t.Errorf("Wrong distance: %v", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.531333983, 0.05) {
		t.Errorf("Wrong weight: %v", voxel.getWeight())
	}

	// Check every voxel has color
	for _, block := range tsdfLayer.getBlocks() {
		for _, voxel := range block.getVoxels() {
			if math.Abs(voxel.getDistance()) < tsdfLayer.VoxelSize {
				color := voxel.getColor()
				if color[0] == 0 && color[1] == 0 && color[2] == 0 {
					t.Errorf("Wrong color: %v %v", block.Index, voxel.Index)
				}
			}
		}
	}

	WriteTsdfLayerToTxtFile(tsdfLayer, "../output/simple_layer.txt")

	// Generate Mesh.
	meshLayer := NewMeshLayer(tsdfLayer)
	meshIntegrator := NewMeshIntegrator(meshConfig, tsdfLayer, meshLayer)
	meshIntegrator.IntegrateMesh()

	if meshLayer.getBlockCount() != tsdfLayer.getBlockCount() {
		t.Errorf("Number of allocated blocks is not correct")
	}
}

func TestTsdfIntegrators(t *testing.T) {
	// Simple integrator
	simpleLayer := NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	simpleTsdfIntegrator := SimpleTsdfIntegrator{tsdfConfig, simpleLayer}

	// Merged integrator
	mergedLayer := NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	mergedTsdfIntegrator := MergedTsdfIntegrator{tsdfConfig, mergedLayer}

	// TODO: Fast integrator

	// Iterate over all poses and integrate.
	for _, pose := range poses {
		pointCloud := world.getPointCloudFromTransform(
			&pose,
			cameraResolution,
			fovHorizontal,
			maxDistance,
		)
		poseInverse := pose.inverse()
		transformedPointCloud := transformPointCloud(poseInverse, pointCloud)
		simpleTsdfIntegrator.IntegratePointCloud(pose, transformedPointCloud)
		mergedTsdfIntegrator.IntegratePointCloud(pose, transformedPointCloud)
	}

	// Check the number of blocks in the layers
	if simpleLayer.getBlockCount() == 0 {
		t.Errorf("No blocks in simple Layer")
	}

	// Check a block Origin
	block01Neg1 := simpleLayer.getBlockByIndex(IndexType{0, 1, -1})
	origin := block01Neg1.Origin
	if origin[0] != 0.0 || origin[1] != 1.6 || origin[2] != -1.6 {
		t.Errorf("Wrong block Origin: %v", origin)
	}

	// Check every voxel has color
	for _, block := range simpleLayer.getBlocks() {
		for _, voxel := range block.getVoxels() {
			if math.Abs(voxel.getDistance()) < simpleLayer.VoxelSize && voxel.getWeight() > 0.0 {
				color := voxel.getColor()
				if color[0] == 0 && color[1] == 0 && color[2] == 0 {
					t.Errorf("Wrong color: %v %v", block.Index, voxel.Index)
				}
			}
		}
	}

	WriteTsdfLayerToTxtFile(simpleLayer, "../output/simple_layer.txt")

	// Generate simple layer mesh.
	simpleMeshLayer := NewMeshLayer(simpleLayer)
	meshIntegrator := NewMeshIntegrator(meshConfig, simpleLayer, simpleMeshLayer)
	meshIntegrator.IntegrateMesh()

	if simpleMeshLayer.getBlockCount() != simpleLayer.getBlockCount() {
		t.Errorf("Number of allocated blocks is not correct")
	}

	WriteMeshLayerToObjFiles(simpleMeshLayer, "../output/simple_mesh")

	// Generate merged layer mesh.
	mergedMeshLayer := NewMeshLayer(mergedLayer)
	meshIntegrator = NewMeshIntegrator(meshConfig, mergedLayer, mergedMeshLayer)
	meshIntegrator.IntegrateMesh()

	if mergedMeshLayer.getBlockCount() != mergedLayer.getBlockCount() {
		t.Errorf("Number of allocated blocks is not correct")
	}

	WriteMeshLayerToObjFiles(mergedMeshLayer, "../output/merged_mesh")
}

func TestUpdateTsdfVoxel(t *testing.T) {
	layer := NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	origin := Point{0.0, 6.0, 2.0}
	pointC := Point{0.714538097, -2.8530097, -1.72378588}
	pointG := Point{-2.66666508, 5.2854619, 1.1920929e-07}
	globalVoxelIndex := IndexType{0, 60, 20}
	_, voxel := getBlockAndVoxelFromGlobalVoxelIndex(layer, globalVoxelIndex)
	weight := calculateWeight(pointC)

	tsdfConfig = TsdfConfig{
		VoxelSize:          0.1,
		VoxelsPerSide:      10,
		TruncationDistance: 0.4,
		MaxWeight:          10000.0,
		ConstWeight:        false,
	}

	simpleTsdfIntegrator := SimpleTsdfIntegrator{tsdfConfig, layer}

	updateTsdfVoxel(
		layer,
		simpleTsdfIntegrator.Config,
		origin,
		pointG,
		globalVoxelIndex,
		Color{},
		weight,
		voxel,
	)

	if !almostEqual(voxel.getDistance(), 0.4, kEpsilon) {
		t.Errorf("Expected Tsdf to be 0.4, got %f", voxel.getDistance())
	}
	if !almostEqual(voxel.getWeight(), 0.336537421, kEpsilon) {
		t.Errorf("Expected weight to be 0.336537421, got %f", voxel.getWeight())
	}
	if len(layer.blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(layer.blocks))
	}
}
