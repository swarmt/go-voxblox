package main

import (
	"encoding/binary"
	"math"
	"sync"
	"time"

	"go-voxblox/voxblox"

	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/ungerik/go3d/float64/quaternion"

	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

// transformToQuaternion converts a goroslib TransformStamped to a voxblox Quaternion
func transformToQuaternion(msg *geometry_msgs.TransformStamped) quaternion.T {
	return quaternion.T{
		msg.Transform.Rotation.X,
		msg.Transform.Rotation.Y,
		msg.Transform.Rotation.Z,
		msg.Transform.Rotation.W,
	}
}

// vector3ToPoint converts a goroslib Vector3 to a voxblox Point
func vector3ToPoint(msg *geometry_msgs.Vector3) voxblox.Point {
	return voxblox.Point{msg.X, msg.Y, msg.Z}
}

// interpolatePoints interpolates between two Points
func interpolatePoints(p1, p2 voxblox.Point, f float64) voxblox.Point {
	return voxblox.Point{
		p1[0] + (p2[0]-p1[0])*f,
		p1[1] + (p2[1]-p1[1])*f,
		p1[2] + (p2[2]-p1[2])*f,
	}
}

// TransformQueue is a queue of goroslib TransformStamped messages
type TransformQueue struct {
	StaticTransform voxblox.Transformation
	sync.Mutex
	transforms []*geometry_msgs.TransformStamped
}

// NewTransformQueue returns a new TransformQueue.
func NewTransformQueue(staticTransform voxblox.Transformation) *TransformQueue {
	return &TransformQueue{
		StaticTransform: staticTransform,
	}
}

// addTransform adds a transform to the TransformQueue.
func (t *TransformQueue) addTransform(transform *geometry_msgs.TransformStamped) {
	t.Lock()
	defer t.Unlock()
	t.transforms = append(t.transforms, transform)
}

// interpolateTransform interpolates a transform from the TransformQueue given a timestamp.
func (t *TransformQueue) interpolateTransform(
	timeStamp time.Time,
) *voxblox.Transformation {
	t.Lock()
	defer t.Unlock()

	// If there are no transforms, return nil.
	if len(t.transforms) == 0 {
		return nil
	}

	// If the timeStamp is before the first t0, return false.
	if timeStamp.Before(t.transforms[0].Header.Stamp) {
		return nil
	}

	// Get the t0 before and after the timeStamp.
	var i int
	found := false
	for i = 1; i < len(t.transforms); i++ {
		if t.transforms[i-1].Header.Stamp.Before(timeStamp) &&
			t.transforms[i].Header.Stamp.After(timeStamp) {
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	// Check that the timestamp is not more than 100ms from either the previous or next t0
	if t.transforms[i-1].Header.Stamp.Add(100*time.Millisecond).Before(timeStamp) ||
		t.transforms[i].Header.Stamp.Add(-100*time.Millisecond).After(timeStamp) {
		return nil
	}

	ts0 := t.transforms[i-1].Header.Stamp
	q0 := transformToQuaternion(t.transforms[i-1])
	p0 := vector3ToPoint(&t.transforms[i-1].Transform.Translation)
	ts1 := t.transforms[i].Header.Stamp
	q1 := transformToQuaternion(t.transforms[i])
	p1 := vector3ToPoint(&t.transforms[i].Transform.Translation)
	tsd := ts1.Sub(ts0)

	// Get the interpolation factor.
	f := float64(timeStamp.Sub(ts0)) / float64(tsd)

	t0 := voxblox.Transformation{}

	// Interpolate the transformation.
	t0.Rotation = quaternion.Slerp(&q0, &q1, f)
	t0.Translation = interpolatePoints(p0, p1, f)

	// Apply the static transformation.
	t1 := voxblox.CombineTransformations(&t0, &t.StaticTransform)

	// Remove all the transforms that are older than the timeStamp.
	for len(t.transforms) > 0 && t.transforms[0].Header.Stamp.Before(timeStamp) {
		t.transforms = t.transforms[1:]
	}

	return &t1
}

// float32ToRGB converts a PCL float32 color to uint8 RGB
func float32ToRGB(f float32) voxblox.Color {
	var r, g, b uint8
	i := math.Float32bits(f)
	r = uint8(i & 0x00FF0000 >> 16)
	g = uint8((i & 0x0000FF00) >> 8)
	b = uint8(i & 0x000000FF)
	return voxblox.Color{r, g, b}
}

type XYZRGB struct {
	X, Y, Z float32
	_       float32
	RGB     float32
}

// pointCloud2ToPointCloud converts a goroslib PointCloud2 to a voxblox PointCloud without reflection
func pointCloud2ToPointCloud(msg *sensor_msgs.PointCloud2) voxblox.PointCloud {
	defer voxblox.TimeTrack(time.Now(), "Convert PointCloud2")

	// Unpacks a PointCloud2 message into a voxblox PointCloud.
	pointCloud := voxblox.PointCloud{}
	pointCloud.Points = make([]voxblox.Point, 0, int(msg.Width)*int(msg.Height))
	pointCloud.Colors = make([]voxblox.Color, 0, int(msg.Width)*int(msg.Height))

	// TODO: Make this dynamic based on the message fields.
	for v := 0; v < int(msg.Height); v++ {
		offset := int(msg.RowStep) * v
		for u := 0; u < int(msg.Width); u++ {
			var p XYZRGB
			p.X = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset : offset+4]))
			p.Y = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset+4 : offset+8]))
			p.Z = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset+8 : offset+12]))
			p.RGB = math.Float32frombits(
				binary.LittleEndian.Uint32(msg.Data[offset+16 : offset+20]),
			)

			if !math.IsNaN(float64(p.X)) &&
				!math.IsNaN(float64(p.Y)) &&
				!math.IsNaN(float64(p.Z)) &&
				!math.IsNaN(float64(p.RGB)) {
				pointCloud.Points = append(pointCloud.Points, voxblox.Point{
					float64(p.X),
					float64(p.Y),
					float64(p.Z),
				})
				pointCloud.Colors = append(pointCloud.Colors, float32ToRGB(p.RGB))
			}
			offset += int(msg.PointStep)
		}
	}
	pointCloud.Width = int(msg.Width)
	pointCloud.Height = int(msg.Height)

	return pointCloud
}
