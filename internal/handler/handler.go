package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dermaddis/op_tournament/internal/handler/routes/auth"
	"github.com/dermaddis/op_tournament/internal/handler/routes/ws"
	"github.com/dermaddis/op_tournament/internal/service"
	"github.com/dermaddis/op_tournament/internal/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	e       *echo.Echo
	service *service.Service
	ws      *websocket.WebsocketHandler
	log     *slog.Logger
	routers routers
}

type routers struct {
	auth *auth.AuthRouter
	ws   *ws.WsRouter
}

func New(service *service.Service, wsHandler *websocket.WebsocketHandler, log *slog.Logger) *Handler {
	e := echo.New()
	handler := &Handler{
		e:       e,
		service: service,
		ws:      wsHandler,
		log:     log,
		routers: routers{
			auth: auth.New(e.Group("/auth"), log),
			ws:   ws.New(e.Group("/ws"), wsHandler, log),
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
