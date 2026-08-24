package mihomo

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/publicsuffix"
)

func configDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("[Config Dir Error]: %w", err)
	}
	return filepath.Join(dir, "go-proxy-tui", "mihomo"), nil
}

func secondLevelDomain(host string) string {
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if host == "" || net.ParseIP(host) != nil {
		return host
	}
	if domain, err := publicsuffix.EffectiveTLDPlusOne(host); err == nil {
		return domain
	}
	return host
}

func uniqueConfigPath(dir, link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", fmt.Errorf("[Parse Error]: invalid subscription link: %w", err)
	}
	base := sanitizeFileNameBase(secondLevelDomain(u.Hostname()))
	if base == "" {
		return "", fmt.Errorf("[Parse Error]: subscription link has no host")
	}
	for i := 0; ; i++ {
		name := base + ".yaml"
		if i > 0 {
			name = fmt.Sprintf("%s(%d).yaml", base, i)
		}
		path := filepath.Join(dir, name)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err == nil {
			f.Close()
			return path, nil
		}
		if !os.IsExist(err) {
			return "", fmt.Errorf("[Create Config Error]: %w", err)
		}
	}
}

func sanitizeFileNameBase(base string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '_'
		}
		return r
	}, base)
}
