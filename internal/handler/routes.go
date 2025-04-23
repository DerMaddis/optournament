package handler

import (
	"errors"
	"net/http"

	"github.com/dermaddis/op_tournament/internal/errs"
	"github.com/dermaddis/op_tournament/internal/handler/customcontext"
	"github.com/dermaddis/op_tournament/internal/templates"
	"github.com/labstack/echo/v4"
)

func (h *Handler) registerRoutes() {
	h.e.GET("/", h.index)

	h.e.GET("/tournament/:id", h.tournament)

	h.e.GET("/connect", h.connect)
}

func (h *Handler) index(c echo.Context) error {
	customContext, ok := c.(*customcontext.CustomContext)
	if !ok {
		customContext = &customcontext.CustomContext{
			LoggedIn: false,
			Context:  c,
		}
	}

	return templates.Render(c, http.StatusOK, templates.Index(customContext))
}

type getTournamentRequest struct {
	Id string `param:"id"`
}

func (h *Handler) tournament(c echo.Context) error {
	customContext, ok := c.(*customcontext.CustomContext)
	if !ok {
		return c.String(http.StatusInternalServerError, "internal server error context")
	}

	var data getTournamentRequest
	if err := c.Bind(&data); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	tournament, err := h.service.Tournament(data.Id, string(customContext.DiscordUser.Id))
	if err != nil {
		var customErr errs.CustomError
		if errors.As(err, &customErr) {
			return c.String(customErr.StatusCode, customErr.Error())
		}
	}

	return templates.Render(c, http.StatusOK, templates.Tournament(customContext, tournament))
}

func (h *Handler) connect(c echo.Context) error {
	return h.ws.NewConnection(c)
}
