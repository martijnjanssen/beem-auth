package middleware

import (
	"beem-auth/internal/pkg/database"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey int

const (
	txContextKey contextKey = iota
)

func NewTransactionInterceptor(db *sqlx.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tx, err := db.Beginx()
		if err != nil {
			log.Printf("unable to start transaction: %s", err)
			return nil, status.Error(codes.Internal, "unable to handle request")
		}

		// Set transaction on the context
		txCtx := SetContextTx(ctx, tx)

		// Call next requesthandler to handle the request, if there was an error,
		// roll back the transaction. the resp and err here are supposed to be
		// returned to the user, unless there was an error with comitting the transaction.
		resp, err = handler(txCtx, req)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("error while rolling back transaction: %s", rollbackErr)
			}
			return
		}

		// Try to commit the transaction
		commitErr := tx.Commit()
		if commitErr != nil {
			log.Printf("unable to commit transaction: %s", commitErr)
			return nil, status.Error(codes.Internal, "")
		}

		return
	}
}

func SetContextTx(ctx context.Context, q database.Queryer) context.Context {
	return context.WithValue(ctx, txContextKey, q)
}

func GetContextTx(ctx context.Context) database.Queryer {
	return ctx.Value(txContextKey).(database.Queryer)
}
