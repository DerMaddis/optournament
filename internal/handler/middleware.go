package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	discordCfg "github.com/dermaddis/op_tournament/internal/config/discord"
	"github.com/dermaddis/op_tournament/internal/handler/customcontext"
	"github.com/dermaddis/op_tournament/internal/model/discord"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func (h *Handler) registerMiddleware() {
	h.e.Use(h.userMiddleware)
}

func (h *Handler) userMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/auth") {
			return next(c)
		}

		jwtCookie, err := c.Cookie("jwt")
		if err != nil {
			customContext := &customcontext.CustomContext{
				DiscordUser: discord.APIUser{},
				LoggedIn:    false,
				Context:     c,
			}
			return next(customContext)
		}

		token, err := jwt.Parse(jwtCookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte("THE_ULT1M4TE_S3CRET"), nil
		})
		if err != nil {
			return c.String(http.StatusUnauthorized, "unauthorized")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.String(http.StatusUnauthorized, "unauthorized")
		}

		expiryInt := int64(claims["expiry"].(float64))
		expiry := time.Unix(expiryInt, 0)

		oauthToken := oauth2.Token{
			AccessToken:  claims["access_token"].(string),
			RefreshToken: claims["refresh_token"].(string),
			Expiry:       expiry,
		}

		tokenSource := discordCfg.DiscordConfig.TokenSource(c.Request().Context(), &oauthToken)

		client := oauth2.NewClient(c.Request().Context(), tokenSource)

		resp, err := client.Get("https://discord.com/api/v10/users/@me")
		if err != nil {
			h.log.Error("auth/middleware: failed to get user", slog.Any("err", err))
			return c.String(http.StatusInternalServerError, "internal server error")
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusTooManyRequests {
				return c.String(http.StatusTooManyRequests, "rate limit exceeded")
			}

			body, _ := io.ReadAll(resp.Body)
			h.log.Error("auth/middleware: failed to get user", slog.Any("body", string(body)))
			return c.String(http.StatusInternalServerError, "internal server error")
		}

		var discordUser discord.APIUser
		err = json.NewDecoder(resp.Body).Decode(&discordUser)
		if err != nil {
			h.log.Error("auth/middleware: failed to decode user", slog.Any("err", err))
			return c.String(http.StatusInternalServerError, "internal server error")
		}

		idInt, err := strconv.Atoi(string(discordUser.Id))
		if err != nil {
			body, _ := io.ReadAll(resp.Body)
			h.log.Debug("body", slog.String("body", string(body)))
			h.log.Error("auth/middleware: failed to parse user id", slog.Any("err", err), slog.Any("discordUser", discordUser))
			return c.String(http.StatusInternalServerError, "internal server error")
		}
		discordUser.IdInt = idInt

		customContext := &customcontext.CustomContext{
			DiscordUser: discordUser,
			LoggedIn:    true,
			Context:     c,
		}

		return next(customContext)

	}
}
