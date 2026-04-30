package app

import (
	helper "lxr-d/internal/Helper"
	"lxr-d/internal/handlers"
)

type App struct {
	Handler *handlers.LXRHandler
	Helper  *helper.Helper
}

func NewApp() *App {
	return &App{
		Handler: handlers.NewHandler(),
		Helper:  helper.NewHelper(),
	}
}
