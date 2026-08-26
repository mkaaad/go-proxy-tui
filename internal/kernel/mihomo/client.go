package mihomo

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/mkaaad/go-proxy-tui/internal/kernel/rest"
	"gopkg.in/yaml.v3"
)

const defaultBaseURL = "http://127.0.0.1:9090"

type Client struct {
	api *rest.Client
}

func kernelConfigDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "kernel"), nil
}

func createKernelConfig(name, url, secret string) error {
	dir, err := kernelConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("[Mkdir Error]: %w", err)
	}
	path, err := uniqueNamePath(dir, name)
	if err != nil {
		return err
	}
	kernelConfig := &kernel.KernelConfig{
		URL:    url,
		Secret: secret,
	}
	data, _ := yaml.Marshal(kernelConfig)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("[Write Config Error]: %w", err)
	}
	return nil
}

func (c *Client) NewConfig(opts kernel.Options) error {
	if opts.URL == "" {
		opts.URL = defaultBaseURL
	}
	if !strings.Contains(opts.URL, "://") {
		opts.URL = "https://" + opts.URL
	}
	if _, err := url.ParseRequestURI(opts.URL); err != nil {
		return fmt.Errorf("[URL Wrong]: %w", err)
	}
	err := createKernelConfig(opts.Name, opts.URL, opts.Secret)
	if err != nil {
		return err
	}
	c.api = rest.New(opts.URL, opts.Secret, 15*time.Second)
	return nil
}

func (c *Client) LoadConfig(name string) error {
	if name == "" || name != filepath.Base(name) {
		return fmt.Errorf("[Invalid Name]: %q", name)
	}
	dir, err := kernelConfigDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("[Read Config Error]: %w", err)
	}
	var cfg kernel.KernelConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("[Parse Config Error]: %w", err)
	}
	c.api = rest.New(cfg.URL, cfg.Secret, 15*time.Second)
	return nil
}
func (c *Client) ListConfig() ([]kernel.ConfigFileInfo, error) {
	dir, err := kernelConfigDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []kernel.ConfigFileInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[List Config Error]: %w", err)
	}

	files := make([]kernel.ConfigFileInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("[List Config Error]: %w", err)
		}
		files = append(files, kernel.ConfigFileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
	return files, nil
}
