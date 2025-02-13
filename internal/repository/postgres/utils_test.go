package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"
)

var (
	db *sql.DB
)

func setupTestDB(t *testing.T) func() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", hostAndPort)

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("pgx", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		balance INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS inventory (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		item TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS history (
		id SERIAL PRIMARY KEY,
		sender_name TEXT NOT NULL,
		receiver_name TEXT NOT NULL,
		amount INTEGER
);
	`)
	if err != nil {
		log.Fatalf("Could not create table: %s", err)
	}

	return func() {
		if err = db.Close(); err != nil {
			log.Fatalf("Could not close database: %s", err)
		}
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func TestMain(m *testing.M) {
	teardown := setupTestDB(nil)

	code := m.Run()

	teardown()

	os.Exit(code)
}
