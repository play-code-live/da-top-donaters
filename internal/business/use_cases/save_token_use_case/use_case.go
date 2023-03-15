package saveTokenUseCase

import (
	"He110/donation-report-manager/internal/business/domain/user"
	"time"
)

type userStorage interface {
	SaveUser(user *user.User) (*user.User, error)
}

type UseCase struct {
	userStorage userStorage
}

func New(us userStorage) *UseCase {
	return &UseCase{userStorage: us}
}

type Parameters struct {
	ChannelId    string
	AccessToken  string
	RefreshToken string
	ExpiresAfter time.Time
}

func (u UseCase) Perform(p Parameters) (*user.User, error) {
	usr := &user.User{
		ChannelID:    p.ChannelId,
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
		ExpiresAfter: p.ExpiresAfter,
	}

	return u.userStorage.SaveUser(usr)
}
