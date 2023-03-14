package saveConfigUseCase

import "He110/donation-report-manager/internal/business/domain/configs"

type twitchService interface {
	GetChannelId(token string) (string, error)
}

type configStorage interface {
	Save(cfg *configs.Config) (*configs.Config, error)
}

type UseCase struct {
	twitchService twitchService
	configStorage configStorage
}

type Parameters struct {
	TwitchAccessToken string
	Title             string
	DonatersCount     int
}

func (u UseCase) Perform(p Parameters) (*configs.Config, error) {
	channelId, err := u.twitchService.GetChannelId(p.TwitchAccessToken)
	if err != nil {
		return nil, err
	}

	cfg := &configs.Config{
		ChannelID:     channelId,
		Title:         p.Title,
		DonatersCount: p.DonatersCount,
	}

	return u.configStorage.Save(cfg)
}
