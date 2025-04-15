package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dermaddis/op_tournament/internal/handler/routes/auth"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	e        *echo.Echo
	log      *slog.Logger
	upgrader *websocket.Upgrader
	routers  routers
}

type routers struct {
	auth *auth.AuthRouter
}

func New(log *slog.Logger) *Handler {
	e := echo.New()
	upgrader := websocket.Upgrader{}

	handler := &Handler{
		e:        e,
		log:      log,
		upgrader: &upgrader,
		routers: routers{
			auth: auth.New(e.Group("/auth"), log),
		},
	}

	handler.registerStatic()
	handler.registerRoutes()
	handler.registerMiddleware()

	return handler
}

func (h *Handler) registerStatic() {
	h.e.Static("/static", "static")
}

func (h *Handler) Serve(address string) error {
	err := h.e.Start(address)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil

}

func (h *Handler) Stop(ctx context.Context) error {
	err := h.e.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
