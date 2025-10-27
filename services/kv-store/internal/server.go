package internal

import (
	"context"
	kv "kv-store/internal/proto/kv/v1"
	"kv-store/internal/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	kv.UnimplementedKeyValueServer
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{store: store}
}

func (s *Server) Put(ctx context.Context, req *kv.PutRequest) (*kv.PutResponse, error) {
	s.store.Put(req.Key, req.Value)
	return &kv.PutResponse{Success: true}, nil
}

func (s *Server) Get(ctx context.Context, req *kv.GetRequest) (*kv.GetResponse, error) {
	value, exists := s.store.Get(req.Key)
	if !exists {
		return nil, status.Error(codes.NotFound, "key not found")
	}
	return &kv.GetResponse{Value: value}, nil
}

func (s *Server) Delete(ctx context.Context, req *kv.DeleteRequest) (*kv.DeleteResponse, error) {
	existed := s.store.Delete(req.Key)
	return &kv.DeleteResponse{Success: existed}, nil
}
