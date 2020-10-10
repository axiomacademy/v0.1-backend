package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	dbPool *pgxpool.Pool
	logger *log.Logger
}

func (r *Repository) Init(logger *log.Logger, dbUrl string) {
	// Setup logger
	r.logger = logger

	var err error
	r.dbPool, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		r.logger.WithFields(log.Fields{
			"service": "repository",
			"dburl":   dbUrl,
			"error":   err.Error(),
		}).Fatal("Unable to connect to database")
	}

	r.logger.WithField("service", "repository").Info("Successfully initialised")
}

func (r *Repository) Close() {
	r.dbPool.Close()
}

// TODO: Reuse the DBPool after typecasting pgxpool.Poor into *sql.DB using the pgx/stdlib library
func (r *Repository) Migrate(dbUrl string) {
	m, err := migrate.New("file://db/migrations", dbUrl)
	if err != nil {
		r.logger.WithFields(log.Fields{
			"service": "repository",
			"dburl":   dbUrl,
			"error":   err.Error(),
		}).Fatal("Unable to setup migrations from db")
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			r.logger.WithField("service", "repository").Info("Database already at latest")
		} else {
			r.logger.WithFields(log.Fields{
				"service": "repository",
				"error":   err.Error(),
			}).Fatal("Unable to apply migrations to db")
		}
	}

	r.logger.WithField("service", "repository").Info("Successfully applied all database migrations")
}
