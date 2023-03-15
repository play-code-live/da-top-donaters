package saveConfigUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/user"
	"He110/donation-report-manager/internal/business/errors"
)

type configStorage interface {
	Save(cfg *configs.Config) (*configs.Config, error)
}

type userStorage interface {
	GetUser(channelId string) (*user.User, error)
}

type UseCase struct {
	userStorage   userStorage
	configStorage configStorage
}

func New(us userStorage, cs configStorage) *UseCase {
	return &UseCase{userStorage: us, configStorage: cs}
}

type Parameters struct {
	ChannelId     string
	Title         string
	DonatersCount int
}

func (u UseCase) Perform(p Parameters) (*configs.Config, error) {
	_, err := u.userStorage.GetUser(p.ChannelId)
	if err != nil {
		return nil, errors.NotAuthorizedError{}
	}

	cfg := &configs.Config{
		ChannelID:     p.ChannelId,
		Title:         p.Title,
		DonatersCount: p.DonatersCount,
	}

	return u.configStorage.Save(cfg)
}
