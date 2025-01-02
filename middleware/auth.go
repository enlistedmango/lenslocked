package middleware

import (
	"net/http"

	"github.com/enlistedmango/lenslocked/models"
	"github.com/gorilla/sessions"
)

const (
	SessionName = "lenslocked-session"
	UserIDKey   = "userID"
)

type AuthMiddleware struct {
	Store       *sessions.CookieStore
	UserService *models.UserService
}

func (amw *AuthMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := amw.Store.Get(r, SessionName)
		userID, ok := session.Values[UserIDKey].(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		user, err := amw.UserService.GetByID(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = models.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
