package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/middlewares"

	"github.com/solderneer/axiom-backend/services/match"
	"github.com/solderneer/axiom-backend/services/notifs"
)

const defaultPort = "8080"
const defaultDbUrl = "postgresql://postgres:axiom@127.0.0.1:5432/postgres?sslmode=disable"
const defaultSecret = "password"

func main() {
	// Setup logger
	var logger = log.New()

	// Get default environment variables
	port := os.Getenv("PORT")
	if port == "" {
		log.WithFields(log.Fields{
			"default_port": defaultPort,
		}).Warn("No PORT environment variable, using default")
		port = defaultPort
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.WithFields(log.Fields{
			"default_db_url": defaultDbUrl,
		}).Warn("No DB_URL environment variable, using default")
		dbUrl = defaultDbUrl
	}

	secret := os.Getenv("SERVER_SECRET")
	if secret == "" {
		log.WithFields(log.Fields{
			"default_secret": defaultSecret,
		}).Warn("No SERVER_SECRET environment variable, using default")
	}

	fbCredentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if fbCredentials == "" {
		log.Error("No GOOGLE_APPLICATION_CREDENTIALS environment variable. Firebase will misbehave")
	}

	repo := db.Repository{}
	repo.Init(logger, dbUrl)
	repo.Migrate(dbUrl)

	defer repo.Close()

	// Initialising all services
	ns := notifs.NotifService{}
	ns.Init(logger)

	ms := match.MatchService{}
	ms.Init(logger, secret, &ns, &repo)

	// Binding services to resolver
	resolver := graph.Resolver{
		Secret: secret,
		Logger: logger,
		Repo:   &repo,
		Ns:     &ns,
		Ms:     &ms,
	}

	graphSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))

	r := mux.NewRouter()
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphSrv)

	// Auth middleware
	amw := middlewares.AuthMiddleware{Secret: secret, Repo: &repo}
	r.Use(amw.Middleware)

	httpSrv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
		// Enforcing timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Server fully initialised. Connect to http://localhost:%s/ for GraphQL playground", port)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.WithField("error", err.Error()).Fatal("Sudden error, terminating server")
	}
}
