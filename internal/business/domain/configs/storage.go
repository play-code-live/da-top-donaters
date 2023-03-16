package configs

import "errors"

type Storage interface {
	GetConfig(channelID string) (*Config, error)
	Save(cfg *Config) (*Config, error)
}

type InMemoryStorage struct {
	memory map[string]Config
}

func NewStorage() Storage {
	return &InMemoryStorage{
		memory: make(map[string]Config),
	}
}

func (s *InMemoryStorage) GetConfig(channelID string) (*Config, error) {
	cfg, exists := s.memory[channelID]
	if !exists {
		return s.createDefaultConfig(channelID), nil
	}

	return &cfg, nil
}

func (s *InMemoryStorage) createDefaultConfig(channelId string) *Config {
	return &Config{
		ChannelID:     channelId,
		Title:         "Топ донатеров",
		DonatersCount: 10,
		NamesToIgnore: make([]string, 0),
	}
}

func (s *InMemoryStorage) Save(cfg *Config) (*Config, error) {
	if cfg == nil || cfg.ChannelID == "" {
		return nil, errors.New("channel id must be specified")
	}

	s.memory[cfg.ChannelID] = *cfg

	return cfg, nil
}
