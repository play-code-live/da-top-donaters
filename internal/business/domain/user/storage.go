package user

import "errors"

type Storage interface {
	GetUser(channelId string) (*User, error)
	SaveUser(user *User) (*User, error)
}

type InMemoryStorage struct {
	memory map[string]User
}

func NewStorage() Storage {
	return &InMemoryStorage{
		memory: make(map[string]User, 0),
	}
}

func (s *InMemoryStorage) GetUser(channelId string) (*User, error) {
	if user, exists := s.memory[channelId]; exists {
		return &user, nil
	}

	return nil, errors.New("user not found")
}

func (s *InMemoryStorage) SaveUser(user *User) (*User, error) {
	if user == nil || user.ChannelID == "" {
		return nil, errors.New("wrong user data")
	}

	s.memory[user.ChannelID] = *user

	return user, nil
}
