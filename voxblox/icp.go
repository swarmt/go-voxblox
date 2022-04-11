package voxblox

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ungerik/go3d/float64/quaternion"

	"github.com/ungerik/go3d/float64/vec3"
	"gonum.org/v1/gonum/mat"
)

func jacobianGaussNewton(src, tgt Point, initial *[6]float64) ([]float64, float64) {
	dx := tgt[0]
	dy := tgt[1]
	dz := tgt[2]

	sx := src[0]
	sy := src[1]
	sz := src[2]

	alpha := initial[0]
	beta := initial[1]
	gamma := initial[2]
	tx := initial[3]
	ty := initial[4]
	tz := initial[5]

	a1 := (-2 * beta * sx * sy) - (2 * gamma * sx * sz) + (2 * alpha * ((sy * sy) + (sz * sz))) + (2 * ((sz * dy) - (sy * dz))) + 2*((sy*tz)-(sz*ty))
	a2 := (-2 * alpha * sx * sy) - (2 * gamma * sy * sz) + (2 * beta * ((sx * sx) + (sz * sz))) + (2 * ((sx * dz) - (sz * dx))) + 2*((sz*tx)-(sx*tz))
	a3 := (-2 * alpha * sx * sz) - (2 * beta * sy * sz) + (2 * gamma * ((sx * sx) + (sy * sy))) + (2 * ((sy * dx) - (sx * dy))) + 2*((sx*ty)-(sy*tx))
	a4 := 2 * (sx - (gamma * sy) + (beta * sz) + tx - dx)
	a5 := 2 * (sy - (alpha * sz) + (gamma * sx) + ty - dy)
	a6 := 2 * (sz - (beta * sx) + (alpha * sy) + tz - dz)

	r := (a4 * a4 / 4) + (a5 * a5 / 4) + (a6 * a6 / 4)

	return []float64{a1, a2, a3, a4, a5, a6}, r
}

func stepICP(src, tgt []Point, transform *[6]float64) error {
	var jacobians []float64
	var residuals []float64

	for i := 0; i < len(src); i++ {
		j, r := jacobianGaussNewton(src[i], tgt[i], transform)
		jacobians = append(jacobians, j...)
		residuals = append(residuals, r)
	}

	jacobianMatrix := mat.NewDense(len(src), 6, jacobians)
	residualVector := mat.NewVecDense(len(src), residuals)

	var update mat.VecDense
	err := update.SolveVec(jacobianMatrix, residualVector)
	if err != nil {
		return err
	}

	for i := 0; i < update.Len(); i++ {
		transform[i] -= update.At(i, 0)
	}

	return nil
}

// getGradient returns the gradient of the voxel.
// Calculated using the neighbor voxels.
func getGradient(tsdfLayer *TsdfLayer, globalVoxelIndex IndexType) (Point, bool) {
	gradient := Point{}
	for i := 0; i < 3; i++ {
		for sign := -1; sign <= 1; sign += 2 {
			neighborIndex := globalVoxelIndex
			neighborIndex[i] = globalVoxelIndex[i] + sign
			_, voxel := getBlockAndVoxelFromGlobalVoxelIndexIfExists(tsdfLayer, neighborIndex)
			if voxel == nil {
				return gradient, false
			}
			gradient[i] += voxel.getDistance() * float64(sign)
		}
	}
	for i := 0; i < 3; i++ {
		gradient[i] /= 2.0 * tsdfLayer.VoxelSize
	}
	return gradient, true
}

func addNormalizedPointInfo(point, normalizedPointNormal Point, infoVector *[6]float64) {
	translation := vec3.Mul(&normalizedPointNormal, &normalizedPointNormal)
	infoVector[0] += 2 * translation[0]
	infoVector[1] += 2 * translation[1]
	infoVector[2] += 2 * translation[2]
	infoVector[3] += 2 * (point[1]*point[1]*normalizedPointNormal[2]*normalizedPointNormal[2] +
		point[2]*point[2]*normalizedPointNormal[1]*normalizedPointNormal[1])
	infoVector[4] += 2 * (point[0]*point[0]*normalizedPointNormal[2]*normalizedPointNormal[2] +
		point[2]*point[2]*normalizedPointNormal[0]*normalizedPointNormal[0])
	infoVector[5] += 2 * (point[0]*point[0]*normalizedPointNormal[1]*normalizedPointNormal[1] +
		point[1]*point[1]*normalizedPointNormal[0]*normalizedPointNormal[0])
}

func computeTargetPoint(pointG, voxelCenter, gradient Point, distance float64) Point {
	delta := vec3.Sub(&pointG, &voxelCenter)
	distance += vec3.Dot(&gradient, &delta)
	return vec3.Sub(&pointG, gradient.Scale(distance))
}

func matchPoints(tsdfLayer *TsdfLayer, pose Transform, pointCloud *PointCloud) ([]Point, []Point) {
	kMinGradMagnitude := 0.1
	infoVector := [6]float64{}

	var srcPoints []Point
	var tgtPoints []Point

	for _, point := range pointCloud.Points {
		pointG := pose.transformPoint(point)

		globalVoxelIndex := getGridIndexFromPoint(pointG, tsdfLayer.VoxelSizeInv)
		block, voxel := getBlockAndVoxelFromGlobalVoxelIndexIfExists(tsdfLayer, globalVoxelIndex)
		if block == nil || voxel == nil {
			continue
		}
		distance := voxel.getDistance()
		if distance >= tsdfLayer.VoxelSize*4 {
			continue
		}
		gradient, ok := getGradient(tsdfLayer, globalVoxelIndex)
		if !ok {
			continue
		}
		if gradient.LengthSqr() < kMinGradMagnitude {
			continue
		}
		gradient = gradient.Normalized()

		addNormalizedPointInfo(vec3.Sub(&pointG, &pose.Translation), gradient, &infoVector)

		voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, tsdfLayer.VoxelSize)
		tgtPoint := computeTargetPoint(pointG, voxelCenter, gradient, distance)

		srcPoints = append(srcPoints, pointG)
		tgtPoints = append(tgtPoints, tgtPoint)
	}
	return srcPoints, tgtPoints
}

func vectorToTransform(vector [6]float64) Transform {
	translation := Point{vector[0], vector[1], vector[2]}
	rotation := quaternion.FromEulerAngles(vector[3], vector[4], vector[5])
	return Transform{Rotation: rotation, Translation: translation}
}

func GetIcpTransform(tsdfLayer *TsdfLayer, pose Transform, pointCloud PointCloud) Transform {
	defer TimeTrack(time.Now(), "ICP")

	// Shuffle the point cloud.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(pointCloud.Points), func(i, j int) {
		pointCloud.Points[i], pointCloud.Points[j] = pointCloud.Points[j], pointCloud.Points[i]
		pointCloud.Colors[i], pointCloud.Colors[j] = pointCloud.Colors[j], pointCloud.Colors[i]
	})

	transform := [6]float64{kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon}

	// TODO: Mini-batch size from config.
	// TODO: Cache gradients for voxels.
	miniBatchSize := 100
	chunks := splitPointCloud(&pointCloud, len(pointCloud.Points)/miniBatchSize)
	for _, pC := range chunks {
		src, tgt := matchPoints(tsdfLayer, pose, &pC)
		matchPercentage := float64(len(src)) / float64(len(pC.Points))
		if matchPercentage < 0.8 {
			continue
		}
		err := stepICP(src, tgt, &transform)
		if err != nil {
			continue
		}
	}
	temp := vectorToTransform(transform)
	fmt.Println("ICP transform:", temp)
	return temp
}
