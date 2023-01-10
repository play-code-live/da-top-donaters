package donation_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Client struct {
	accessToken, clientId, clientSecret string
	scope                               []string
	client                              *http.Client
}

const (
	host              = "https://www.donationalerts.com"
	endpointToken     = "/oauth/token"
	endpointDonations = "/api/v1/alerts/donations"
)

func NewClient(clientId, clientSecret, accessToken string) *Client {
	return &Client{
		accessToken:  accessToken,
		clientId:     clientId,
		clientSecret: clientSecret,
		scope:        []string{"oauth-donation-index"},
		client:       &http.Client{},
	}
}

func (c *Client) GetAllDonations() ([]Donation, error) {
	donations, meta, err := c.GetDonations(1)
	if err != nil {
		return nil, err
	}

	for i := 2; i <= meta.LastPage; i++ {
		nextPageDonations, _, err := c.GetDonations(i)
		if err != nil {
			return nil, err
		}
		donations = append(donations, nextPageDonations...)
	}

	return donations, nil
}

func (c *Client) GetDonations(page int) ([]Donation, *Meta, error) {
	url := fmt.Sprintf("%s?page=%d", endpointDonations, page)
	response, err := c.performRequest("GET", url, map[string]string{})
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
	url := fmt.Sprintf("%s%s", host, endpoint)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for key, value := range data {
		_ = writer.WriteField(key, value)
	}
	_ = writer.WriteField("client_id", c.clientId)
	_ = writer.WriteField("client_secret", c.clientSecret)
	_ = writer.WriteField("redirect_uri", "http://localhost/")
	_ = writer.WriteField("scope", strings.Join(c.scope, " "))
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.accessToken)
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
