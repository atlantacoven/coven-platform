package users

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func findByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}

	db := ctx.Value("db").(*sqlx.DB)
	q, args := sq.Select("*").From("users").Where(sq.Eq{"email": email}).Limit(1).MustSql()
	if err := db.Get(&user, q, args); err != nil {
		return nil, err
	}
	return &user, nil
}
