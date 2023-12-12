package main

import (
	"fmt"
	"go-voxblox/proto"
	"go-voxblox/voxblox"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
)

// onPointCloud2 is called when a PointCloud2 message is received.
// Converts the message to a Voxblox PointCloud and integrates it.
func onPointCloud2(
	msg *sensor_msgs.PointCloud2,
	tsdfIntegrator voxblox.TsdfIntegrator,
	tf *TransformListener,
) {
	transform, err := tf.LookupTransform(msg.Header.Stamp)
	if err != nil {
		log.Println(err)
		return
	}
	voxbloxPointCloud := PointCloud2ToPointCloud(msg)
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

	// Transformer
	tfListener := NewTransformListener(voxblox.Transform{
		Rotation:    config.Rotation,
		Translation: config.Translation,
	})

	// Integrators
	tsdfLayer := voxblox.NewTsdfLayer(config.VoxelSize, config.VoxelsPerSide)
	tsdfIntegrator := voxblox.NewFastTsdfIntegrator(&config, tsdfLayer)
	meshLayer := voxblox.NewMeshLayer(tsdfLayer)
	meshIntegrator := voxblox.NewMeshIntegrator(config, tsdfLayer, meshLayer)

	// PointCloud2 subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: config.TopicPointCloud2,
		Callback: func(msg *sensor_msgs.PointCloud2) {
			onPointCloud2(msg, tsdfIntegrator, tfListener)
		},
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// Transform subscriber
	sub, err = goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:  n,
		Topic: config.TopicTransform,
		Callback: func(msg *geometry_msgs.TransformStamped) {
			tfListener.addTransform(msg)
		},
	})
	if err != nil {
		panic(err)
	}
	defer sub.Close()

	// gRPC mesh server
	meshServer := NewMeshServer(&meshIntegrator)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterMeshServiceServer(grpcServer, meshServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()
	defer grpcServer.Stop()

	// Wait for CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// TODO: This is temporary.
	meshIntegrator.Integrate()
	voxblox.WriteMeshLayerToObjFiles(meshLayer, "output")
}
