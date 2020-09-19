package middlewares

import (
	"net/http"
)

type Middleware interface {
	Middleware(http.Handler) http.Handler
}
