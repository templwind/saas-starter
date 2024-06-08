package account

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/ui/components/apperror"
	"github.com/templwind/sass-starter/internal/utils"
)

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	account := middleware.AccountFromContext(r)
	primaryUser, err := models.UserByID(r.Context(), c.svcCtx.SqlxDB, account.PrimaryUserID)

	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
	}

	utils.Render(w, r, http.StatusOK, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
		WithAccount(account),
		WithPrimaryUser(primaryUser),
	))
}
