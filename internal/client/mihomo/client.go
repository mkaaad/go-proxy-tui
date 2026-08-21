package mihomo

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/mkaaad/go-proxy-tui/internal/client"
)

type mihomo struct {
	apiClient apiClient
}

func (m *mihomo) Init(config client.Config) error {
	if config.URL != "" {
		m.apiClient.Client = resty.New().
			SetBaseURL(config.URL)
		return nil
	}
	m.apiClient.Client = resty.New().
		SetBaseURL("http://127.0.0.1:9090")
	return nil
}

func (m *mihomo) Ping() error {
	_, err := m.apiClient.Get("/")
	if err != nil {
		return err
	}
	return nil
}

func (m *mihomo) Start() error {
	return nil
}

func (m *mihomo) Version() (string, error) {
	resp := map[string]string{}
	respData, err := m.apiClient.Get("/version")
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(respData, &resp); err != nil {
		return "", err
	}
	return resp["version"], err
}

func (m *mihomo) Restart() error {
	_, err := m.apiClient.Post("/restart", nil)
	return err
}

func (m *mihomo) Stop() error {
	return nil
}

func (m *mihomo) GetModes() ([]string, error) {
	return []string{"Rule", "Global", "Direct"}, nil
}

func (m *mihomo) SwitchModes(mode string) error {
	_, err := m.apiClient.Patch("/configs", map[string]string{
		"mode": mode,
	})
	if err != nil {
		return err
	}
	return nil
}
