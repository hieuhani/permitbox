package handler

import "log/slog"

type BaseHandler struct {
	Logger *slog.Logger
}

func NewBaseHandler(logger *slog.Logger) BaseHandler {
	return BaseHandler{
		Logger: logger,
	}
}
