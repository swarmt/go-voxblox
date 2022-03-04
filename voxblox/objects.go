package voxblox

import (
	"github.com/ungerik/go3d/float64/vec2"
	"github.com/ungerik/go3d/float64/vec3"
	"math"
)

type Object interface {
	DistanceToPoint(Point) float64
	RayIntersection(vec3.T, vec3.T, float64) (bool, Point, float64)
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
	minZ := c.Center.X - c.Height/2.0
	maxZ := c.Center.X + c.Height/2.0
	if point.Z > minZ && point.Z < maxZ {
		a := point.asVec2()
		b := c.Center.asVec2()
		distance = a.Sub(b).Length() - c.Radius
	} else if point.Z > maxZ {
		distance = math.Sqrt(
			math.Max((point.asVec2().Sub(c.Center.asVec2())).LengthSqr()-c.Radius*c.Radius, 0.0) +
				(point.Z-maxZ)*(point.Z-maxZ))

	} else {
		distance = math.Sqrt(
			math.Max((point.asVec2().Sub(c.Center.asVec2())).LengthSqr()-c.Radius*c.Radius, 0.0) +
				(point.Z-minZ)*(point.Z-minZ))
	}
	return distance
}

func (c *Cylinder) RayIntersection(
	rayOrigin vec3.T,
	rayDirection vec3.T,
	maxDistance float64,
) (bool, Point, float64) {
	var intersectPoint Point
	var intersectDist float64

	var vectorE = vec3.Sub(&rayOrigin, c.Center.asVec3())
	vectorD := rayDirection

	A := vectorD[0]*vectorD[0] + vectorD[1]*vectorD[1]
	B := 2*vectorE[0]*vectorD[0] + 2*vectorE[1]*vectorD[1]
	C := vectorE[0]*vectorE[0] + vectorE[1]*vectorE[1] - c.Radius*c.Radius

	// t = (-b +- sqrt(b^2 - 4ac))/2a
	// t only has solutions if b^2 - 4ac >= 0
	t1 := -1.0
	t2 := -1.0

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

		s := vectorD.Scaled(t3)
		q3 := vec3.Add(&vectorE, &s)
		s = vectorD.Scaled(t4)
		q4 := vec3.Add(&vectorE, &s)

		q3Head := vec2.T{q3[0], q3[1]}
		if t3 >= 0.0 && q3Head.Length() < c.Radius {
			t3Valid = true
		}
		q4Head := vec2.T{q4[0], q4[1]}
		if t4 >= 0.0 && q4Head.Length() < c.Radius {
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
		T = math.Min(T, t4)
	}

	// Intersection greater than max dist, so no intersection in the sensor range.
	if T >= maxDistance {
		return false, intersectPoint, intersectDist
	}

	iV := rayOrigin.Add(rayDirection.Scale(T))
	intersectPoint = Point{X: iV[0], Y: iV[1], Z: iV[2]}
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

func (plane *Plane) RayIntersection(
	rayOrigin vec3.T,
	rayDirection vec3.T,
	maxDistance float64,
) (bool, Point, float64) {
	return false, Point{}, 0.0
}
