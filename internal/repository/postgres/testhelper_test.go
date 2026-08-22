// internal/repository/postgres/testhelper_test.go
package postgres_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/taskdb_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test db: %v", err)
	}

	// Clean slate before each test — don't rely on test execution order.
	t.Cleanup(func() {
		db.Exec("TRUNCATE TABLE tasks")
		db.Close()
	})

	return db
}
