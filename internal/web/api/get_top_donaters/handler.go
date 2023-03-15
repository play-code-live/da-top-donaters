package getTopDonaters

import (
	getTopDonatersUseCase "He110/donation-report-manager/internal/business/use_cases/get_top_donaters_use_case"
	"encoding/json"
	"log"
	"net/http"
)

func NewTopDonatersHandler(uc *getTopDonatersUseCase.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		request := &Request{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Println(err.Error())
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		donationsResult, err := uc.Perform(request.ChannelId)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Cannot get top donaters", http.StatusInternalServerError)
			return
		}

		response := &Response{
			Data: ResponseData{
				Donations: donationsResult.Donaters,
				Config:    donationsResult.Config,
			},
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err.Error())
			http.Error(w, "Cannot get top donaters", http.StatusInternalServerError)
			return
		}
	}
}
