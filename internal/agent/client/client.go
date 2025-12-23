package client

import (
	"context"
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
	Authenticate(ctx context.Context, body *requests.AuthRequest, address string) (*string, error)
	Ping(ctx context.Context, bearerToken string, address string) error
	AddPassword(ctx context.Context, bearerToken, address, login, password string) error
	AddSecret(ctx context.Context, bearerToken, address string, body *requests.SecretRequest) error
	GetSecret(ctx context.Context, bearerToken, address string, secretID int) (*models.Secret, error)
	GetAllSecrets(ctx context.Context, bearerToken, address string) ([]models.SecretInfo, error)
	DeleteSecret(ctx context.Context, bearerToken, address string, secretID int) error
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

func (c *Client) Authenticate(ctx context.Context, body *requests.AuthRequest, address string) (*string, error) {
	endpoint := "/api/login"
	if body.IsNewUser {
		endpoint = "/api/auth"
	}
	reqBuilder, err := NewRequestBuilder(ctx, http.MethodPost, address, endpoint)

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

func (c *Client) Ping(ctx context.Context, bearerToken string, address string) error {
	reqBuilder, err := NewRequestBuilder(ctx, http.MethodPost, address, "/api/ping")
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

func (c *Client) AddPassword(ctx context.Context, bearerToken, address, login, password string) error {
	body := &requests.PasswordRequest{Login: login, Password: password}
	reqBuilder, err := NewRequestBuilder(ctx, http.MethodPost, address, "/api/passwords")
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

func (c *Client) AddSecret(ctx context.Context, bearerToken, address string, body *requests.SecretRequest) error {
	endpointPath := "/api/secrets"
	if body.ID != nil {
		endpointPath += fmt.Sprintf("/%d", *body.ID)
	}

	reqBuilder, err := NewRequestBuilder(ctx, http.MethodPost, address, endpointPath)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	reqBuilder = reqBuilder.WithBearer(bearerToken).WithBody(body)

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

func (c *Client) GetAllSecrets(ctx context.Context, bearerToken, address string) ([]models.SecretInfo, error) {
	endpointPath := "/api/secrets"

	reqBuilder, err := NewRequestBuilder(ctx, http.MethodGet, address, endpointPath)
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

func (c *Client) GetSecret(ctx context.Context, bearerToken, address string, secretID int) (*models.Secret, error) {
	endpointPath := fmt.Sprintf("/api/secrets/%d", secretID)

	reqBuilder, err := NewRequestBuilder(ctx, http.MethodGet, address, endpointPath)
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

func (c *Client) DeleteSecret(ctx context.Context, bearerToken, address string, secretID int) error {
	endpointPath := fmt.Sprintf("/api/secrets/%d", secretID)

	reqBuilder, err := NewRequestBuilder(ctx, http.MethodDelete, address, endpointPath)
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
