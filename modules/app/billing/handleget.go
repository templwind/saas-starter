package billing

import (
	"net/http"

	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/ui/components/apperror"

	"github.com/templwind/sass-starter/internal/utils"
)

func (c *Controller) HandleGet(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r)
	accounts, err := models.FindAllAccountsByUserID(r.Context(), c.svcCtx.SqlxDB, user.ID, 0, 0)
	if err != nil {
		utils.Render(w, r, http.StatusOK, apperror.New(
			err.Error(),
		))
	}

	utils.Render(w, r, http.StatusOK, New(
		WithConfig(c.svcCtx.Config),
		WithRequest(r),
		WithAccounts(accounts),
	))
}
