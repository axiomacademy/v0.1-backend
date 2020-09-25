package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/middlewares"

	"github.com/solderneer/axiom-backend/services/notifs"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	repo := db.Repository{}
	repo.InitDb()
	repo.Migrate()

	defer repo.DbPool.Close()

	// Initialising all services
	ns := notifs.NotifService{Nchans: map[string]chan *model.Notification{}, Nmutex: sync.Mutex{}}

	// Binding services to resolver
	resolver := graph.Resolver{
		Secret: "password",
		Repo:   &repo,
		Ns:     &ns,
	}

	graphSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver}))

	r := mux.NewRouter()
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphSrv)

	// Auth middleware
	amw := middlewares.AuthMiddleware{Secret: "password"}
	r.Use(amw.Middleware)

	httpSrv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
		// Enforcing timeouts
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(httpSrv.ListenAndServe())
}
