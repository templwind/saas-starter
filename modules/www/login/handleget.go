package login

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/utils"
)

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	err := utils.Render(w, r, http.StatusOK, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
		WithForm(c.form),
	))
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
