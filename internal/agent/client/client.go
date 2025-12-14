package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type RequestBuilder struct {
	req *retryablehttp.Request
}

func NewRequestBuilder(method, baseURL, suffix string) (*RequestBuilder, error) {
	path, err := url.JoinPath(baseURL, suffix)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(method, path, nil)
	if err != nil {
		return nil, err
	}
	return &RequestBuilder{req: req}, nil
}

func (r *RequestBuilder) WithBody(body interface{}) *RequestBuilder {
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	r.req.Body = io.NopCloser(bytes.NewReader(data))
	r.req.ContentLength = int64(len(data))
	r.req.Header.Set("Content-Type", "application/json")

	r.req.SetBody(data)

	return r
}

func (r *RequestBuilder) WithBearer(token string) *RequestBuilder {
	r.req.Header.Add("Authorization", "Bearer "+token)
	return r
}

func (r *RequestBuilder) Build() *retryablehttp.Request {
	return r.req
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type Client struct {
	client *retryablehttp.Client
}

func NewClient() *Client {
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = time.Millisecond * 300
	httpClient.RetryMax = 3
	httpClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Second * time.Duration(2*attemptNum+1)
	}
	return &Client{client: httpClient}
}

func (c *Client) Authenticate(body *AuthRequest, address string) (*string, error) {
	reqBuilder, err := NewRequestBuilder(http.MethodPost, address, "/api/auth")

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

	var result AuthResponse
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
	body := &PasswordRequest{Login: login, Password: password}
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
