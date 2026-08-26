package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	*resty.Client
}

type APIError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("[Api Error]: %s", http.StatusText(e.Status))
	}
	return fmt.Sprintf("[Api Error]: [%d]: %s", e.Status, e.Message)
}

func New(baseURL, secret string, timeout time.Duration) *Client {
	rc := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(timeout)
	if secret != "" {
		rc.SetAuthToken(secret)
	}
	return &Client{Client: rc}
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.do(resty.MethodGet, path, nil)
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.do(resty.MethodPost, path, body)
}

func (c *Client) Put(path string, body any) ([]byte, error) {
	return c.do(resty.MethodPut, path, body)
}

func (c *Client) Patch(path string, body any) ([]byte, error) {
	return c.do(resty.MethodPatch, path, body)
}

func (c *Client) Delete(path string, body any) ([]byte, error) {
	return c.do(resty.MethodDelete, path, body)
}

func (c *Client) do(method, path string, body any) ([]byte, error) {
	resp, err := c.R().SetBody(body).Execute(method, path)
	if err != nil {
		return nil, fmt.Errorf("[Network Error]: %w", err)
	}
	if resp.IsError() {
		apiErr := &APIError{Status: resp.StatusCode(), Message: string(resp.Body())}
		return nil, apiErr
	}
	return resp.Body(), nil
}
