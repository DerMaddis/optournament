package discord

// https://discord.com/developers/docs/resources/user#user-object-user-flags
type UserFlags uint64

// https://discord.com/developers/docs/resources/user#user-object-premium-types
type UserPremiumType uint16

type APIUser struct {
	Id    Snowflake `json:"id"`
	IdInt int
	// The user's username, not unique across the platform
	Username string `json:"username"`
	// The user's 4-digit discord-tag
	Discriminator string `json:"discriminator"`
	// The user's avatar hash
	// See https://discord.com/developers/docs/reference#image-formatting
	Avatar string `json:"avatar"`
	// Whether the user belongs to an OAuth2 application
	Bot bool `json:"bot"`
	// Whether the user is an Official Discord System user (part of the urgent message system)
	System bool `json:"system"`
	// Whether the user has two factor enabled on their account
	MfaEnabled bool `json:"mfa_enabled"`
	// The user's banner hash
	// See https://discord.com/developers/docs/reference#image-formatting
	Banner string `json:"banner"`
	// The user's banner color encoded as an integer representation of hexadecimal color code
	AccentColor int `json:"accent_color"`
	// Whether the email on this account has been verified
	Verified bool `json:"verified"`
	// The user's email
	Email string `json:"email"`
	// The flags on a user's account
	// See https://discord.com/developers/docs/resources/user#user-object-user-flags
	Flags UserFlags `json:"flags"`
	// The type of Nitro subscription on a user's account
	// See https://discord.com/developers/docs/resources/user#user-object-premium-types
	PremiumType UserPremiumType `json:"premium_type"`
	// The public flags on a user's account
	// See https://discord.com/developers/docs/resources/user#user-object-user-flags
	PublicFlags UserFlags `json:"public_flags"`
}
