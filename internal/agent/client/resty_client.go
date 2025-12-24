package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gophkeep/internal/agent/compression"
	"gophkeep/internal/core/models"
	"gophkeep/internal/core/requests"
	"gophkeep/internal/core/responses"
	"gophkeep/internal/logger"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	client      *resty.Client
	bearerToken *string
	log         logger.ILogger
	compressor  compression.ICompressor
}

func NewRestyClient(log logger.ILogger, bearerToken string) IClient {
	httpClient := &RestyClient{
		bearerToken: &bearerToken,
		log:         log,
		compressor:  compression.NewCompressor(),
	}

	client := resty.New().
		SetTimeout(time.Millisecond * 300).
		SetRetryCount(3).
		SetRetryWaitTime(time.Second * 5).
		SetRetryMaxWaitTime(20 * time.Second).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			req.SetHeader("Authorization", "Bearer "+bearerToken)
			return nil
		})
	httpClient.client = client

	return httpClient
}

func (r RestyClient) Authenticate(ctx context.Context, body *requests.AuthRequest, address string) (*string, error) {
	endpoint := "/api/login"
	if body.IsNewUser {
		endpoint = "/api/login"
	}
	path, _ := url.JoinPath(address, endpoint)

	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status())
	}

	var result responses.AuthResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return &result.Token, nil
}

func (r RestyClient) Ping(ctx context.Context, address string) error {
	endpoint := "/api/ping"
	path, _ := url.JoinPath(address, endpoint)

	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsSuccess() {
		return fmt.Errorf("bad status: %s", resp.Status())
	}

	fmt.Println(resp.String())
	return nil
}

func (r RestyClient) AddPassword(ctx context.Context, address, login, password string) error {
	// TODO implement me
	panic("implement me")
}

func (r RestyClient) AddSecret(ctx context.Context, address string, body *requests.SecretRequest) error {
	endpoint := "/api/secrets"
	if body.ID != nil {
		endpoint += fmt.Sprintf("/%d", *body.ID)
	}
	path, _ := url.JoinPath(address, endpoint)

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error encode body: %w", err)
	}
	compressed, err := r.compressor.Compress(data)
	if err != nil {
		return fmt.Errorf("error compress body: %w", err)
	}

	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(compressed).
		Post(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status())
	}

	return nil
}

func (r RestyClient) GetSecret(ctx context.Context, address string, secretID int) (*models.Secret, error) {
	endpoint := "/api/secrets"
	path, _ := url.JoinPath(address, endpoint, fmt.Sprintf("%d", secretID))

	resp, err := r.client.R().
		SetContext(ctx).
		Get(path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status())
	}

	var secret models.Secret
	if err := json.Unmarshal(resp.Body(), &secret); err != nil {
		return nil, fmt.Errorf("error unmarshalling secret: %w", err)
	}

	return &secret, nil
}

func (r RestyClient) GetAllSecrets(ctx context.Context, address string) ([]models.SecretInfo, error) {
	endpoint := "/api/secrets"
	path, _ := url.JoinPath(address, endpoint)

	resp, err := r.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+r.getToken()).
		Get(path)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status())
	}

	var secrets responses.SecretsListResponseObject
	if err := json.Unmarshal(resp.Body(), &secrets); err != nil {
		return nil, fmt.Errorf("error unmarshalling secret: %w", err)
	}

	return secrets.Secrets, nil
}

func (r RestyClient) DeleteSecret(ctx context.Context, address string, secretID int) error {
	endpoint := "/api/secrets"
	path, _ := url.JoinPath(address, endpoint, fmt.Sprintf("%d", secretID))

	resp, err := r.client.R().
		SetContext(ctx).
		Delete(path)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status())
	}

	return nil
}

func (r RestyClient) getToken() string {
	if r.bearerToken != nil {
		return *r.bearerToken
	}
	return ""
}
