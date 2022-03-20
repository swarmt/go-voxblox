package main

import (
	"go-voxblox/voxblox"
	"os"
	"os/signal"
	"runtime"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

var (
	tsdfConfig     voxblox.TsdfConfig
	tsdfLayer      *voxblox.TsdfLayer
	tsdfIntegrator voxblox.TsdfIntegrator
	meshConfig     voxblox.MeshConfig
	meshLayer      *voxblox.MeshLayer
)

func onMessage(msg *sensor_msgs.PointCloud2) {
	tsdfIntegrator.IntegratePointCloud(voxblox.Transformation{}, voxblox.PointCloud{})
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
		VoxelSize:          0.1,
		VoxelsPerSide:      16,
		MinRange:           0.1,
		MaxRange:           5.0,
		TruncationDistance: 0.1 * 4.0,
		AllowClearing:      true,
		AllowCarving:       true,
		ConstWeight:        false,
		MaxWeight:          10000.0,
		Threads:            runtime.NumCPU(),
	}

	meshConfig = voxblox.MeshConfig{
		UseColor:  true,
		MinWeight: 2.0,
	}

	// Create integrators
	tsdfLayer = voxblox.NewTsdfLayer(tsdfConfig.VoxelSize, tsdfConfig.VoxelsPerSide)
	tsdfIntegrator = &voxblox.MergedTsdfIntegrator{Config: tsdfConfig, Layer: tsdfLayer}

	// Create a subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     n,
		Topic:    "test_topic",
		Callback: onMessage,
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
