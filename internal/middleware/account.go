package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/security"
	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/internal/tokens"
)

const AccountCookieName = "account"

func SetAccountToken(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext, userAccount *models.UserAccount) error {
	token, err := tokens.NewAccountToken(svcCtx, userAccount)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AccountCookieName,
		Value:    token,
		Path:     "/",
		Secure:   r.URL.Scheme == "https",
		HttpOnly: true,
	})

	return nil
}

func LoadAccountContextFromCookie(svcCtx *svc.ServiceContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie(AccountCookieName)
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

			account, err := models.AccountByID(r.Context(), svcCtx.SqlxDB, id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if _, err := security.ParseJWT(token, tokens.AccountTokenPrefix+svcCtx.Config.Auth.TokenSecret); err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), ContextAccountKey, account)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AccountGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/app/choose-account") {
			next.ServeHTTP(w, r)
			return
		}

		if AccountFromContext(r) != nil {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/app/choose-account", http.StatusFound)
	})
}
