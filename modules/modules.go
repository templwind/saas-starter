package modules

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/svc"
)

// Module interface that all modules must implement
type Module interface {
	Register(svcCtx *svc.ServiceContext, mux *http.ServeMux) error
}

func RegisterAll(svcCtx *svc.ServiceContext, mux *http.ServeMux) error {
	for _, m := range registry {
		err := m.Register(svcCtx, mux)
		if err != nil {
			return err
		}
	}
	return nil
}
