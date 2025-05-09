package grpcservers

import (
	"context"
	pb "go-keeper/proto"
)

var _ pb.StorageServer = (*StorageServer)(nil)

type StorageServer struct {
	pb.UnimplementedStorageServer
}

func NewStorageServer() *StorageServer {
	return &StorageServer{}
}

func (s *StorageServer) Store(ctx context.Context, request *pb.StoreRequest) (*pb.StoreResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageServer) Load(ctx context.Context, request *pb.LoadRequest) (*pb.LoadResponse, error) {
	//TODO implement me
	panic("implement me")
}
