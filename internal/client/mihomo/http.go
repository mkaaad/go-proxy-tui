package mihomo

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type apiClient struct {
	*resty.Client
}
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[api error]: [%d]: %s", e.Code, e.Message)
}
func (c *apiClient) Get(path string) ([]byte, error) {
	errResp := &APIError{}
	resp, err := c.R().SetError(&errResp).Get(path)
	if err != nil {
		return nil, fmt.Errorf("[network error]: %w", err)
	}
	if resp.IsError() {
		return nil, errResp
	}
	return resp.Body(), nil
}
func (c *apiClient) Post(path string, body interface{}) ([]byte, error) {
	errResp := &APIError{}
	resp, err := c.R().SetBody(body).SetError(&errResp).Post(path)
	if err != nil {
		return nil, fmt.Errorf("[network error]: %w", err)
	}
	if resp.IsError() {
		return nil, errResp
	}
	return resp.Body(), nil
}
func (c *apiClient) Patch(path string, body interface{}) ([]byte, error) {
	errResp := &APIError{}
	resp, err := c.R().SetBody(body).SetError(&errResp).Patch(path)
	if err != nil {
		return nil, fmt.Errorf("[network error]: %w", err)
	}
	if resp.IsError() {
		return nil, errResp
	}
	return resp.Body(), nil
}
func (c *apiClient) Put(path string, body interface{}) ([]byte, error) {
	errResp := &APIError{}
	resp, err := c.R().SetBody(body).SetError(&errResp).Put(path)
	if err != nil {
		return nil, fmt.Errorf("[network error]: %w", err)
	}
	if resp.IsError() {
		return nil, errResp
	}
	return resp.Body(), nil
}
func (c *apiClient) Delete(path string, body interface{}) ([]byte, error) {
	errResp := &APIError{}
	resp, err := c.R().SetBody(body).SetError(&errResp).Delete(path)
	if err != nil {
		return nil, fmt.Errorf("[network error]: %w", err)
	}
	if resp.IsError() {
		return nil, errResp
	}
	return resp.Body(), nil
}
