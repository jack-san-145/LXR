package handlers

import (
	helper "lxr-d/internal/helper"
)

type Handler struct {
	Helper *helper.Helper
}

func NewHandler(helper *helper.Helper) *Handler {
	return &Handler{
		Helper: helper,
	}
}
