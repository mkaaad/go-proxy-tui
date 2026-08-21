package mihomo

import (
	"github.com/go-resty/resty/v2"
	"github.com/mkaaad/go-proxy-tui/internal/client"
)

type mihomo struct {
	apiClient apiClient
}

func (m *mihomo) Init(config client.Config) error {
	m.apiClient.Client = resty.New().
		SetBaseURL(config.URL)
	return nil
}
func (m *mihomo) Ping() error {
	_, err := m.apiClient.Get("/version")
	if err != nil {
		return err
	}
	return nil
}
func (m *mihomo) Start() error {
	return nil
}
func (m *mihomo) restart() error {

	return nil
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
