package users

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/atlantacoven/door-lock/server/database"
)

func findByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}

	db := database.Get(ctx)
	q, args := sq.Select("*").From("users").Where(sq.Eq{"email": email}).Limit(1).MustSql()
	if err := db.GetContext(ctx, &user, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func create(ctx context.Context, user *User) error {
	db := database.Get(ctx)
	q, args := sq.Insert("users").
		Columns("name", "email", "encrypted_password").
		Values(user.Name, user.Email, user.EncryptedPassword).
		Suffix(`RETURNING "id"`).
		MustSql()
	return db.QueryRowContext(ctx, q, args...).Scan(&user.Id)
}
