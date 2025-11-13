package types

import "time"

type UserSessionData struct {
	Email        string
	Name         string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
