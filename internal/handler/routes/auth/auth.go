package auth

import (
	"log/slog"
	"net/http"

	"github.com/dermaddis/op_tournament/internal/config/discord"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type AuthRouter struct {
	group *echo.Group
	log   *slog.Logger
}

func New(group *echo.Group, log *slog.Logger) *AuthRouter {
	router := AuthRouter{
		group: group,
		log:   log,
	}

	router.registerRoutes()

	return &router
}

func (r *AuthRouter) registerRoutes() {
	r.group.GET("/login", r.login)
	r.group.GET("/redirect", r.redirect)
}

func (r *AuthRouter) login(c echo.Context) error {
	url := discord.DiscordConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusFound, url)
}

type getApiAuthRedirectRequest struct {
	Code string `query:"code"`
}

func (r *AuthRouter) redirect(c echo.Context) error {
	var data getApiAuthRedirectRequest
	if err := c.Bind(&data); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	token, err := discord.DiscordConfig.Exchange(c.Request().Context(), data.Code)
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
