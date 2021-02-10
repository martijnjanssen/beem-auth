package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Valid    bool      `db:"valid"`
}

func UserAdd(ctx context.Context, db Queryer, email string, hashedPassword string) (uuid.UUID, error) {
	var id uuid.UUID
	err := db.GetContext(ctx, &id, "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", email, hashedPassword)
	if err != nil {
		return uuid.Nil, dbAccessError(err)
	}

	return id, nil
}

func UserGetOnEmail(ctx context.Context, db Queryer, email string) (*User, error) {
	user := &User{}
	err := db.GetContext(ctx, user, "SELECT id, email, valid FROM users WHERE email=$1", email)
	if err != nil {
		return nil, dbAccessError(err)
	}

	return user, nil
}

func UserSetValid(ctx context.Context, db Queryer, id uuid.UUID) error {
	res, err := db.ExecContext(ctx, "UPDATE users SET valid = TRUE WHERE id=$1", id)
	if err != nil {
		return dbAccessError(err)
	}
	if num, _ := res.RowsAffected(); num != 1 {
		return fmt.Errorf("%d users were updated instead of 1", num)
	}

	return nil
}
