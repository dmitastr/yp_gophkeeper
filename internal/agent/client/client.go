package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/core/responses"
)

type IClient interface {
	Authenticate(body *requests.AuthRequest, address string) (*string, error)
	Ping(bearerToken string, address string) error
	AddPassword(bearerToken, address, login, password string) error
	AddSecret(bearerToken, address string, body *requests.SecretRequest) error
	GetSecret(bearerToken, address string, secretID int) (*models.Secret, error)
	GetAllSecrets(bearerToken, address string) ([]models.SecretInfo, error)
	UpdateSecret(bearerToken, address string, body *requests.SecretRequest) error
	DeleteSecret(bearerToken, address string, secretID int) error
}

type Client struct {
	client *retryablehttp.Client
}

func NewClient() IClient {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = time.Millisecond * 300
	httpClient.RetryMax = 3
	httpClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Second * time.Duration(2*attemptNum+1)
	}
	return &Client{client: httpClient}
}

func (c *Client) Authenticate(body *requests.AuthRequest, address string) (*string, error) {
	endpoint := "login"
	if body.IsNewUser {
		endpoint = "auth"
	}
	reqBuilder, err := NewRequestBuilder(http.MethodPost, address, "/api/"+endpoint)

	if err != nil {
		return nil, err
	}

	req := reqBuilder.WithBody(body).Build()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var result responses.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result.Token, nil
}

func (c *Client) Ping(bearerToken string, address string) error {
	reqBuilder, err := NewRequestBuilder(http.MethodGet, address, "/api/ping")
	if err != nil {
		return err
	}

	req := reqBuilder.WithBearer(bearerToken).Build()
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
	return nil
}

func (c *Client) AddPassword(bearerToken, address, login, password string) error {
	body := &requests.PasswordRequest{Login: login, Password: password}
	reqBuilder, err := NewRequestBuilder(http.MethodPost, address, "/api/secrets/passwords")
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req := reqBuilder.WithBearer(bearerToken).WithBody(body).Build()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error adding password: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}

func (c *Client) AddSecret(bearerToken, address string, body *requests.SecretRequest) error {
	reqBuilder, err := NewRequestBuilder(http.MethodPost, address, "/api/secrets")
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	reqBuilder = reqBuilder.WithBearer(bearerToken).WithBody(body)

	// switch body.SecretType {
	// case models.PASSWORD, models.BANK_CARD:
	// 	reqBuilder = reqBuilder.WithBody(body)
	// default:
	// 	reqBuilder = reqBuilder.WithRawBody(body.Body)
	// }

	req := reqBuilder.Build()

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error adding secret: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetAllSecrets(bearerToken, address string) ([]models.SecretInfo, error) {
	reqBuilder, err := NewRequestBuilder(http.MethodGet, address, "/api/secrets")
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req := reqBuilder.WithBearer(bearerToken).Build()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting secrets: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var secrets responses.SecretsListResponseObject
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, fmt.Errorf("error parsing secrets: %w", err)
	}

	return secrets.Secrets, nil
}

func (c *Client) GetSecret(bearerToken, address string, secretID int) (*models.Secret, error) {
	reqBuilder, err := NewRequestBuilder(http.MethodGet, address, fmt.Sprintf("/api/secrets/%d", secretID))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req := reqBuilder.WithBearer(bearerToken).Build()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting secrets: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var secret responses.SecretResponseObject
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return nil, fmt.Errorf("error parsing secrets: %w", err)
	}

	return secret.Secret, nil
}

func (c *Client) UpdateSecret(bearerToken, address string, body *requests.SecretRequest) error {
	reqBuilder, err := NewRequestBuilder(http.MethodPut, address, fmt.Sprintf("/api/secrets/%d", body.ID))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	reqBuilder = reqBuilder.WithBearer(bearerToken)

	switch body.SecretType {
	case models.PASSWORD, models.BANK_CARD:
		reqBuilder = reqBuilder.WithBody(body)
	default:
		reqBuilder = reqBuilder.WithRawBody(body.Body)
	}

	req := reqBuilder.Build()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error updating secret: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}

func (c *Client) DeleteSecret(bearerToken, address string, secretID int) error {
	reqBuilder, err := NewRequestBuilder(http.MethodDelete, address, fmt.Sprintf("/api/secrets/%d", secretID))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req := reqBuilder.WithBearer(bearerToken).Build()
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting secret: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}
