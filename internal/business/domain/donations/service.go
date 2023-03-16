package donations

import (
	"He110/donation-report-manager/internal/business/domain/user"
	donationClient "He110/donation-report-manager/internal/pkg/donation-alerts-client"
	"container/heap"
	"errors"
	"fmt"
	"github.com/eko/gocache/store"
	"time"
)

type cache interface {
	Get(key interface{}) (interface{}, error)
	Set(key interface{}, value interface{}, options *store.Options) error
}

const cacheKeyTpl = "top_donaters_%s"

type DonationService struct {
	client *donationClient.Client
	cache  cache
}

func NewService(client *donationClient.Client, cacheManager cache) *DonationService {
	return &DonationService{client: client, cache: cacheManager}
}

func (s *DonationService) EnsureTokenRefreshed(usr *user.User) (*user.User, error) {
	if usr.RefreshToken == "" {
		return nil, errors.New("cannot refresh token without required refresh-token")
	}

	if usr.ExpiresAfter.After(time.Now()) {
		return usr, nil
	}

	tokenData, err := s.client.RefreshToken(usr.RefreshToken)
	if err != nil {
		return nil, err
	}

	usr.RefreshToken = tokenData.RefreshToken
	usr.AccessToken = tokenData.AccessToken
	usr.ExpiresAfter = time.Unix(int64(tokenData.ExpiresIn), 0)

	return usr, nil
}

func (s *DonationService) GetTopDonaters(channelId string, count int, accessToken string) ([]DonationItem, error) {
	donations, err := s.getTopDonaters(channelId, count, accessToken)
	if err != nil {
		return nil, err
	}

	return donations[:count], nil
}

func (s *DonationService) getTopDonaters(channelId string, count int, accessToken string) ([]DonationItem, error) {
	donaters, err := s.getDonatersSums(channelId, accessToken)
	if err != nil {
		return nil, err
	}

	h := &MaxDonationHeap{}
	for name, sum := range donaters {
		heap.Push(h, DonationItem{
			Name:   name,
			Amount: sum,
		})
	}

	result := make([]DonationItem, 0, 100)
	i := 0
	for h.Len() > 0 && i < 100 {
		result = append(result, heap.Pop(h).(DonationItem))
		i++
	}

	return result, nil
}

func (s *DonationService) getDonatersSums(channelId, accessToken string) (map[string]float64, error) {
	cacheKey := fmt.Sprintf(cacheKeyTpl, channelId)
	cached, err := s.cache.Get(cacheKey)
	if err == nil {
		return cached.(map[string]float64), nil
	}

	donations, err := s.client.GetAllDonations(channelId, accessToken)
	if err != nil {
		return nil, err
	}

	donaters := make(map[string]float64, 100)
	for _, d := range donations {
		if _, exists := donaters[d.Username]; !exists {
			donaters[d.Username] = 0
		}

		donaters[d.Username] += d.Amount
	}

	return donaters, s.cache.Set(cacheKey, donaters, &store.Options{})
}
