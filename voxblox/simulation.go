package voxblox

import (
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
	}
}

func (w SimulationWorld) AddObject(object interface{}) {
	w.Objects = append(w.Objects, object)
}

func (w SimulationWorld) AddGroundLevel(f float64) {
}

type Object interface {
	Center() Point
	DistanceToPoint(Point) float64
}

type Cylinder struct {
	Center Point
	Radius float64
	Height float64
}

func (c Cylinder) DistanceToPoint(point Point) float64 {
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

func (c Cylinder) RayIntersection() bool {
	return false
}

type Plane struct {
	Center Point
	Normal Point
}

func (plane Plane) DistanceToPoint(point Point) float64 {
	d := -vec3.Dot(plane.Normal.asVec3(), point.asVec3())
	normalized := plane.Normal.asVec3().Normalized()
	p := d / normalized.Length()
	distance := vec3.Dot(plane.Normal.asVec3(), plane.Center.asVec3()) - p
	return distance
}
