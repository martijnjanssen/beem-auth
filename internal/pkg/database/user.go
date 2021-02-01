package database

type User struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}

func UserAdd(db Queryer, email string, password string) error {
	_, err := db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	if err != nil {
		return dbAccessError(err)
	}

	return nil
}

func UserGetOnEmail(db Queryer, email string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT email FROM users WHERE email=$1", email)
	if err != nil {
		return nil, dbAccessError(err)
	}

	return user, nil
}
