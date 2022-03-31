package main

import (
	"context"
	"go-voxblox/proto"
	"go-voxblox/voxblox"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/aler9/goroslib/pkg/msgs/sensor_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go3d/float64/quaternion"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// countResponses counts the number of responses received from the server
func countResponses(stream proto.MeshService_GetMeshBlocksClient) int {
	count := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		count++
	}
	return count
}

// TestGetMeshBlocks tests the GetMeshBlocks RPC
func TestGetMeshBlocks(t *testing.T) {
	config, _ := voxblox.ReadConfig("testdata/test.yaml")
	tsdfLayer := voxblox.NewTsdfLayer(config.VoxelSize, config.VoxelsPerSide)
	tsdfIntegrator := voxblox.NewFastTsdfIntegrator(&config, tsdfLayer)
	meshLayer := voxblox.NewMeshLayer(tsdfLayer)
	meshIntegrator := voxblox.NewMeshIntegrator(config, tsdfLayer, meshLayer)
	meshServer := NewMeshServer(&meshIntegrator)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterMeshServiceServer(s, meshServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Spin up an in memory server
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	// Create a client
	client := proto.NewMeshServiceClient(conn)
	resp, err := client.GetMeshBlocks(ctx, &proto.GetMeshRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 0, countResponses(resp))

	data, err := os.ReadFile("testdata/PointCloud2.bin")
	assert.NoError(t, err)
	pointCloud2 := sensor_msgs.PointCloud2{
		Header:    std_msgs.Header{},
		Height:    480,
		Width:     640,
		PointStep: 32,
		RowStep:   20480,
	}
	pointCloud2.Data = data
	pointCloud := PointCloud2ToPointCloud(&pointCloud2)
	tsdfIntegrator.IntegratePointCloud(
		voxblox.Transform{
			Translation: voxblox.Point{6, -0.2, -2},
			Rotation: quaternion.T{
				-0.03534060950936696,
				0.03534060950936697,
				0.7062230818371107,
				0.7062230818371108,
			},
		},
		pointCloud,
	)

	time.Sleep(100 * time.Millisecond)

	resp, err = client.GetMeshBlocks(ctx, &proto.GetMeshRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 4, countResponses(resp))

	resp, err = client.GetMeshBlocks(ctx, &proto.GetMeshRequest{})
	assert.NoError(t, err)
	assert.Equal(t, 0, countResponses(resp))
}
