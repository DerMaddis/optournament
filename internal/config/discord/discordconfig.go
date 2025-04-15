package discord

import "golang.org/x/oauth2"

var DiscordConfig = oauth2.Config{
	ClientID:     "1361645240305455184",
	ClientSecret: "MKUrjpUgaH0q5MWcRP2gAl0q4U-WPCtc",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://discord.com/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	},
	RedirectURL: "http://localhost:3000/auth/redirect",
	Scopes:      []string{"identify"},
}
