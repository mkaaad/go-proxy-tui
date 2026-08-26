package mihomo

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/mkaaad/go-proxy-tui/internal/kernel/rest"
)

func subscriptionDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "subscription"), nil
}

func (c *Client) ParseSubLink(link string) error {
	body, err := rest.New("", "", 15*time.Second).Get(link)
	if err != nil {
		return err
	}
	payload, err := parseLinksToPayload(body)
	if err != nil {
		return err
	}

	dir, err := subscriptionDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("[Mkdir Error]: %w", err)
	}
	path, err := uniqueSubPath(dir, link)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(payload), 0o644); err != nil {
		return fmt.Errorf("[Write Config Error]: %w", err)
	}
	return nil
}

func (c *Client) LoadSubConfig(name string) error {
	if name == "" || name != filepath.Base(name) {
		return fmt.Errorf("[Invalid Name]: %q", name)
	}
	dir, err := subscriptionDir()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("[Read Config Error]: %w", err)
	}
	_, err = c.api.Patch("/configs?force=true", map[string]string{"payload": string(data)})
	return err
}

func (c *Client) ListSubConfig() ([]kernel.ConfigFileInfo, error) {
	dir, err := subscriptionDir()
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
