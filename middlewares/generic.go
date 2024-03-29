// Package middlewares implements all the key HTTP middlewares, especially auth
package middlewares

import (
	"net/http"
)

type Middleware interface {
	Middleware(http.Handler) http.Handler
}
