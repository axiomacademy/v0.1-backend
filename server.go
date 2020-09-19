package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/middlewares"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db.InitDb()
	db.Migrate()

	defer db.DbPool.Close()

	graphSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Secret: "password"}}))

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
