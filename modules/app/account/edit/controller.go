package edit

import (
	"github.com/templwind/sass-starter/internal/middleware"
	"github.com/templwind/sass-starter/internal/models"
	"github.com/templwind/sass-starter/internal/svc"
	"github.com/templwind/sass-starter/internal/utils"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	svcCtx *svc.ServiceContext
}

func NewController(svcCtx *svc.ServiceContext) *Controller {
	return &Controller{
		svcCtx: svcCtx,
	}
}

func (c *Controller) Index(e echo.Context) error {
	account := middleware.AccountFromContext(ctx)
	primaryUser, _ := models.UserByID(ctx.Request().Context(), c.svcCtx.SqlxDB, account.PrimaryUserID)

	// return utils.Render(ctx, 200, views.Index(ctx, c.svcCtx.Config, account, primaryUser))
	return nil
}
