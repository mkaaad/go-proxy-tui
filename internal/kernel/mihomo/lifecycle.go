package mihomo

import (
	"encoding/json"
	"fmt"
)

const (
	modeRule   = "rule"
	modeGlobal = "global"
	modeDirect = "direct"
)

type tun struct {
	Enable    bool   `json:"enable"`
	Stack     string `json:"stack"`
	Device    string `json:"device"`
	AutoRoute bool   `json:"auto-route"`
}

func (c *Client) Ping() error {
	_, err := c.api.Get("/")
	return err
}

func (c *Client) Version() (string, error) {
	data, err := c.api.Get("/version")
	if err != nil {
		return "", err
	}
	var resp struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("[Parse Error]: %w", err)
	}
	return resp.Version, nil
}

func (c *Client) ListModes() ([]string, error) {
	return []string{modeRule, modeGlobal, modeDirect}, nil
}

func (c *Client) SwitchMode(mode string) error {
	_, err := c.api.Patch("/configs", map[string]string{"mode": mode})
	return err
}

func (c *Client) Enable() error {
	t := tun{
		Enable:    true,
		Stack:     "mixed",
		Device:    "mihomo",
		AutoRoute: true,
	}
	_, err := c.api.Patch("/configs", map[string]tun{
		"tun": t,
	})
	return err
}

func (c *Client) Restart() error {
	_, err := c.api.Post("/restart", nil)
	return err
}

func (c *Client) Disable() error {
	t := tun{
		Enable: false,
	}
	_, err := c.api.Patch("/configs", map[string]tun{
		"tun": t,
	})
	return err
}
