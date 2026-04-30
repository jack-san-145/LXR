package handlers

import (
	helper "lxr-d/internal/Helper"
)

type LXRHandler struct {
	Helper *helper.Helper
}

func NewHandler() *LXRHandler {
	return &LXRHandler{
		Helper: helper.NewHelper(),
	}
}
