package middlewares

import (
	"context"
	"net/http"
	"uiassignment/internal/pkg/auth"

	"github.com/gorilla/mux"
)

func AccessTokenCheckMW() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accesstoken := r.Header.Get("X-Accesstoken")
			isTokenValid, tokenOwner := auth.IsAccessTokenValid(accesstoken)
			if !isTokenValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "tokenOwner", tokenOwner)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func OwnerAccessCheckMW() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accesstoken := r.Header.Get("X-Accesstoken")
			isTokenValid, tokenOwner := auth.IsAccessTokenValid(accesstoken)
			if !isTokenValid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(r)
			account := vars["account"]
			if account != tokenOwner {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "tokenOwner", tokenOwner)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}
