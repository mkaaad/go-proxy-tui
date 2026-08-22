package mihomo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/mkaaad/go-proxy-tui/internal/kernel/rest"
)

const defaultBaseURL = "http://127.0.0.1:9090"
const (
	modeRule   = "rule"
	modeGlobal = "global"
	modeDirect = "direct"
)

type Client struct {
	api *rest.Client
}

func New(opts kernel.Options) (*Client, error) {
	if opts.URL == "" {
		opts.URL = defaultBaseURL
	}
	if !strings.Contains(opts.URL, "://") {
		opts.URL = "https://" + opts.URL
	}
	if _, err := url.ParseRequestURI(opts.URL); err != nil {
		return nil, fmt.Errorf("[URL wrong]: %w", err)
	}
	return &Client{api: rest.New(opts.URL, opts.Secret, 15*time.Second)}, nil
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
		return "", fmt.Errorf("[pase error]: %w", err)
	}
	return resp.Version, nil
}

func (c *Client) GetModes() ([]string, error) {
	return []string{modeRule, modeGlobal, modeDirect}, nil
}

func (c *Client) SwitchMode(mode string) error {
	_, err := c.api.Patch("/configs", map[string]string{"mode": mode})
	return err
}

func (c *Client) Restart() error {
	_, err := c.api.Post("/restart", nil)
	return err
}

func (c *Client) Stop() error {
	return nil
}

func (c *Client) ParseSubLink(link string) error {
	body, err := rest.New("", "", 15*time.Second).Get(link)
	if err != nil {
		return err
	}
	payload, err := parseLink2payload(body)
	if err != nil {
		return err
	}

	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("[mkdir error]: %w", err)
	}
	path, err := uniqueConfigPath(dir, link)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		return fmt.Errorf("[write config error]: %w", err)
	}
	return nil
}

func (c *Client) ReadConfigFile(name string) (string, error) {
	if name == "" || name != filepath.Base(name) {
		return "", fmt.Errorf("[invalid name]: %q", name)
	}
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", fmt.Errorf("[read config error]: %w", err)
	}
	return string(data), nil
}

func (c *Client) ListConfigFiles() ([]kernel.ConfigFileInfo, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []kernel.ConfigFileInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[list config error]: %w", err)
	}

	files := make([]kernel.ConfigFileInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("[list config error]: %w", err)
		}
		files = append(files, kernel.ConfigFileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if !files[i].ModTime.Equal(files[j].ModTime) {
			return files[i].ModTime.After(files[j].ModTime)
		}
		return files[i].Name < files[j].Name
	})
	return files, nil
}

func configDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("[config dir error]: %w", err)
	}
	return filepath.Join(dir, "go-proxy-tui", "mihomo"), nil
}
