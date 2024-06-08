package changepassword

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/utils"
)

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	account := middleware.AccountFromContext(r)
	primaryUser, _ := models.UserByID(r.Context(), c.svcCtx.SqlxDB, account.PrimaryUserID)

	utils.Render(w, r, http.StatusOK, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
		WithAccount(account),
		WithPrimaryUser(primaryUser),
	))
}
