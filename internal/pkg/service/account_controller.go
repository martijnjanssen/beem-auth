package service

import (
	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"context"
	"log"

	"beem-auth/internal/pkg/util"

	"beem-auth/internal/pkg/middleware"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AccountController implements the GRPC AccountService
type accountController struct {
	pb.UnimplementedAccountServiceServer
}

func NewAccountController() pb.AccountServiceServer {
	return &accountController{}
}

// Create creates a new user
func (a accountController) Create(ctx context.Context, req *pb.AccountCreateRequest) (*empty.Empty, error) {
	tx := middleware.GetContextTx(ctx)

	hashPassword, err := util.HashAndSalt(req.GetPassword())
	if err != nil {
		log.Printf("unable to hash password: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	err = database.UserAdd(ctx, tx, req.GetEmail(), hashPassword)
	if err != nil {
		log.Printf("unable to create account: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &empty.Empty{}, nil
}
