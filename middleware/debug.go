package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func Debug(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n=== Debug Request ===\n")
		fmt.Printf("Method: %s\n", r.Method)
		fmt.Printf("Path: %s\n", r.URL.Path)
		fmt.Printf("Handler: %T\n", next)

		rctx := chi.RouteContext(r.Context())
		if rctx == nil {
			fmt.Printf("Route Context: nil\n")
		} else {
			fmt.Printf("Route Pattern: %s\n", rctx.RoutePattern())
			fmt.Printf("URL Params: %+v\n", rctx.URLParams)
		}
		fmt.Printf("===================\n")

		next.ServeHTTP(w, r)
	})
}
