package voxblox

import (
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"math"
)

type ShapeType int

type SimulationWorld struct {
	VoxelSize float64
	MinBound  Point
	MaxBound  Point
	Objects   []interface{}
}

func NewSimulationWorld(voxelSize float64, minBound Point, maxBound Point) *SimulationWorld {
	return &SimulationWorld{
		VoxelSize: voxelSize,
		MinBound:  minBound,
		MaxBound:  maxBound,
		Objects:   make([]interface{}, 0),
	}
}

func (w *SimulationWorld) AddObject(object interface{}) {
	w.Objects = append(w.Objects, object)
}

func (w *SimulationWorld) AddGroundLevel(height float64) {
	ground := Plane{
		Center: Point{0.0, 0.0, height},
		Normal: Point{0.0, 0.0, 1.0},
		Color:  [3]uint8{},
	}
	w.Objects = append(w.Objects, ground)
}

func (w *SimulationWorld) GetPointCloudFromViewpoint(
	viewOrigin *Point,
	viewDirection *Point,
	cameraResolution vec2.T,
	fovHorizontalRad float64,
	maxDistance float64,
) sensor_msgs.PointCloud2 {
	focalLength := cameraResolution[0] / (2.0 * math.Tan(fovHorizontalRad/2.0))

	nominalViewDirection := Point{1.0, 0.0, 0.0}
	diff := quaternion.Vec3Diff(nominalViewDirection.asVec3(), viewDirection.asVec3())
	rotationQuaternion := diff.Normalized()

	// Iterate over all the pixels
	for u := -cameraResolution[0] / 2; u < cameraResolution[0]/2; u++ {
		for v := -cameraResolution[1] / 2; v < cameraResolution[1]/2; v++ {
			rayCameraDirection := Point{
				x: 1.0,
				y: u / focalLength,
				z: v / focalLength,
			}.asVec3()
			rotationQuaternion.RotateVec3(rayCameraDirection.Normalize())
		}
	}

	return sensor_msgs.PointCloud2{}
}

type Object interface {
	Center() Point
	DistanceToPoint(Point) float64
}

type Cylinder struct {
	Center Point
	Radius float64
	Height float64
	Color  Color
}

func (c *Cylinder) DistanceToPoint(point Point) float64 {
	// TODO: This seems like a simplified distance to cylinder.
	// TODO: May or may not matter.
	distance := 0.0
	minZ := c.Center.x - c.Height/2.0
	maxZ := c.Center.x + c.Height/2.0
	if point.z > minZ && point.z < maxZ {
		a := point.asVec2()
		b := c.Center.asVec2()
		distance = a.Sub(b).Length() - c.Radius
	} else if point.z > maxZ {
		distance = math.Sqrt(
			math.Max((point.asVec2().Sub(c.Center.asVec2())).LengthSqr()-c.Radius*c.Radius, 0.0) +
				(point.z-maxZ)*(point.z-maxZ))

	} else {
		distance = math.Sqrt(
			math.Max((point.asVec2().Sub(c.Center.asVec2())).LengthSqr()-c.Radius*c.Radius, 0.0) +
				(point.z-minZ)*(point.z-minZ))
	}
	return distance
}

func (c *Cylinder) RayIntersection(rayOrigin Point, rayDirection Point, maxDistance float64) (bool, Point, float64) {
	var intersectPoint Point
	var intersectDist float64

	vectorE := subtractPoints(rayOrigin, c.Center).asVec3()
	vectorD := rayDirection.asVec3()

	A := vectorD[0]*vectorD[0] + vectorD[1]*vectorD[1]
	B := 2*vectorE[0]*vectorD[0] + 2*vectorE[1]*vectorD[1]
	C := vectorE[0]*vectorE[0] + vectorE[1]*vectorE[1] - c.Radius*c.Radius

	// t = (-b +- sqrt(b^2 - 4ac))/2a
	// t only has solutions if b^2 - 4ac >= 0
	t1 := -1.0
	t2 := -1.0

	// Don't divide by 0
	if A < kEpsilon {
		// NOTE: Voxblox returns false here, but this is a valid intersection.
		A = kEpsilon
	}

	underSqrt := B*B - 4*A*C
	if underSqrt < 0 {
		return false, intersectPoint, intersectDist
	}
	if underSqrt <= kEpsilon {
		t1 = -B / (2 * A)
	} else {
		t1 = (-B + math.Sqrt(underSqrt)) / (2 * A)
		t2 = (-B - math.Sqrt(underSqrt)) / (2 * A)
	}

	// Check if hit is on the cylinder or end caps
	T := maxDistance
	z1 := vectorE[2] + t1*vectorD[2]
	z2 := vectorE[2] + t2*vectorD[2]

	t1Valid := false
	if t1 >= 0.0 && z1 >= -c.Height/2.0 && z1 <= c.Height/2.0 {
		t1Valid = true
	}
	t2Valid := false
	if t2 >= 0.0 && z2 >= -c.Height/2.0 && z2 <= c.Height/2.0 {
		t2Valid = true
	}

	var t3, t4 float64
	t3Valid := false
	t4Valid := false

	// Don't divide by 0
	if math.Abs(vectorD[2]) > kEpsilon {
		// t3 is the bottom end-cap, t4 is the top.
		t3 = (-c.Height/2.0 - vectorE[2]) / vectorD[2]
		t4 = (c.Height/2.0 - vectorE[2]) / vectorD[2]

		q3 := vectorE.Add(vectorD.Scale(t3))
		q4 := vectorE.Add(vectorD.Scale(t4))

		q3Head := vec2.T{q3[0], q3[1]}
		if t3 >= 0.0 && q3Head.Normalize().Length() < c.Radius {
			t3Valid = true
		}
		q4Head := vec2.T{q4[0], q4[1]}
		if t4 >= 0.0 && q4Head.Normalize().Length() < c.Radius {
			t4Valid = true
		}
	}

	if !(t1Valid || t2Valid || t3Valid || t4Valid) {
		return false, intersectPoint, intersectDist
	}
	if t1Valid {
		T = math.Min(T, t1)
	}
	if t2Valid {
		T = math.Min(T, t2)
	}
	if t3Valid {
		T = math.Min(T, t3)
	}
	if t4Valid {
		T = math.Min(T, t3)
	}

	// Intersection greater than max dist, so no intersection in the sensor range.
	if T >= maxDistance {
		return false, intersectPoint, intersectDist
	}

	iV := rayOrigin.asVec3().Add(rayDirection.asVec3().Scale(T))
	intersectPoint = Point{x: iV[0], y: iV[1], z: iV[2]}
	intersectDist = T

	return true, intersectPoint, intersectDist
}

type Plane struct {
	Center Point
	Normal Point
	Color  Color
}

func (plane *Plane) DistanceToPoint(point Point) float64 {
	norm := plane.Normal.asVec3().Normalized()
	d := -vec3.Dot(&norm, point.asVec3())
	p := d / norm.Length()
	distance := vec3.Dot(&norm, plane.Center.asVec3()) - p
	return distance
}
