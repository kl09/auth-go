package pg_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/kl09/auth-go/internal/pg"
)

var (
	PostgresHost           = getEnv("POSTGRES_HOST", "localhost")
	PostgresPort           = getEnv("POSTGRES_PORT", "5432")
	PostgresDB             = getEnv("POSTGRES_DB", "auth_test")
	PostgresDBTest         = getEnv("POSTGRES_DB_TEST", "auth_test_only")
	PostgresUser           = getEnv("POSTGRES_USER", "auth")
	PostgresPassword       = getEnv("POSTGRES_PASSWORD", "auth")
	PostgresConnectTimeout = getEnv("POSTGRES_CONNECT_TIMEOUT", "3")

	PostgresSys = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s connect_timeout=%s sslmode=disable",
		PostgresUser, PostgresPassword, PostgresHost, PostgresPort, PostgresDB, PostgresConnectTimeout)

	PostgresTest = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s connect_timeout=%s sslmode=disable",
		PostgresUser, PostgresPassword, PostgresHost, PostgresPort, PostgresDBTest, PostgresConnectTimeout)
)

func setUp(t *testing.T) {
	t.Helper()

	clearSQLDb(t)
}

func clearSQLDb(t *testing.T) {
	var err error

	pool, err := sql.Open("postgres", PostgresSys)
	if err != nil {
		t.Fatal("can't connect to db")
	}
	defer func() { _ = pool.Close() }()

	_, err = pool.Exec("DROP DATABASE IF EXISTS " + PostgresDBTest)
	if err != nil {
		t.Fatal(err)
	}
	_, err = pool.Exec("CREATE DATABASE " + PostgresDBTest)
	if err != nil {
		t.Fatal(err)
	}

	// Create schema
	c := pg.NewClient()
	if err = c.Open(PostgresTest); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	if err = c.Schema(); err != nil {
		t.Fatal(err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}
	return value
}
