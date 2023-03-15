package getTopDonatersUseCase

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/donations"
)

type TopDonatersResult struct {
	Config   *configs.Config          `json:"config"`
	Donaters []donations.DonationItem `json:"donaters"`
}
