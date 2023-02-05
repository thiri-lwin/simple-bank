package gapi

import (
	"context"

	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
	"github.com/thiri-lwin/thiri-bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	server.extractMetadata(ctx)
	acc, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.GetOwner(),
		Currency: req.GetCurrency(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.CreateAccountResponse{
		Account: &pb.Account{
			ID:        acc.ID,
			Owner:     acc.Owner,
			Currency:  acc.Currency,
			Balance:   acc.Balance,
			CreatedAt: timestamppb.New(acc.CreatedAt),
		},
	}, nil
}
