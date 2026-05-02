package handlers

import (
	helper "lxr-d/internal/helper"
)

type Handler struct {
	Helper *helper.Helper
}

func NewHandler() *Handler {
	return &Handler{
		Helper: helper.NewHelper(),
	}
}
