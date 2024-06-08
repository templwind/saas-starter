package register

import (
	"github.com/templwind/sass-starter/internal/svc"
)

type Controller struct {
	svcCtx *svc.ServiceContext
	form   *RegisterForm
}

func NewController(svcCtx *svc.ServiceContext) *Controller {
	return &Controller{
		svcCtx: svcCtx,
		form:   new(RegisterForm),
	}
}
