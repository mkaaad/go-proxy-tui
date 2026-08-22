package mihomo

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
)

func (c *Client) ReloadConfigFile(name string) error {
	if name == "" || name != filepath.Base(name) {
		return fmt.Errorf("[invalid name]: %q", name)
	}
	dir, err := configDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("[read config error]: %w", err)
	}
	_, err = c.api.Put("/configs?force=true", map[string]string{"payload": string(data)})
	return err
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
