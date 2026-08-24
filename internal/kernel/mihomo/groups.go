package mihomo

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
)

func (c *Client) ReloadConfigGroup(group string) error {
	if group == "" || group == "." || group == ".." || group != filepath.Base(group) {
		return fmt.Errorf("[Invalid Group Name]: %q", group)
	}
	groupsDir, err := groupDir()
	if err != nil {
		return err
	}
	groupPath := filepath.Join(groupsDir, group)
	if info, err := os.Stat(groupPath); err != nil || !info.IsDir() {
		return fmt.Errorf("[Group Not Found]: %q", group)
	}
	entries, err := os.ReadDir(groupPath)
	if err != nil {
		return fmt.Errorf("[Read Group Error]: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	if len(names) == 0 {
		return fmt.Errorf("[Group Empty]: %q", group)
	}
	var parts []string
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(groupPath, name))
		if err != nil {
			return fmt.Errorf("[Read Config Error]: %w", err)
		}
		parts = append(parts, string(data))
	}
	_, err = c.api.Put("/configs?force=true", map[string]string{"payload": strings.Join(parts, "\n")})
	return err
}

func (c *Client) CreateGroup(name string) error {
	if name == "" || name == "." || name == ".." || name != filepath.Base(name) {
		return fmt.Errorf("[Invalid Group Name]: %q", name)
	}
	groupsDir, err := groupDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(groupsDir, 0o755); err != nil {
		return fmt.Errorf("[Mkdir Groups Error]: %w", err)
	}
	for i := 0; ; i++ {
		groupName := name
		if i > 0 {
			groupName = fmt.Sprintf("%s(%d)", name, i)
		}
		if err := os.Mkdir(filepath.Join(groupsDir, groupName), 0o755); err == nil {
			return nil
		} else if !os.IsExist(err) {
			return fmt.Errorf("[Mkdir Group Error]: %w", err)
		}
	}
}

func groupDir() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "groups"), nil
}

func (c *Client) ListGroups() ([]kernel.GroupInfo, error) {
	groupsDir, err := groupDir()
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(groupsDir)
	if os.IsNotExist(err) {
		return []kernel.GroupInfo{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("[List Groups Error]: %w", err)
	}

	groups := make([]kernel.GroupInfo, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("[List Groups Error]: %w", err)
		}
		groupPath := filepath.Join(groupsDir, entry.Name())
		groupEntries, err := os.ReadDir(groupPath)
		if err != nil {
			return nil, fmt.Errorf("[List Groups Error]: %w", err)
		}
		memberCount := 0
		var names []string
		for _, groupEntry := range groupEntries {
			if !groupEntry.IsDir() {
				memberCount++
				names = append(names, groupEntry.Name())
			}
		}
		groups = append(groups, kernel.GroupInfo{
			Name:        entry.Name(),
			MemberCount: memberCount,
			ModTime:     info.ModTime(),
			Configs:     names,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups, nil
}

func (c *Client) AddToGroup(group, file string, overwrite bool) error {
	if group == "" || group == "." || group == ".." || group != filepath.Base(group) {
		return fmt.Errorf("[Invalid Group Name]: %q", group)
	}
	if file == "" || file == "." || file == ".." || file != filepath.Base(file) {
		return fmt.Errorf("[Invalid File Name]: %q", file)
	}
	groupsDir, err := groupDir()
	if err != nil {
		return err
	}
	groupPath := filepath.Join(groupsDir, group)
	if info, err := os.Stat(groupPath); err != nil || !info.IsDir() {
		return fmt.Errorf("[Group Not Found]: %q", group)
	}
	configsDir, err := configDir()
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(configsDir, file)); err != nil {
		return fmt.Errorf("[Config Not Found]: %q", file)
	}
	linkPath := filepath.Join(groupPath, file)
	if _, err := os.Lstat(linkPath); err == nil {
		if !overwrite {
			return nil
		}
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("[Remove Link Error]: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("[Stat Link Error]: %w", err)
	}
	if err := os.Symlink(filepath.Join("..", "..", file), linkPath); err != nil {
		return fmt.Errorf("[Symlink Error]: %w", err)
	}
	return nil
}
