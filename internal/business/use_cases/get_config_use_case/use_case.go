package getConfigUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/user"
)

type twitchService interface {
	GetChannelId(token string) (string, error)
}

type userStorage interface {
	GetUser(channelId string) (*user.User, error)
}

type configStorage interface {
	GetConfig(channelID string) (*configs.Config, error)
}

type UseCase struct {
	twitchService twitchService
	userStorage   userStorage
	configStorage configStorage
}

func (u UseCase) Perform(twitchAccessToken string) (*configs.Config, error) {
	channelId, err := u.twitchService.GetChannelId(twitchAccessToken)
	if err != nil {
		return nil, err
	}

	if _, err = u.userStorage.GetUser(channelId); err != nil {
		return nil, err
	}
	
	return u.configStorage.GetConfig(channelId)
}
