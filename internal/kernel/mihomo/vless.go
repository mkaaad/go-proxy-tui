package mihomo

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type proxyConfig struct {
	Name           string       `yaml:"name"`
	Type           string       `yaml:"type"`
	Server         string       `yaml:"server"`
	Port           int          `yaml:"port"`
	UUID           string       `yaml:"uuid"`
	UDP            bool         `yaml:"udp"`
	Network        string       `yaml:"network,omitempty"`
	Flow           string       `yaml:"flow,omitempty"`
	TLS            bool         `yaml:"tls,omitempty"`
	Servername     string       `yaml:"servername,omitempty"`
	Fingerprint    string       `yaml:"fingerprint,omitempty"`
	SkipCertVerify bool         `yaml:"skip-cert-verify,omitempty"`
	WSOpts         *wsOpts      `yaml:"ws-opts,omitempty"`
	GRPCOpts       *grpcOpts    `yaml:"grpc-opts,omitempty"`
	RealityOpts    *realityOpts `yaml:"reality-opts,omitempty"`
}

type wsOpts struct {
	Path    string            `yaml:"path,omitempty"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type grpcOpts struct {
	ServiceName string `yaml:"grpc-service-name"`
}

type realityOpts struct {
	PublicKey string `yaml:"public-key"`
	ShortID   string `yaml:"short-id"`
	SpiderX   string `yaml:"spider-x"`
}

type clashConfig struct {
	MixedPort   int            `yaml:"mixed-port"`
	Mode        string         `yaml:"mode"`
	LogLevel    string         `yaml:"log-level"`
	Proxies     []*proxyConfig `yaml:"proxies"`
	ProxyGroups []proxyGroup   `yaml:"proxy-groups"`
	Rules       []string       `yaml:"rules"`
}

type proxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
}

func buildConfigPayload(links []string) (string, error) {
	cfg := clashConfig{
		MixedPort: 7890,
		Mode:      modeRule,
		LogLevel:  "info",
	}
	group := proxyGroup{Name: "PROXY", Type: "select"}
	for _, link := range links {
		proxy, err := parseVlessLink(link)
		if err != nil {
			continue
		}
		cfg.Proxies = append(cfg.Proxies, proxy)
		group.Proxies = append(group.Proxies, proxy.Name)
	}
	if len(cfg.Proxies) == 0 {
		return "", fmt.Errorf("[Parse Error]: no valid vless links found")
	}
	group.Proxies = append(group.Proxies, "DIRECT")
	cfg.ProxyGroups = []proxyGroup{group}
	cfg.Rules = []string{"MATCH,PROXY"}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return "", fmt.Errorf("[Marshal Error]: %w", err)
	}
	return string(data), nil
}

func parseVlessLink(link string) (*proxyConfig, error) {
	u, err := url.Parse(link)
	if err != nil || u.Scheme != "vless" || u.User == nil || u.Hostname() == "" {
		return nil, fmt.Errorf("[Parse Error]: unsupported link: %s", link)
	}

	port := 443
	if p := u.Port(); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			port = n
		}
	}

	q := u.Query()
	network := q.Get("type")
	if network == "" {
		network = "tcp"
	}

	name := u.Fragment
	if name == "" {
		name = u.Hostname()
	}

	proxy := &proxyConfig{
		Name:    name,
		Type:    "vless",
		Server:  u.Hostname(),
		Port:    port,
		UUID:    u.User.Username(),
		UDP:     true,
		Network: network,
	}

	switch security := q.Get("security"); security {
	case "tls", "reality":
		proxy.TLS = true
		if sni := q.Get("sni"); sni != "" {
			proxy.Servername = sni
		} else if host := q.Get("host"); host != "" && network == "ws" {
			proxy.Servername = host
		}
		if fp := q.Get("fp"); fp != "" {
			proxy.Fingerprint = fp
		}
		if security == "reality" {
			opts := &realityOpts{
				PublicKey: q.Get("pbk"),
				ShortID:   q.Get("sid"),
				SpiderX:   q.Get("spx"),
			}
			if opts.PublicKey != "" || opts.ShortID != "" || opts.SpiderX != "" {
				proxy.RealityOpts = opts
			}
		}
	}

	if insecure := q.Get("insecure"); insecure == "1" || strings.EqualFold(insecure, "true") {
		proxy.SkipCertVerify = true
	}
	if flow := q.Get("flow"); flow != "" {
		proxy.Flow = flow
	}

	switch network {
	case "ws":
		opts := &wsOpts{Path: q.Get("path")}
		if host := q.Get("host"); host != "" {
			opts.Headers = map[string]string{"Host": host}
		}
		if opts.Path != "" || opts.Headers != nil {
			proxy.WSOpts = opts
		}
	case "grpc":
		if svc := q.Get("serviceName"); svc != "" {
			proxy.GRPCOpts = &grpcOpts{ServiceName: svc}
		}
	}

	return proxy, nil
}

func parseLinksToPayload(data []byte) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(string(data))
		if err != nil {
			return "", fmt.Errorf("[Parse Error]: invalid base64 subscription: %w", err)
		}
	}

	var rawLinks []string
	for line := range strings.SplitSeq(string(decoded), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			rawLinks = append(rawLinks, line)
		}
	}
	if len(rawLinks) == 0 {
		return "", fmt.Errorf("[Parse Error]: subscription contains no links")
	}

	payload, err := buildConfigPayload(rawLinks)
	if err != nil {
		return "", err
	}
	return payload, nil
}
