package handlers

import (
	helper "lxr-d/internal/Helper"
)

type Handler struct {
	Helper *helper.Helper
}

func NewHandler() *Handler {
	return &Handler{
		Helper: helper.NewHelper(),
	}
}
