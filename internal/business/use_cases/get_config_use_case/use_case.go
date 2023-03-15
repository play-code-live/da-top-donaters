package getConfigUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/user"
	"He110/donation-report-manager/internal/business/errors"
)

type userStorage interface {
	GetUser(channelId string) (*user.User, error)
}

type configStorage interface {
	GetConfig(channelID string) (*configs.Config, error)
}

type UseCase struct {
	userStorage   userStorage
	configStorage configStorage
}

func New(us userStorage, cs configStorage) *UseCase {
	return &UseCase{userStorage: us, configStorage: cs}
}

func (u UseCase) Perform(channelId string) (*configs.Config, error) {
	if _, err := u.userStorage.GetUser(channelId); err != nil {
		return nil, errors.NotAuthorizedError{}
	}

	return u.configStorage.GetConfig(channelId)
}
