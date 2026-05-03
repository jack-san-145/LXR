package app

import (
	"lxr-d/internal/handlers"
	helper "lxr-d/internal/helper"
)

type App struct {
	Handler *handlers.Handler
	Helper  *helper.Helper
}

func NewApp() *App {
	helper := helper.NewHelper()

	return &App{
		Handler: handlers.NewHandler(helper),
		Helper:  helper,
	}
}
