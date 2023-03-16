package donations

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/user"
	donationClient "He110/donation-report-manager/internal/pkg/donation-alerts-client"
	"container/heap"
	"errors"
	"fmt"
	"github.com/eko/gocache/store"
	"strings"
	"time"
)

type cache interface {
	Get(key interface{}) (interface{}, error)
	Set(key interface{}, value interface{}, options *store.Options) error
	Delete(key interface{}) error
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

func (s *DonationService) GetTopDonaters(channelId string, cfg *configs.Config, accessToken string) ([]DonationItem, error) {
	donations, err := s.getTopDonaters(channelId, cfg, accessToken)
	if err != nil {
		return nil, err
	}

	return donations[:cfg.DonatersCount], nil
}

func (s *DonationService) getTopDonaters(channelId string, cfg *configs.Config, accessToken string) ([]DonationItem, error) {
	donaters, err := s.getDonatersSums(channelId, cfg, accessToken)
	if err != nil {
		return nil, err
	}

	h := &MaxDonationHeap{}
	for _, d := range donaters {
		heap.Push(h, *d)
	}

	result := make([]DonationItem, 0, 100)
	i := 0
	for h.Len() > 0 && i < 100 {
		result = append(result, heap.Pop(h).(DonationItem))
		i++
	}

	return result, nil
}

func (s *DonationService) getDonatersSums(channelId string, cfg *configs.Config, accessToken string) (map[string]*DonationItem, error) {
	cacheKey := fmt.Sprintf(cacheKeyTpl, channelId)
	cached, err := s.cache.Get(cacheKey)
	if err == nil {
		return cached.(map[string]*DonationItem), nil
	}

	prepareName := func(v string) string {
		return strings.ToLower(strings.TrimSpace(v))
	}

	ignoredNames := make(map[string]struct{}, len(cfg.NamesToIgnore))
	for _, name := range cfg.NamesToIgnore {
		ignoredNames[prepareName(name)] = struct{}{}
	}

	donations, err := s.client.GetAllDonations(channelId, accessToken)
	if err != nil {
		return nil, err
	}

	donaters := make(map[string]*DonationItem, 100)
	for _, d := range donations {
		namePrepared := prepareName(d.Username)
		if _, skip := ignoredNames[namePrepared]; skip {
			continue
		}
		if _, exists := donaters[namePrepared]; !exists {
			donaters[namePrepared] = &DonationItem{Name: d.Username, Amount: 0}
		}

		donaters[namePrepared].Amount += d.Amount
	}

	return donaters, s.cache.Set(cacheKey, donaters, &store.Options{})
}

func (s *DonationService) ResetCache(channelId string) error {
	cacheKey := fmt.Sprintf(cacheKeyTpl, channelId)

	return s.cache.Delete(cacheKey)
}
