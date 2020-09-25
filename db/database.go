package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	DbPool *pgxpool.Pool
}

func (r *Repository) InitDb() {
	var err error
	r.DbPool, err = pgxpool.Connect(context.Background(), "postgresql://postgres:axiom@127.0.0.1:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

// TODO: Reuse the DBPool after typecasting pgxpool.Poor into *sql.DB using the pgx/stdlib library
func (r *Repository) Migrate() {
	m, err := migrate.New(
		"file://db/migrations",
		"postgresql://postgres:axiom@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		fmt.Println(err)
	}
}
