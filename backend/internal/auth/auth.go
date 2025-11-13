package auth

type AuthMethod string

const (
	AuthMethodNone   AuthMethod = "none"
	AuthMethodOAuth2 AuthMethod = "oauth2"
	AuthMethodPlain  AuthMethod = "plain"
)
