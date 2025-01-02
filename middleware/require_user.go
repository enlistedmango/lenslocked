package middleware

import (
	"net/http"

	"github.com/enlistedmango/lenslocked/models"
)

type RequireUser struct {
	User *models.User
}

func (mw *RequireUser) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := models.UserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
