package www

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/modules/www/index"
	"github.com/templwind/sass-starter/modules/www/login"
	"github.com/templwind/sass-starter/modules/www/register"
)

func Module() *WwwModule {
	return &WwwModule{}
}

type WwwModule struct {
	Name string
}

func (m *WwwModule) Register(svcCtx *svc.ServiceContext, mux *http.ServeMux) error {
	m.Name = "www"

	// Home
	mux.Handle("GET /", http.HandlerFunc(index.NewController(svcCtx).HandleGet))

	// Login
	mux.Handle("GET /login", http.HandlerFunc(login.NewController(svcCtx).HandleGet))
	mux.Handle("POST /login", http.HandlerFunc(login.NewController(svcCtx).HandlePost))

	// Register
	mux.Handle("GET /register", http.HandlerFunc(register.NewController(svcCtx).HandleGet))
	mux.Handle("POST /register", http.HandlerFunc(register.NewController(svcCtx).HandlePost))

	return nil
}
