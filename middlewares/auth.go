package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

type AuthMiddleware struct {
	Secret string
	Repo   *db.Repository
}

// Standard HTTP.Middleware interface compliant authorisation middleware
func (amw *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")

		// Allow unauthenticated users in
		if err != nil || c == nil {
			next.ServeHTTP(w, r)
			return
		}

		tokenstr := c.Value
		id, err := auth.ParseToken(tokenstr, amw.Secret)
		if err != nil {
			http.Error(w, "Invalid auth token", http.StatusForbidden)
			return
		}

		// Determine user type
		idSplit := strings.Split(id, ":")

		var ctx context.Context

		// Retrieve user from database
		if idSplit[0] == "s" {
			s, err := amw.Repo.GetStudentById(id)

			if err != nil {
				http.Error(w, "Malformed auth token", http.StatusForbidden)
				return
			}

			ctx = context.WithValue(r.Context(), "user", map[string]interface{}{
				"user": s,
				"type": idSplit[0],
			})

		} else if idSplit[0] == "t" {
			t, err := amw.Repo.GetTutorById(id)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Malformed auth token", http.StatusForbidden)
				return
			}

			ctx = context.WithValue(r.Context(), "user", map[string]interface{}{
				"user": t,
				"type": idSplit[0],
			})
		} else {
			http.Error(w, "Malformed auth token", http.StatusForbidden)
			return
		}

		// Continue with new context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}
