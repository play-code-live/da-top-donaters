package getConfig

import (
	"net/http"
)

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelId := r.URL.Query().Get("channelId")
		if channelId == "" {
			
		}

	}
}
