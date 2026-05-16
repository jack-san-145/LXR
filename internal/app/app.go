package app

import (
	"lxr-d/internal/handlers"
	helper "lxr-d/internal/helper"
	"lxr-d/internal/ip"
)

type App struct {
	Handler       *handlers.Handler
	Helper        *helper.Helper
	NetworkConfig *ip.NetworkConfig
}

func NewApp(newNetwork *ip.NetworkConfig) *App {

	helper := helper.NewHelper(newNetwork)

	return &App{
		Handler:       handlers.NewHandler(helper),
		Helper:        helper,
		NetworkConfig: newNetwork,
	}
}
