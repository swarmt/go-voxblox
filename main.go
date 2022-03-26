package main

import (
	"fmt"
	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/ungerik/go3d/float64/quaternion"
	"go-voxblox/voxblox"
	"os"
	"os/signal"
	"runtime"
)

// onPointCloud2 is called when a PointCloud2 message is received.
func onPointCloud2(
	msg *sensor_msgs.PointCloud2,
	tsdfIntegrator voxblox.TsdfIntegrator,
	tf *transformQueue,
) {
	transform := tf.interpolateTransform(msg.Header.Stamp)
	if transform == nil {
		fmt.Println("No transform")
		return
	}

	// Convert goroslib point cloud to voxblox point cloud
	voxbloxPointCloud := pointCloud2ToPointCloud(msg)

	// Integrate
	tsdfIntegrator.IntegratePointCloud(*transform, voxbloxPointCloud)
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
	tsdfConfig := voxblox.TsdfConfig{
		VoxelSize:                   0.05,
		VoxelsPerSide:               16,
		MinRange:                    0.1,
		MaxRange:                    5.0,
		AllowClearing:               true,
		AllowCarving:                true,
		WeightConstant:              false,
		WeightDropOff:               true,
		MaxWeight:                   10000.0,
		StartVoxelSubsamplingFactor: 1.0,
		MaxConsecutiveRayCollisions: 2,
		Threads:                     runtime.NumCPU(),
	}

	meshConfig := voxblox.MeshConfig{
		UseColor:  true,
		MinWeight: 0.5,
	}

	// Transformer.
	// TODO: Pull this out of a config file or switch to TF2.
	// TODO: This is the Vicon to Kinect transform for the cow and lady dataset.
	staticTransform := voxblox.Transformation{
		Rotation:    quaternion.T{0.0924132, 0.0976455, 0.0702949, 0.9884249},
		Translation: voxblox.Point{0.00114049, 0.0450936, 0.0430765},
	}
	tf := NewTF(staticTransform)

	// Create integrators
	tsdfLayer := voxblox.NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	tsdfIntegrator := voxblox.NewFastTsdfIntegrator(tsdfConfig, tsdfLayer)

	meshLayer := voxblox.NewMeshLayer(tsdfLayer)
	meshIntegrator := voxblox.NewMeshIntegrator(meshConfig, tsdfLayer, meshLayer)

	// Create a PointCloud2 subscriber
	// TODO: Topic name should be configurable
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: "/camera/depth_registered/points",
		Callback: func(msg *sensor_msgs.PointCloud2) {
			onPointCloud2(msg, tsdfIntegrator, tf)
		},
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Create a Transform subscriber
	// TODO: Topic name should be configurable
	sub, err = goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: "/kinect/vrpn_client/estimated_transform",
		Callback: func(msg *geometry_msgs.TransformStamped) {
			tf.addTransform(msg)
		},
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
