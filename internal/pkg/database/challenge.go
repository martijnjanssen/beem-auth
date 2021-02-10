package database

import (
	"context"
	"fmt"

	"beem-auth/internal/pkg/util/random"

	"github.com/google/uuid"
)

type Challenge struct {
	Key    string    `db:"key"`
	UserId uuid.UUID `db:"user_id"`
}

func ChallengeCreate(ctx context.Context, db Queryer, userId uuid.UUID) (string, error) {
	str, err := random.RandomHash(10)
	if err != nil {
		return "", fmt.Errorf("unable to generate random string: %w", err)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO challenges (user_id, key) VALUES ($1, $2)", userId, str)
	if err != nil {
		return "", dbAccessError(err)
	}

	return str, nil
}
