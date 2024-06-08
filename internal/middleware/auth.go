package middleware

import (
	"context"
	"net/http"

	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/security"
	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/internal/tokens"
)

const AuthCookieName = "auth"

func SetAuthToken(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext, user *models.User) error {
	token, err := tokens.NewUserAuthToken(svcCtx, user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})

	return nil
}

func LoadAuthContextFromCookie(svcCtx *svc.ServiceContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie(AuthCookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token := tokenCookie.Value

			unverifiedClaims, err := security.ParseUnverifiedJWT(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			id, _ := unverifiedClaims["id"].(string)

			user, err := models.UserByID(r.Context(), svcCtx.SqlxDB, id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if _, err := security.ParseJWT(token, user.TokenKey+svcCtx.Config.Auth.TokenSecret); err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := UserFromContext(r); user != nil {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	})
}
