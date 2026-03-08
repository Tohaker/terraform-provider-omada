package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	Host         string
	HTTPClient   *http.Client
	CustomerId   string
	ClientId     string
	ClientSecret string
}

type Client struct {
	accessToken string
	baseURL     *url.URL
	httpClient  *http.Client
}

type AuthorizationResult struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	RefreshToken string `json:"refreshToken"`
}

type AuthorizationResponse struct {
	ErrorCode int                  `json:"errorCode"`
	Message   string               `json:"msg"`
	Result    *AuthorizationResult `json:"result"`
}

func GetAccessToken(cfg ClientConfig) (string, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/openapi/authorize/token", cfg.Host))
	if err != nil {
		return "", err
	}

	q := endpoint.Query()
	q.Add("grant_type", "client_credentials")

	endpoint.RawQuery = q.Encode()

	body := fmt.Appendf(nil, `{
			"omadacId": "%s",
			"client_id": "%s",
			"client_secret": "%s"
		}`, cfg.CustomerId, cfg.ClientId, cfg.ClientSecret)

	res, err := cfg.HTTPClient.Post(endpoint.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	response := &AuthorizationResponse{}
	derr := json.NewDecoder(res.Body).Decode(response)
	if derr != nil {
		return "", derr
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Authorization returned status %d and error code %d: %s", res.StatusCode, response.ErrorCode, response.Message)
	}

	if response.Result == nil {
		return "", fmt.Errorf("No Result was returned from the API, something else may be wrong")
	}

	return response.Result.AccessToken, nil
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("Missing Host")
	}

	openapiPath, err := url.Parse("openapi")
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(cfg.Host)
	if err != nil {
		return nil, err
	}

	// TODO: Remove this when we have figured out certificates in the test environment
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var httpClient *http.Client
	if cfg.HTTPClient != nil {
		httpClient = cfg.HTTPClient
		httpClient.Transport = tr
	} else {
		httpClient = &http.Client{
			Transport: tr,
		}
		cfg.HTTPClient = httpClient
	}

	accessToken, err := GetAccessToken(cfg)
	if err != nil {
		return nil, err
	}

	client := &Client{
		accessToken: accessToken,
		baseURL:     url.ResolveReference(openapiPath),
		httpClient:  httpClient,
	}

	return client, nil
}
