package donation_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Client struct {
	clientId, clientSecret string
	redirectHost           string
	scope                  []string
	client                 *http.Client
}

const (
	host              = "https://www.donationalerts.com"
	endpointToken     = "/oauth/token"
	endpointDonations = "/api/v1/alerts/donations"
	endpointAuthorize = "/oauth/authorize"

	redirectEndpoint = "/donationAlertsRedirectUri/"
)

func NewClient(clientId, clientSecret, redirectHost string) (*Client, error) {
	c := &Client{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectHost: redirectHost,
		scope:        []string{"oauth-donation-index"},
		client:       &http.Client{},
	}

	return c, nil
}

func (c *Client) getRedirectUri(channelId string) string {
	return c.redirectHost + fmt.Sprintf(redirectEndpoint+"%s", channelId)
}

func (c *Client) GetAuthLink(channelId string) string {
	//query := url.Values{}
	//query.Add("response_type", "code")
	//query.Add("client_id", c.clientId)
	//query.Add("redirect_uri", c.getRedirectUri(channelId))
	//query.Add("scope", strings.Join(c.scope, " "))

	query := fmt.Sprintf(
		"response_type=%s&client_id=%s&redirect_uri=%s&scope=%s",
		"code",
		c.clientId,
		c.getRedirectUri(channelId),
		strings.Join(c.scope, " "),
	)
	return host + endpointAuthorize + "?" + query
}

func (c *Client) ListenAndServerAuthHandler(address string, onSuccess func(channelID string, response TokenResponse) error) error {
	var server http.Server
	router := mux.NewRouter()
	router.HandleFunc(redirectEndpoint+"{channelID}", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		channelID := mux.Vars(r)["channelID"]
		tokenResponse, err := c.ObtainAccessToken(code, channelID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//TODO Respond with some DTO
			return
		}

		if err = onSuccess(channelID, *tokenResponse); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//TODO Respond with some DTO
			return
		}
		//TODO Respond with some DTO
	})

	server = http.Server{
		Addr:    address,
		Handler: router,
	}
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (c *Client) ObtainAccessToken(code, channelID string) (*TokenResponse, error) {
	response, err := c.performRequest("POST", endpointToken, map[string]string{
		"code":       code,
		"grant_type": "authorization_code",
		"channel_id": channelID,
	})
	if err != nil {
		return nil, err
	}
	var result TokenResponse
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) RefreshToken(refreshToken, channelID string) (*TokenResponse, error) {
	response, err := c.performRequest("POST", endpointToken, map[string]string{
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
		"channel_id":    channelID,
	})
	if err != nil {
		return nil, err
	}

	var result TokenResponse
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetAllDonations(channelID, accessToken string) ([]Donation, error) {
	donations, meta, err := c.GetDonations(1, channelID, accessToken)
	if err != nil {
		return nil, err
	}

	for i := 2; i <= meta.LastPage; i++ {
		nextPageDonations, _, err := c.GetDonations(i, channelID, accessToken)
		if err != nil {
			return nil, err
		}
		donations = append(donations, nextPageDonations...)
	}

	return donations, nil
}

func (c *Client) GetDonations(page int, channelId, accessToken string) ([]Donation, *Meta, error) {
	requestUrl := fmt.Sprintf("%s?page=%d", endpointDonations, page)
	response, err := c.performRequest("GET", requestUrl, map[string]string{
		"access_token": accessToken,
		"client_id":    channelId,
	})
	if err != nil {
		return nil, nil, err
	}

	result := DonationResponse{}
	if err = json.Unmarshal(response, &result); err != nil {
		return nil, nil, err
	}

	return result.Data, &result.Meta, nil
}

func (c *Client) performRequest(method, endpoint string, data map[string]string) ([]byte, error) {
	requestUrl := host + endpoint
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for key, value := range data {
		_ = writer.WriteField(key, value)
	}
	_ = writer.WriteField("client_id", c.clientId)
	_ = writer.WriteField("client_secret", c.clientSecret)
	_ = writer.WriteField("scope", strings.Join(c.scope, " "))
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, requestUrl, payload)
	if err != nil {
		return nil, err
	}

	if channelId, exists := data["channel_id"]; exists {
		_ = writer.WriteField("redirect_uri", c.getRedirectUri(channelId))
	}

	if accessToken, exists := data["access_token"]; exists {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
