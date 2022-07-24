package middlewares

import (
	"net/http"
	"uiassignment/internal/pkg/auth"

	"github.com/gorilla/mux"
)

func AccessTokenCheckMW() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accesstoken := r.Header.Get("X-Accesstoken")
			if !auth.IsAccessTokenValid(accesstoken) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
