package handler

import (
	"net/http"

	"github.com/dermaddis/op_tournament/internal/handler/customcontext"
	"github.com/dermaddis/op_tournament/internal/templates"
	"github.com/labstack/echo/v4"
)

func (h *Handler) registerRoutes() {
	h.e.GET("/", h.index)

	h.e.GET("/connect", h.connect)
}

func (h *Handler) index(c echo.Context) error {
	customContext, ok := c.(*customcontext.CustomContext)
	if !ok {
		return c.String(http.StatusInternalServerError, "internal server error context")
	}

	return templates.Render(c, http.StatusOK, templates.Index(customContext))
}

func (h *Handler) connect(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	defer ws.Close()

	return nil
}
