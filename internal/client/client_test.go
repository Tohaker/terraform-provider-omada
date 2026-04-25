package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_Success(t *testing.T) {
	var receivedGrant, receivedClientId string
	accessToken := "mock-access-token"
	controllerId := "mock-controller-id"
	clientId := "mock-client-id"
	clientSecret := "mock-client-secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedGrant = r.URL.Query().Get("grant_type")

		var body struct {
			ClientId string `json:"client_id"`
		}

		_ = json.NewDecoder(r.Body).Decode(&body)
		receivedClientId = body.ClientId

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": 0,
			"result":    map[string]any{"accessToken": accessToken},
		})
	}))

	// Close the server after the test runs
	t.Cleanup(server.Close)

	meta, err := New(context.Background(), Config{
		Host:         server.URL,
		ControllerID: controllerId,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})

	if err != nil {
		t.Fatal(err)
	}

	if meta.OmadacID != controllerId {
		t.Errorf("OmadaID value: %q", meta.OmadacID)
	}

	if receivedGrant != "client_credentials" {
		t.Errorf("Grant type: %q", receivedGrant)
	}

	if receivedClientId != clientId {
		t.Errorf("Client ID: %q", receivedClientId)
	}

	receivedToken := meta.Client.GetConfig().DefaultHeader["Authorization"]

	if receivedToken != fmt.Sprintf("AccessToken=%s", accessToken) {
		t.Errorf("Recieved Token: %q", receivedToken)
	}
}

func TestNew_AuthFailure(t *testing.T) {
	controllerId := "mock-controller-id"
	clientId := "mock-client-id"
	clientSecret := "mock-client-secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": -1001,
			"msg":       "Could not get access token",
		})
	}))

	// Close the server after the test runs
	t.Cleanup(server.Close)

	_, err := New(context.Background(), Config{
		Host:         server.URL,
		ControllerID: controllerId,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})

	if err == nil {
		t.Fatal("No error returned")
	}

	if err.Error() != "400 Bad Request" {
		t.Errorf("Received Error: %s", err.Error())
	}
}

func TestNew_NoResult(t *testing.T) {
	controllerId := "mock-controller-id"
	clientId := "mock-client-id"
	clientSecret := "mock-client-secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": 0,
			"msg":       "No result",
		})
	}))

	// Close the server after the test runs
	t.Cleanup(server.Close)

	_, err := New(context.Background(), Config{
		Host:         server.URL,
		ControllerID: controllerId,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})

	if err == nil {
		t.Fatal("No error returned")
	}

	if err.Error() != "token response missing result" {
		t.Errorf("Received Error: %s", err.Error())
	}
}

func TestNew_NoAccessToken(t *testing.T) {
	controllerId := "mock-controller-id"
	clientId := "mock-client-id"
	clientSecret := "mock-client-secret"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"errorCode": 0,
			"result":    map[string]any{},
		})
	}))

	// Close the server after the test runs
	t.Cleanup(server.Close)

	_, err := New(context.Background(), Config{
		Host:         server.URL,
		ControllerID: controllerId,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})

	if err == nil {
		t.Fatal("No error returned")
	}

	if err.Error() != "token response missing access token" {
		t.Errorf("Received Error: %s", err.Error())
	}
}
