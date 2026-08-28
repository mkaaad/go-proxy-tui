package mihomo

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
)

func (c *Client) ListProxies() ([]kernel.ProxyInfo, error) {
	data, err := c.api.Get("/proxies")
	if err != nil {
		return nil, err
	}
	var resp struct {
		Proxies map[string]struct {
			Type  string   `json:"type"`
			Now   string   `json:"now"`
			All   []string `json:"all"`
			Alive bool     `json:"alive"`
			Hist  []struct {
				Delay int `json:"delay"`
			} `json:"history"`
		} `json:"proxies"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("[Parse Error]: %w", err)
	}
	proxies := make([]kernel.ProxyInfo, 0, len(resp.Proxies))
	for name, p := range resp.Proxies {
		delay := 0
		if n := len(p.Hist); n > 0 {
			delay = p.Hist[n-1].Delay
		}
		proxies = append(proxies, kernel.ProxyInfo{
			Name:  name,
			Type:  p.Type,
			Now:   p.Now,
			All:   p.All,
			Delay: delay,
			Alive: p.Alive,
		})
	}
	sort.Slice(proxies, func(i, j int) bool {
		return proxies[i].Name < proxies[j].Name
	})
	return proxies, nil
}
