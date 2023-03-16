package getTopDonatersUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/donations"
	"He110/donation-report-manager/internal/business/domain/user"
)

type daService interface {
	EnsureTokenRefreshed(usr *user.User) (*user.User, error)
	GetTopDonaters(channelId string, cfg *configs.Config, accessToken string) ([]donations.DonationItem, error)
}

type userStorage interface {
	GetUser(channelId string) (*user.User, error)
	SaveUser(user *user.User) (*user.User, error)
}

type configStorage interface {
	GetConfig(channelID string) (*configs.Config, error)
}

type UseCase struct {
	userStorage   userStorage
	configStorage configStorage
	daService     daService
}

func New(us userStorage, cs configStorage, daService daService) *UseCase {
	return &UseCase{userStorage: us, configStorage: cs, daService: daService}
}

func (u UseCase) Perform(channelId string) (*TopDonatersResult, error) {
	usr, err := u.userStorage.GetUser(channelId)
	if err != nil {
		return nil, err
	}

	usr, err = u.daService.EnsureTokenRefreshed(usr)
	if err != nil {
		return nil, err
	}
	usr, err = u.userStorage.SaveUser(usr)
	if err != nil {
		return nil, err
	}

	cfg, err := u.configStorage.GetConfig(usr.ChannelID)
	if err != nil {
		return nil, err
	}

	topDonaters, err := u.daService.GetTopDonaters(usr.ChannelID, cfg, usr.AccessToken)
	if err != nil {
		return nil, err
	}
	return &TopDonatersResult{
		Config:   cfg,
		Donaters: topDonaters,
	}, nil
}
