package app

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/modules/app/account"
	"github.com/templwind/sass-starter/modules/app/account/changepassword"
	"github.com/templwind/sass-starter/modules/app/account/selectaccount"
	"github.com/templwind/sass-starter/modules/app/billing"
	"github.com/templwind/sass-starter/modules/app/dashboard"
	"github.com/templwind/sass-starter/modules/app/settings"

	"github.com/templwind/templwind/htmx"
)

func Module() *AppModule {
	return &AppModule{}
}

type AppModule struct {
	Name string
}

func (m *AppModule) Register(svcCtx *svc.ServiceContext, mux *http.ServeMux) error {
	m.Name = "app"

	// Middleware for authentication and account loading
	authMiddleware := middleware.Chain(
		middleware.LoadAuthContextFromCookie(svcCtx),
		middleware.AuthGuard,
		middleware.LoadAccountContextFromCookie(svcCtx),
		middleware.AccountGuard,
	)

	// logout
	mux.Handle("GET /app/logout", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.ClearCookies(w,
			middleware.AuthCookieName,    // clear auth cookie
			middleware.AccountCookieName, // clear account cookie
		)
		htmx.Redirect(w, r, "/login")
	})))

	// dashboard
	mux.Handle("GET /app/dashboard", authMiddleware(dashboard.NewController(svcCtx).HandleGet))
	// g.GET("/dashboard", dashboard.NewController(svcCtx).HandleGet)

	// settings
	mux.Handle("GET /app/settings", authMiddleware(settings.NewController(svcCtx).HandleGet))
	mux.Handle("GET /app/settings/change-password", authMiddleware(changepassword.NewController(svcCtx).HandleGet))

	// account
	mux.Handle("GET /app/account", authMiddleware(account.NewController(svcCtx).HandleGet))

	// billing
	mux.Handle("GET /app/account/billing", authMiddleware(billing.NewController(svcCtx).HandleGet))

	// select-account
	mux.Handle("GET /app/account/select-account", authMiddleware(selectaccount.NewController(svcCtx).HandleGet))
	mux.Handle("POST /app/account/select-account/:id", authMiddleware(selectaccount.NewController(svcCtx).HandlePost))

	return nil
}
