package donation_client

type Donation struct {
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type DonationResponse struct {
	Data []Donation `json:"data"`
	Meta Meta       `json:"meta"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
