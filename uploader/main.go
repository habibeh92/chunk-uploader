package main

import (
	"chunk-uploader/repository/uploaderpb"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type uploaderServer struct {
	log *logrus.Logger
}

func main() {
	fmt.Println("Uploader service is listening to port 50051 ...")

	// Make a listener
	lis, err := net.Listen("tcp", "uploader_service:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// SSL config
	var opts []grpc.ServerOption

	// Make a gRPC server
	grpcServer := grpc.NewServer(opts...)
	uploaderpb.RegisterChunkServiceServer(grpcServer, &uploaderServer{log: logrus.New()})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Run the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// ChunkUpload validate and save the uploaded chunk
func (u uploaderServer) ChunkUpload(ctx context.Context, r *uploaderpb.ChunkRequest) (*uploaderpb.ChunkResponse, error) {
	sha256 := r.GetSha256()
	id := r.GetId()

	u.log.WithFields(logrus.Fields{
		"sha256":   sha256,
		"chunk-id": id,
	}).Info("incoming uploader request")

	if !u.validateImage(sha256) {
		return u.withError(status.Error(codes.NotFound, "Image not found"))
	}

	if !u.validateChunk(sha256, id) {
		return u.withError(status.Error(codes.AlreadyExists, "Chunk already exists"))
	}

	err := os.WriteFile(u.chunkPath(sha256, id), []byte(r.GetData()), 0644)
	if err != nil {
		return u.withError(status.Error(codes.Internal, "Failed to create chunk file"))
	}

	return &uploaderpb.ChunkResponse{}, nil
}

// validateImage check if image exists
func (u uploaderServer) validateImage(sha256 string) bool {
	if _, err := os.Stat(u.imagePath(sha256)); os.IsNotExist(err) {
		return false
	}

	return true
}

// imagePath get the image path
func (u uploaderServer) imagePath(sha256 string) string {
	return fmt.Sprintf("./repository/images/%s", sha256)
}

// chunkPath get the chunk path
func (u uploaderServer) chunkPath(sha256 string, id int64) string {
	return fmt.Sprintf("%s/%d.chunk", u.imagePath(sha256), id)
}

// validateChunk check if the chunk path is already exist
func (u uploaderServer) validateChunk(sha256 string, id int64) bool {
	if _, err := os.Stat(u.chunkPath(sha256, id)); !os.IsNotExist(err) {
		return false
	}

	return true
}

// withError log and return err
func (u uploaderServer) withError(err error) (*uploaderpb.ChunkResponse, error) {
	u.log.Error(err)
	return nil, err
}
