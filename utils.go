package main

import (
	"encoding/binary"
	"fmt"
	"go-voxblox/voxblox"
	"math"
	"sync"
	"time"

	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/ungerik/go3d/float64/quaternion"

	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

func TransformStampedToTransform(msg *geometry_msgs.TransformStamped) *voxblox.Transform {
	return &voxblox.Transform{
		Translation: voxblox.Point{
			msg.Transform.Translation.X,
			msg.Transform.Translation.Y,
			msg.Transform.Translation.Z,
		},
		Rotation: quaternion.T{
			msg.Transform.Rotation.X,
			msg.Transform.Rotation.Y,
			msg.Transform.Rotation.Z,
			msg.Transform.Rotation.W,
		},
	}
}

// TransformQueue is a queue of goroslib TransformStamped messages
type TransformQueue struct {
	StaticTransform voxblox.Transform
	sync.Mutex
	transforms []*geometry_msgs.TransformStamped
}

// NewTransformQueue returns a new TransformQueue.
func NewTransformQueue(staticTransform voxblox.Transform) *TransformQueue {
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

// removePreviousTransforms removes transforms older than the given timestamp.
func (t *TransformQueue) removePreviousTransforms(timestamp time.Time) {
	for len(t.transforms) > 0 && t.transforms[0].Header.Stamp.Before(timestamp) {
		t.transforms = t.transforms[1:]
	}
}

// LookupTransform interpolates a transform from the TransformQueue given a timestamp.
func (t *TransformQueue) LookupTransform(
	timeStamp time.Time,
) (*voxblox.Transform, error) {
	t.Lock()
	defer t.Unlock()

	// If there are no transforms, return nil.
	if len(t.transforms) == 0 {
		return nil, fmt.Errorf("no transforms in queue")
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
		return nil, fmt.Errorf("no transform before and after timestamp")
	}

	// Check that the timestamp is not more than 100ms from either the previous or next t0
	if t.transforms[i-1].Header.Stamp.Add(100*time.Millisecond).Before(timeStamp) ||
		t.transforms[i].Header.Stamp.Add(-100*time.Millisecond).After(timeStamp) {
		return nil, fmt.Errorf("timestamp too far from t0 and t1")
	}

	// Calculate the interpolation factor.
	tsLower := t.transforms[i-1].Header.Stamp
	tsUpper := t.transforms[i].Header.Stamp
	tsDelta := tsUpper.Sub(tsLower)
	f := float64(timeStamp.Sub(tsLower)) / float64(tsDelta)

	t0 := voxblox.InterpolateTransform(
		*TransformStampedToTransform(t.transforms[i-1]),
		*TransformStampedToTransform(t.transforms[i]),
		f,
	)
	t1 := voxblox.ApplyTransform(&t0, &t.StaticTransform)
	t.removePreviousTransforms(timeStamp)

	return &t1, nil
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

// PointCloud2ToPointCloud converts a goroslib PointCloud2 to a voxblox PointCloud
// TODO: Make this dynamic based on the message fields.
func PointCloud2ToPointCloud(msg *sensor_msgs.PointCloud2) voxblox.PointCloud {
	defer voxblox.TimeTrack(time.Now(), "Convert PointCloud2")

	pointCloud := voxblox.PointCloud{}
	pointCloud.Points = make([]voxblox.Point, 0, int(msg.Width)*int(msg.Height))
	pointCloud.Colors = make([]voxblox.Color, 0, int(msg.Width)*int(msg.Height))

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
