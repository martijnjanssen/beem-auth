package service

import (
	"context"
	"fmt"
	"log"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"beem-auth/internal/pkg/middleware"
	"beem-auth/internal/pkg/util/email"
	"beem-auth/internal/pkg/util/hash"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validateEmail string = `<html><head></head><h1>New beem-auth account</h1><p>You have created a new beem-auth account. Please validate your email address <a href=\"%s/challenge/complete?key=%s\">here</a>.</p></html>`

// AccountController implements the GRPC AccountService
type accountController struct {
	pb.UnimplementedAccountServiceServer

	mailer          email.Mailer
	mailCallbackURL string
}

func NewAccountController(mailer email.Mailer, mailCallbackURL string) pb.AccountServiceServer {
	return &accountController{mailer: mailer, mailCallbackURL: mailCallbackURL}
}

// Create creates a new user
func (a accountController) Create(ctx context.Context, req *pb.AccountCreateRequest) (*empty.Empty, error) {
	tx := middleware.GetContextTx(ctx)

	hashPassword, err := hash.HashAndSalt(req.GetPassword())
	if err != nil {
		log.Printf("unable to hash password: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	userId, err := database.UserAdd(ctx, tx, req.GetEmail(), hashPassword)
	if err != nil {
		log.Printf("unable to create account: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	key, err := database.ChallengeCreate(ctx, tx, userId)
	if err != nil {
		log.Printf("unable to create challenge: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	email := email.Email{
		Recipient: req.GetEmail(),
		Subject:   "New account",
		Content:   fmt.Sprintf(validateEmail, a.mailCallbackURL, key),
	}
	err = a.mailer.SendEmail(email)
	if err != nil {
		log.Printf("unable to send email: %s", err)
		return nil, status.Errorf(codes.Internal, "")
	}

	return &empty.Empty{}, nil
}
