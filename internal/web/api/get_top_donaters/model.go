package getTopDonaters

import (
	"He110/donation-report-manager/internal/business/domain/configs"
	"He110/donation-report-manager/internal/business/domain/donations"
)

type Request struct {
	ChannelId string `json:"channel_id"`
}

type Response struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Donations []donations.DonationItem `json:"donations"`
	Config    *configs.Config          `json:"config"`
}
