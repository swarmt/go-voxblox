package main

import (
	"fmt"
	"go-voxblox/voxblox"
	"os"
	"os/signal"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

// onPointCloud2 is called when a PointCloud2 message is received.
func onPointCloud2(
	msg *sensor_msgs.PointCloud2,
	tsdfIntegrator voxblox.TsdfIntegrator,
	tf *TransformQueue,
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
	config, err := voxblox.ReadConfig("voxblox.yaml")
	if err != nil {
		panic(err)
	}

	// Create a node and connect to the master
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          "go-voxblox",
		MasterAddress: config.RosMaster,
	})
	if err != nil {
		panic(err)
	}
	defer n.Close()

	// Transformer.
	staticTransform := voxblox.Transformation{
		Rotation:    config.Rotation,
		Translation: config.Translation,
	}
	tf := NewTransformQueue(staticTransform)

	// Create integrators
	tsdfLayer := voxblox.NewTsdfLayer(config.VoxelSize, config.VoxelsPerSide)
	tsdfIntegrator := voxblox.NewFastTsdfIntegrator(&config, tsdfLayer)

	meshLayer := voxblox.NewMeshLayer(tsdfLayer)
	meshIntegrator := voxblox.NewMeshIntegrator(config, tsdfLayer, meshLayer)

	// Create a PointCloud2 subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: config.TopicPointCloud2,
		Callback: func(msg *sensor_msgs.PointCloud2) {
			onPointCloud2(msg, tsdfIntegrator, tf)
		},
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Create a Transform subscriber
	sub, err = goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: config.TopicTransform,
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

	// TODO: This is temporary. Should be a service.
	meshIntegrator.IntegrateMesh()
	voxblox.WriteMeshLayerToObjFiles(meshLayer, "output/cow_lady")
}
