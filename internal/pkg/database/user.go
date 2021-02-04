package database

import "context"

type User struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

func UserAdd(ctx context.Context, db Queryer, email string, hashedPassword string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		return dbAccessError(err)
	}

	return nil
}

func UserGetOnEmail(ctx context.Context, db Queryer, email string) (*User, error) {
	user := &User{}
	err := db.GetContext(ctx, user, "SELECT email FROM users WHERE email=$1", email)
	if err != nil {
		return nil, dbAccessError(err)
	}

	return user, nil
}
