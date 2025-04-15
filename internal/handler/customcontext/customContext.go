package customcontext

import (
	"github.com/dermaddis/op_tournament/internal/model/discord"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	DiscordUser discord.APIUser
	LoggedIn    bool
	echo.Context
}
