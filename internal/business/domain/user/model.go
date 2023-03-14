package user

import "time"

type User struct {
	ChannelID    string    `json:"channel_id"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAfter time.Time `json:"-"`
}
