package settings

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/utils"
)

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	utils.Render(w, r, http.StatusOK, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
		// WithSubMenuItems(c.svcCtx.),
	))
}
