package app

import (
	"lxr-d/internal/handlers"
	helper "lxr-d/internal/helper"
	"lxr-d/internal/ip"
)

type App struct {
	Handler *handlers.Handler
	Helper  *helper.Helper
	IpStack *ip.IpStack
}

func NewApp(ipStack *ip.IpStack) *App {

	helper := helper.NewHelper(ipStack)

	return &App{
		Handler: handlers.NewHandler(helper),
		Helper:  helper,
		IpStack: ipStack,
	}
}
