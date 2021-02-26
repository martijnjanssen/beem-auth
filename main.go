package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"beem-auth/internal/pb"
	"beem-auth/internal/pkg/database"
	"beem-auth/internal/pkg/middleware"
	"beem-auth/internal/pkg/service"
	"beem-auth/internal/pkg/util/email"
	"beem-auth/internal/pkg/web"
)

type Config struct {
	Port string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string

	Profile string

	SendinblueKey   string
	MailCallbackURL string
}

func main() {
	log.Println("Running GRPC with version", grpc.Version)

	var conf Config
	err := envconfig.Process("beemauth", &conf)
	if err != nil {
		log.Fatalf("failed to parse environment variables: %v", err)
	}

	db, err := database.Connect(conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPassword, conf.DbName)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = database.ApplyMigrations(db)
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	// Listen to tcp port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		// Interceptors are executed from left to right (or top to bottom as seen here)
		grpc_middleware.WithUnaryServerChain(
			middleware.NewTransactionInterceptor(db),
			// TODO: needs a panic recoveryhandler with a custom RecoveryHandlerFuncContext after the
			// TransactionInterceptor
			// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/recovery/interceptors.go
		),
	)

	sendinblue := email.NewSendinblue(conf.SendinblueKey)

	pb.RegisterAccountServiceServer(grpcServer, service.NewAccountController(sendinblue, conf.MailCallbackURL))

	// Enable reflection for https://github.com/grpc/grpc/blob/master/doc/command_line_tool.md#usage
	if conf.Profile == "dev" {
		reflection.Register(grpcServer)
		csrf.Secure(false)
	}

	web.SetDB(db)
	cancelListen := web.ListenHTTP()

	log.Printf("started...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		cancelListen()
	}
}
