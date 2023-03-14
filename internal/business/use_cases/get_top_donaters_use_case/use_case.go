package getTopDonatersUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/donations"
	"He110/donation-report-manager/internal/business/domain/user"
)

type daService interface {
	EnsureTokenRefreshed(usr *user.User) (string, error)
	GetTopDonaters(count int, accessToken string) ([]donations.DonationItem, error)
}

type userStorage interface {
	GetUser(channelId string) (*user.User, error)
}

type configStorage interface {
	GetConfig(channelID string) (*configs.Config, error)
}

type UseCase struct {
	userStorage   userStorage
	configStorage configStorage
	daService     daService
}

func (u UseCase) Perform(channelId string) ([]donations.DonationItem, error) {
	usr, err := u.userStorage.GetUser(channelId)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.daService.EnsureTokenRefreshed(usr)
	if err != nil {
		return nil, err
	}

	cfg, err := u.configStorage.GetConfig(usr.ChannelID)
	if err != nil {
		return nil, err
	}

	return u.daService.GetTopDonaters(cfg.DonatersCount, accessToken)
}
