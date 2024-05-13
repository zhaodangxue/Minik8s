package cri_proxy

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func getContext() context.Context {
	return context.Background()
}

func getCriGrpcClient() (conn *grpc.ClientConn, err error) {
	// Create a gRPC client connection
	conn, err = grpc.Dial("unix:///run/containerd/containerd.sock", grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

func getRuntimeServiceClient() (runtimeServiceClient cri.RuntimeServiceClient, err error) {
	// Create a gRPC client connection
	conn, err := getCriGrpcClient()
	if err != nil {
		return
	}
	// Create the runtime service client using the gRPC client connection
	runtimeServiceClient = cri.NewRuntimeServiceClient(conn)
	return
}

func getImageServiceClient() (imageSeviceClient cri.ImageServiceClient, err error) {
	// Create a gRPC client connection
	conn, err := getCriGrpcClient()
	if err != nil {
		return
	}
	// Create the runtime service client using the gRPC client connection
	imageSeviceClient = cri.NewImageServiceClient(conn)
	return
}