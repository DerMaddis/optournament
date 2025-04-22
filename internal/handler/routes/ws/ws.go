package ws

import (
	"log/slog"

	"github.com/dermaddis/op_tournament/internal/websocket"
	"github.com/labstack/echo/v4"
)

type WsRouter struct {
	group     *echo.Group
	log       *slog.Logger
	wsHandler *websocket.WebsocketHandler
}

func New(group *echo.Group, wsHandler *websocket.WebsocketHandler, log *slog.Logger) *WsRouter {
	router := WsRouter{
		group:     group,
		log:       log,
		wsHandler: wsHandler,
	}

	router.registerRoutes()

	return &router
}

func (r *WsRouter) registerRoutes() {
	r.group.GET("/connect", r.connect)
}

func (r *WsRouter) connect(c echo.Context) error {
	return r.wsHandler.NewConnection(c)
}
