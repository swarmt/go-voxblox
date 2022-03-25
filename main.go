package main

import (
	"encoding/binary"
	"fmt"
	"go-voxblox/voxblox"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/ungerik/go3d/float64/quaternion"
)

var (
	tsdfConfig      voxblox.TsdfConfig
	tsdfLayer       *voxblox.TsdfLayer
	tsdfIntegrator  voxblox.TsdfIntegrator
	meshConfig      voxblox.MeshConfig
	meshLayer       *voxblox.MeshLayer
	meshIntegrator  voxblox.MeshIntegrator
	transformsMutex sync.Mutex
	transforms      []*geometry_msgs.TransformStamped
)

type XYZRGB struct {
	X, Y, Z float32
	_       float32
	RGB     float32
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

// pointCloud2ToPointCloud converts a goroslib PointCloud2 to a voxblox PointCloud without reflection
func pointCloud2ToPointCloud(msg *sensor_msgs.PointCloud2) voxblox.PointCloud {
	defer voxblox.TimeTrack(time.Now(), "Convert PointCloud2")

	// Unpacks a PointCloud2 message into a voxblox PointCloud.
	pointcloud := voxblox.PointCloud{}
	pointcloud.Points = make([]voxblox.Point, 0, int(msg.Width)*int(msg.Height))
	pointcloud.Colors = make([]voxblox.Color, 0, int(msg.Width)*int(msg.Height))

	// TODO: Make this dynamic based on the message fields.
	for v := 0; v < int(msg.Height); v++ {
		offset := int(msg.RowStep) * v
		for u := 0; u < int(msg.Width); u++ {
			var p XYZRGB
			p.X = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset : offset+4]))
			p.Y = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset+4 : offset+8]))
			p.Z = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset+8 : offset+12]))
			p.RGB = math.Float32frombits(binary.LittleEndian.Uint32(msg.Data[offset+16 : offset+20]))

			if !math.IsNaN(float64(p.X)) &&
				!math.IsNaN(float64(p.Y)) &&
				!math.IsNaN(float64(p.Z)) &&
				!math.IsNaN(float64(p.RGB)) {
				pointcloud.Points = append(pointcloud.Points, voxblox.Point{
					float64(p.X),
					float64(p.Y),
					float64(p.Z),
				})
				pointcloud.Colors = append(pointcloud.Colors, float32ToRGB(p.RGB))
			}
			offset += int(msg.PointStep)
		}
	}
	pointcloud.Width = int(msg.Width)
	pointcloud.Height = int(msg.Height)

	return pointcloud
}

// transformToQuaternion converts a goroslib TransformStamped to a voxblox Quaternion
func transformToQuaternion(msg *geometry_msgs.TransformStamped) quaternion.T {
	return quaternion.T{
		msg.Transform.Rotation.X,
		msg.Transform.Rotation.Y,
		msg.Transform.Rotation.Z,
		msg.Transform.Rotation.W,
	}
}

func vector3ToPoint(msg *geometry_msgs.Vector3) voxblox.Point {
	return voxblox.Point{msg.X, msg.Y, msg.Z}
}

func interpolatePoints(p1, p2 voxblox.Point, f float64) voxblox.Point {
	return voxblox.Point{
		p1[0] + (p2[0]-p1[0])*f,
		p1[1] + (p2[1]-p1[1])*f,
		p1[2] + (p2[2]-p1[2])*f,
	}
}

// interpolateTransform interpolates the transform between the
// two bounding transforms in the transforms slice.
// TODO: Use tf2 to interpolate the transforms.
func interpolateTransform(stamp time.Time, transform *voxblox.Transformation) bool {
	transformsMutex.Lock()
	defer transformsMutex.Unlock()

	// If the stamp is before the first transform, return false.
	if stamp.Before(transforms[0].Header.Stamp) {
		return false
	}

	var i int
	found := false
	for i = 1; i < len(transforms); i++ {
		if transforms[i-1].Header.Stamp.Before(stamp) && transforms[i].Header.Stamp.After(stamp) {
			found = true
			break
		}
	}

	if !found {
		return false
	}

	// Check that the timestamp is not more than 100ms from either the previous or next transform
	if transforms[i-1].Header.Stamp.Add(100*time.Millisecond).Before(stamp) ||
		transforms[i].Header.Stamp.Add(-100*time.Millisecond).After(stamp) {
		return false
	}

	ts0 := transforms[i-1].Header.Stamp
	q0 := transformToQuaternion(transforms[i-1])
	p0 := vector3ToPoint(&transforms[i-1].Transform.Translation)
	ts1 := transforms[i].Header.Stamp
	q1 := transformToQuaternion(transforms[i])
	p1 := vector3ToPoint(&transforms[i].Transform.Translation)
	tsd := ts1.Sub(ts0)

	// Get the interpolation factor.
	f := float64(stamp.Sub(ts0)) / float64(tsd)

	// Interpolate the quaternion.
	transform.Rotation = quaternion.Slerp(&q0, &q1, f)

	// Interpolate the translation.
	transform.Translation = interpolatePoints(p0, p1, f)

	// Static transform.
	// TODO: Pull this out of a config file or switch to TF2.
	staticTransform := voxblox.Transformation{
		Rotation:    quaternion.T{0.0924132, 0.0976455, 0.0702949, 0.9884249},
		Translation: voxblox.Point{0.00114049, 0.0450936, 0.0430765},
	}

	transform.Rotation = quaternion.Mul(&transform.Rotation, &staticTransform.Rotation)
	transform.Translation = voxblox.Point{
		transform.Translation[0] + staticTransform.Translation[0],
		transform.Translation[1] + staticTransform.Translation[1],
		transform.Translation[2] + staticTransform.Translation[2],
	}

	// TODO: Transforms in a queue to keep RAM usage down. Remove old transforms.

	return true
}

// onPointCloud2 is called when a PointCloud2 message is received.
func onPointCloud2(msg *sensor_msgs.PointCloud2) {
	// Convert goroslib transform to go3d transform
	// TODO: Replace go3d with internal transform methods that work directly on messages?
	if len(transforms) == 0 {
		return
	}

	var transform voxblox.Transformation
	if interpolateTransform(msg.Header.Stamp, &transform) {

		// Convert goroslib point cloud to voxblox point cloud
		voxbloxPointCloud := pointCloud2ToPointCloud(msg)

		// Integrate
		tsdfIntegrator.IntegratePointCloud(transform, voxbloxPointCloud)

	} else {
		fmt.Println("No transform")
	}
}

func onTransform(msg *geometry_msgs.TransformStamped) {
	transformsMutex.Lock()
	defer transformsMutex.Unlock()
	transforms = append(transforms, msg)
}

func main() {
	// Create a node and connect to the master
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "go-voxblox",
		MasterAddress: "127.0.0.1:11311",
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	// TODO: Read from a config file
	tsdfConfig = voxblox.TsdfConfig{
		VoxelSize:                   0.05,
		VoxelsPerSide:               16,
		MinRange:                    0.1,
		MaxRange:                    5.0,
		AllowClearing:               true,
		AllowCarving:                true,
		WeightConstant:              false,
		WeightDropOff:               true,
		MaxWeight:                   10000.0,
		StartVoxelSubsamplingFactor: 2.0,
		ClearChecksEveryNFrames:     1,
		MaxConsecutiveRayCollisions: 2,
		Threads:                     runtime.NumCPU(),
	}

	meshConfig = voxblox.MeshConfig{
		UseColor:  true,
		MinWeight: 0.5,
	}

	// Create integrators
	tsdfLayer = voxblox.NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	tsdfIntegrator = voxblox.NewFastTsdfIntegrator(tsdfConfig, tsdfLayer)

	meshLayer = voxblox.NewMeshLayer(tsdfLayer)
	meshIntegrator = voxblox.NewMeshIntegrator(meshConfig, tsdfLayer, meshLayer)

	// Create a PointCloud2 subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     n,
		Topic:    "/camera/depth_registered/points",
		Callback: onPointCloud2,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Create a Transform subscriber
	sub, err = goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     n,
		Topic:    "/kinect/vrpn_client/estimated_transform",
		Callback: onTransform,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	meshIntegrator.IntegrateMesh()
	voxblox.WriteMeshLayerToObjFiles(meshLayer, "output/cow_lady")
}
