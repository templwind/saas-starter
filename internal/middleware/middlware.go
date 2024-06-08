package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/templwind/sass-starter/internal/models"
)

const (
	ContextAccountKey string = "accountCtx"
	ContextUserKey    string = "userCtx"
)

func AccountFromContext(r *http.Request) *models.Account {
	account, ok := r.Context().Value(ContextAccountKey).(*models.Account)
	if !ok {
		return nil
	}
	return account
}

func UserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(ContextUserKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func ClearCookies(w http.ResponseWriter, cookieNames ...string) {
	for _, cookieName := range cookieNames {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   -1,
		})
	}
}

// Middleware to remove trailing slash and redirect
func RemoveTrailingSlash(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Middleware to inject account and user into request context
func InjectContextData(next http.Handler, account *models.Account, user *models.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextAccountKey, account)
		ctx = context.WithValue(ctx, ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Chain(middlewares ...func(http.Handler) http.Handler) func(http.HandlerFunc) http.HandlerFunc {
	return func(final http.HandlerFunc) http.HandlerFunc {
		h := http.Handler(final)
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h.ServeHTTP
	}
}
