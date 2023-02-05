package gapi

import (
	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
	"github.com/thiri-lwin/thiri-bank/pb"
)

type Server struct {
	pb.UnimplementedThiriBankServer
	store *db.Store
}

// NewServer creates a new gRPC server
func NewServer(store *db.Store) *Server {
	return &Server{
		store: store,
	}
}
