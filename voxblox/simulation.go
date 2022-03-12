package voxblox

import (
	"math"

	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
)

type SimulationWorld struct {
	VoxelSize float64
	MinBound  Point
	MaxBound  Point
	Objects   []Object
}

func NewSimulationWorld(voxelSize float64, minBound, maxBound Point) *SimulationWorld {
	return &SimulationWorld{
		VoxelSize: voxelSize,
		MinBound:  minBound,
		MaxBound:  maxBound,
		Objects:   make([]Object, 0),
	}
}

func (w *SimulationWorld) AddObject(object Object) {
	w.Objects = append(w.Objects, object)
}

func (w *SimulationWorld) GetPointCloudFromViewpoint(
	viewOrigin vec3.T,
	viewDirection vec3.T,
	cameraResolution vec2.T,
	fovHorizontal float64,
	maxDistance float64,
) PointCloud {
	fovHorizontalRad := fovHorizontal * math.Pi / 180.0
	focalLength := cameraResolution[0] / (2.0 * math.Tan(fovHorizontalRad/2.0))

	nominalViewDirection := vec3.T{1.0, 0.0, 0.0}
	rotationQuaternion := quaternion.Vec3Diff(&nominalViewDirection, &viewDirection)

	// Create a slice to store the points
	// TODO: Structured points
	var points []Point
	var colors []Color

	// Iterate over all the pixels
	for u := -cameraResolution[0] / 2; u < cameraResolution[0]/2; u++ {
		for v := -cameraResolution[1] / 2; v < cameraResolution[1]/2; v++ {
			rayCameraDirection := vec3.T{1.0, u / focalLength, v / focalLength}
			rotationQuaternion.RotateVec3(rayCameraDirection.Normalize())

			rayValid := false
			rayDistance := maxDistance
			// Iterate over all the objects
			for _, object := range w.Objects {
				intersects, objectIntersect, objectDistance := object.RayIntersection(
					viewOrigin,
					rayCameraDirection,
					maxDistance,
				)
				if intersects {
					if !rayValid || objectDistance < rayDistance {
						rayValid = true
						rayDistance = objectDistance
						points = append(points, objectIntersect)
						colors = append(colors, object.getColor())
					}
				}
			}
		}
	}
	return PointCloud{
		Width:  int(cameraResolution[0]),
		Height: int(cameraResolution[1]),
		Points: points,
		Colors: colors,
	}
}

func (w *SimulationWorld) GetPointCloudFromTransform(
	pose *Transformation,
	cameraRes vec2.T,
	fovH float64,
	maxDistance float64,
) PointCloud {
	viewDirection := vec3.T{1.0, 0.0, 0.0}
	pose.Rotation.RotateVec3(&viewDirection)
	return w.GetPointCloudFromViewpoint(
		pose.Position,
		viewDirection,
		cameraRes,
		fovH,
		maxDistance,
	)
}
