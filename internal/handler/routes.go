package handler

import (
	"net/http"

	"github.com/dermaddis/op_tournament/internal/handler/customcontext"
	"github.com/dermaddis/op_tournament/internal/templates"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (h *Handler) registerRoutes() {
	h.e.GET("/", h.index)

	h.e.GET("/auth/login", h.login)
	h.e.GET("/auth/redirect", h.redirect)

	h.e.GET("/connect", h.connect)
}

func (h *Handler) index(c echo.Context) error {
	customContext, ok := c.(*customcontext.CustomContext)
	if !ok {
		return c.String(http.StatusInternalServerError, "internal server error context")
	}

	return templates.Render(c, http.StatusOK, templates.Index(customContext))
}

func (h *Handler) login(c echo.Context) error {
	url := discordConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

type getApiAuthRedirectRequest struct {
	Code string `query:"code"`
}

func (h *Handler) redirect(c echo.Context) error {
	var data getApiAuthRedirectRequest
	if err := c.Bind(&data); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	token, err := discordConfig.Exchange(c.Request().Context(), data.Code)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry.Unix(),
	})

	jwtString, err := jwt.SignedString([]byte("THE_ULT1M4TE_S3CRET"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    jwtString,
		Path:     "/",
		HttpOnly: true,
	})

	return c.Redirect(http.StatusFound, "/")
}

func (h *Handler) connect(c echo.Context) error {
	ws, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	defer ws.Close()

	return nil
}
