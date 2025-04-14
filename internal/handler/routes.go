package handler

import (
	"net/http"

	"github.com/dermaddis/op_tournament/internal/templates"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func (h *Handler) registerRoutes() {
	h.e.GET("/", func(c echo.Context) error {
		return templates.Render(c, http.StatusOK, templates.Index())
	})

	h.e.GET("/connect", h.connect)
}

func (h *Handler) connect(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("<h1 id=\"header\">Connected</h1>"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return nil
}
