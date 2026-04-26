package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Tohaker/omada-go-sdk/omada"
)

type Config struct {
	Host, ControllerID, ClientID, ClientSecret string
	HTTPClient                                 *http.Client
}

type Meta struct {
	Client   *omada.APIClient
	OmadacId string
}

func New(ctx context.Context, cfg Config) (*Meta, error) {
	// Create a new Omada client using the configuration values
	config := omada.NewConfiguration()
	config.Servers = omada.ServerConfigurations{
		{URL: cfg.Host},
	}
	config.HTTPClient = cfg.HTTPClient

	client := omada.NewAPIClient(config)

	tokenResp, _, err := client.AuthorizeAPI.AuthorizeToken(ctx).GrantType("client_credentials").TokenRequest(omada.TokenRequest{
		ClientId:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		OmadacId:     &cfg.ControllerID,
	}).Execute()

	if err != nil {
		return nil, err
	}

	result, ok := tokenResp.GetResultOk()
	if !ok {
		return nil, fmt.Errorf("token response missing result")
	}

	accessToken, ok := result.GetAccessTokenOk()
	if !ok {
		return nil, fmt.Errorf("token response missing access token")
	}

	config.DefaultHeader["Authorization"] = "AccessToken=" + *accessToken

	return &Meta{
		Client:   client,
		OmadacId: cfg.ControllerID,
	}, nil
}
