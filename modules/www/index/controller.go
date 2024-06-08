package index

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/internal/utils"
)

type Controller struct {
	svcCtx *svc.ServiceContext
}

func NewController(svcCtx *svc.ServiceContext) *Controller {
	return &Controller{
		svcCtx: svcCtx,
	}
}

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	if err := utils.Render(w, r, 200, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
	)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
