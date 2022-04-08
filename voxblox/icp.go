package voxblox

import (
	"time"

	"github.com/ungerik/go3d/float64/vec3"
	"gonum.org/v1/gonum/mat"
)

func stepICP(src, tgt []Point, initial *[6]float64) error {
	var jacobian []float64
	var residual []float64

	for i := 0; i < len(src); i++ {
		dx := tgt[i][0]
		dy := tgt[i][1]
		dz := tgt[i][2]

		sx := src[i][0]
		sy := src[i][1]
		sz := src[i][2]

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

		jacobian = append(jacobian, a1, a2, a3, a4, a5, a6)
		residual = append(residual, r)
	}

	jacobianMatrix := mat.NewDense(len(src), 6, jacobian)
	residualVector := mat.NewVecDense(len(src), residual)

	var update mat.VecDense
	err := update.SolveVec(jacobianMatrix, residualVector)
	if err != nil {
		return err
	}

	for i := 0; i < update.Len(); i++ {
		initial[i] -= update.At(i, 0)
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

func matchPoints(tsdfLayer *TsdfLayer, pointCloud *PointCloud, pose *Transform) ([]Point, []Point) {
	kMinGradMagnitude := 0.1
	infoVector := [6]float64{kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon}

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
		if distance <= 0 {
			continue
		}
		gradient, ok := getGradient(tsdfLayer, globalVoxelIndex)
		if !ok {
			continue
		}
		if gradient.LengthSqr() < kMinGradMagnitude {
			continue
		}
		gradient.Normalize()

		addNormalizedPointInfo(vec3.Sub(&pointG, &pose.Translation), gradient, &infoVector)

		voxelCenter := getCenterPointFromGridIndex(globalVoxelIndex, tsdfLayer.VoxelSize)

		delta := vec3.Sub(&pointG, &voxelCenter)
		distance += vec3.Dot(&delta, &gradient)

		srcPoints = append(srcPoints, pointG)
		tgtPoints = append(tgtPoints, vec3.Sub(&pointG, gradient.Scale(distance)))
	}
	return srcPoints, tgtPoints
}

func GetIcpTransform(tsdfLayer *TsdfLayer, pose Transform, pointCloud PointCloud) Transform {
	defer TimeTrack(time.Now(), "ICP")

	src, tgt := matchPoints(tsdfLayer, &pointCloud, &pose)
	// Check match percentage
	if len(src) < 1000 {
		return pose
	}

	transform := [6]float64{kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon, kEpsilon}
	_ = stepICP(src, tgt, &transform)

	return Transform{}
}
