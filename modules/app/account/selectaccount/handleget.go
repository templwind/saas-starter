package selectaccount

import (
	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/utils"

	"github.com/labstack/echo/v4"
)

func (c *Controller) HandleGet(e echo.Context) error {
	user := middleware.UserFromContext(e)
	accounts, err := models.FindAllAccountsByUserID(e.Request().Context(), c.svcCtx.SqlxDB, user.ID, 0, 0)
	if err != nil {
		return err
	}

	return utils.Render(e, 200, New(
		WithConfig(c.svcCtx.Config),
		WithEcho(e),
		WithAccounts(accounts),
	))
}
