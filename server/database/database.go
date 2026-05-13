package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"rabidaudio.com/coven-door/server/api"

	"github.com/integralist/go-findroot/find"

	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB is an interface implemented by both *sqlx.DB and *sqlx.Tx
// so that either can be used by our code
type DB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	Queryx(query string, args ...any) (*sqlx.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowx(query string, args ...any) *sqlx.Row

	Get(dest any, query string, args ...any) error
	Select(dest any, query string, args ...any) error

	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row

	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

func Create() (*sqlx.DB, error) {
	return CreateEnv(api.Env())
}

func CreateEnv(env api.Environment) (*sqlx.DB, error) {
	return CreateSrc(fmt.Sprintf("server/%v.db", env))
}

func CreateSrc(src string) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", src)
}

func NewMigrator(db DB) (*migrate.Migrate, error) {
	driver, err := sqlite3.WithInstance(db.(*sqlx.DB).DB, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}
	// need to find the root project directory in case we aren't running tooling
	// from there
	root, err := find.Repo()
	if err != nil {
		return nil, fmt.Errorf("find project root: %w", err)
	}
	return migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v/server/migrations", root.Path), "sqlite3", driver)
}

func createTest() *sqlx.DB {
	db := must(CreateSrc(":memory:"))

	// Wipe the database and migrate up to the latest schema
	m := must(NewMigrator(db))
	err := m.Up()
	if err != nil {
		panic(err)
	}
	return db
}

func WithDB(db DB, ctx context.Context) context.Context {
	return context.WithValue(ctx, "db", db)
}

func Get(ctx context.Context) DB {
	return ctx.Value("db").(DB)
}

var testdb *sqlx.DB

// PrepareForTest creates a clean, empty database and a transaction
// to run queries in which will be automatically rolled back at the
// end of the test. It attaches this to the test Context so it can
// be accessed as usual.
func PrepareForTest(t *testing.T) context.Context {
	t.Helper()
	if testdb == nil {
		testdb = createTest()
		// probably fine to not close since it's just in memory
	}
	tx := must(testdb.BeginTxx(t.Context(), &sql.TxOptions{}))
	t.Cleanup(func() {
		tx.Rollback()
	})
	return WithDB(tx, t.Context())
}

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}
