package web

import "He110/donation-report-manager/internal/business/domain/configs"

type GetConfigData struct {
	IsAuthorized  bool
	ChannelId     string
	SocketHost    string
	Config        *configs.Config
	NamesToIgnore string
}

type RedirectData struct {
	ClientID     string
	ClientSecret string
}
