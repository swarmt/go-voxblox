package main

import (
	"fmt"
	"go-voxblox/proto"
	"go-voxblox/voxblox"
	"log"
)

// MeshServer is used to implement gRPC Server
type MeshServer struct {
	proto.UnimplementedMeshServiceServer
	meshIntegrator *voxblox.MeshIntegrator
}

// NewMeshServer creates a new MeshServer
func NewMeshServer(meshIntegrator *voxblox.MeshIntegrator) *MeshServer {
	return &MeshServer{
		meshIntegrator: meshIntegrator,
	}
}

// GetMeshBlocks streams the glTF binary data over gRPC
func (s MeshServer) GetMeshBlocks(
	in *proto.GetMeshRequest,
	srv proto.MeshService_GetMeshBlocksServer,
) error {
	s.meshIntegrator.Integrate()
	for _, meshBlock := range s.meshIntegrator.MeshLayer.GetBlocks() {
		if !meshBlock.HasData() {
			continue
		}
		buf, err := meshBlock.Gltf()
		if err != nil {
			log.Print(err)
			continue
		}
		err = srv.Send(&proto.GetMeshResult{
			Index: fmt.Sprintf(meshBlock.String()),
			Bytes: buf.Bytes(),
		})
		meshBlock.Clear()
		if err != nil {
			log.Print(err)
			return err
		}

	}
	return nil
}
