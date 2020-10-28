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

type EnvVar struct {
	Value    string
	Required bool
}

// Returns all the required env variables
func getEnvVariables() map[string]EnvVar {
	envars := map[string]EnvVar{
		"PORT":                           EnvVar{Value: defaultPort, Required: false},
		"DB_URL":                         EnvVar{Value: defaultDbUrl, Required: false},
		"SERVER_SECRET":                  EnvVar{Value: defaultSecret, Required: false},
		"GOOGLE_APPLICATION_CREDENTIALS": EnvVar{Value: "", Required: true},
	}

	for name, envar := range envars {
		raw := os.Getenv(name)

		if raw == "" {
			// Must have it set, no default
			if envar.Required == true {
				log.WithFields(log.Fields{
					"Name": name,
				}).Fatal("Cannot find required environmental variable")
			} else {
				log.WithFields(log.Fields{
					"Name":    name,
					"Default": envar.Value,
				}).Warn("This environmental variable is not found, using default")
			}
		}

		envar.Value = raw
	}

	return envars
}

func main() {
	// Setup logger
	var logger = log.New()

	// Get default environment variables
	envars := getEnvVariables()

	repo := db.Repository{}
	repo.Init(logger, envars["DB_URL"].Value)
	repo.Migrate(envars["DB_URL"].Value)

	defer repo.Close()

	// Initialising all services
	ns := notifs.NotifService{}
	ns.Init(logger)

	ms := match.MatchService{}
	ms.Init(logger, &ns, &repo)

	// Binding services to resolver
	resolver := graph.Resolver{
		Secret: envars["SERVER_SECRET"].Value,
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
	amw := middlewares.AuthMiddleware{Secret: envars["SERVER_SECRET"].Value, Repo: &repo}
	r.Use(amw.Middleware)

	httpSrv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + envars["PORT"].Value,
		// Enforcing timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Server fully initialised. Connect to http://localhost:%s/ for GraphQL playground", envars["PORT"].Value)
	if err := httpSrv.ListenAndServe(); err != nil {
		log.WithField("error", err.Error()).Fatal("Sudden error, terminating server")
	}
}
