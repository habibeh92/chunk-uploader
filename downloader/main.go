package main

import (
	"chunk-uploader/repository/downloaderpb"
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

type downloaderServer struct {
	log *logrus.Logger
}

func main() {
	fmt.Println("Downloader service is listening to port 50052 ...")

	// Make a listener
	lis, err := net.Listen("tcp", "downloader_service:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// SSL config
	var opts []grpc.ServerOption

	// Make a gRPC server
	grpcServer := grpc.NewServer(opts...)
	downloaderpb.RegisterDownloadServiceServer(grpcServer, &downloaderServer{log: logrus.New()})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Run the gRPC server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Download merge the uploaded chunks
func (d downloaderServer) Download(ctx context.Context, r *downloaderpb.DownloadRequest) (*downloaderpb.DownloadResponse, error) {
	sha256 := r.GetSha256()
	d.log.WithFields(logrus.Fields{
		"sha256": sha256,
	}).Info("incoming downloader request")

	if !d.validateImage(sha256) {
		return d.withError(status.Error(codes.NotFound, "Image not found"))
	}

	files, err := os.ReadDir(d.imagePath(sha256))
	if err != nil {
		return d.withError(status.Error(codes.Internal, "Failed to read image"))
	}

	data := ""
	for i := 0; i < len(files); i++ {
		bytes, err := os.ReadFile(d.chunkPath(sha256, int64(i)))
		if err != nil {
			return d.withError(status.Error(codes.Internal, "Failed to read chunk file"))
		}
		data += string(bytes)
	}

	res := &downloaderpb.DownloadResponse{
		Data: data,
	}

	return res, err
}

// validateImage check if image exists
func (d downloaderServer) validateImage(sha256 string) bool {
	if _, err := os.Stat(d.imagePath(sha256)); os.IsNotExist(err) {
		return false
	}

	return true
}

// imagePath get the image path
func (d downloaderServer) imagePath(sha256 string) string {
	return fmt.Sprintf("./repository/images/%s", sha256)
}

// chunkPath get the chunk path
func (d downloaderServer) chunkPath(sha256 string, id int64) string {
	return fmt.Sprintf("%s/%d.chunk", d.imagePath(sha256), id)
}

// withError log and return err
func (d downloaderServer) withError(err error) (*downloaderpb.DownloadResponse, error) {
	d.log.Error(err)
	return nil, err
}
