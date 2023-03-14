package saveTokenUseCase

import (
	"He110/donation-report-manager/internal/business/domain/user"
	"time"
)

type twitchService interface {
	GetChannelId(token string) (string, error)
}

type userStorage interface {
	SaveUser(user *user.User) (*user.User, error)
}

type UseCase struct {
	twitchService twitchService
	userStorage   userStorage
}

type Parameters struct {
	TwitchAccessToken string
	AccessToken       string
	RefreshToken      string
	ExpiresAfter      time.Time
}

func (u UseCase) Perform(p Parameters) (*user.User, error) {
	channelId, err := u.twitchService.GetChannelId(p.TwitchAccessToken)
	if err != nil {
		return nil, err
	}

	usr := &user.User{
		ChannelID:    channelId,
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
		ExpiresAfter: p.ExpiresAfter,
	}

	return u.userStorage.SaveUser(usr)
}