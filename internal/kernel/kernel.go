package kernel

import "time"

type Options struct {
	Name   string
	URL    string
	Secret string
}

type ConfigFileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}
type KernelConfig struct {
	Name   string
	URL    string
	Secret string
}
type GroupInfo struct {
	Name        string
	MemberCount int
	ModTime     time.Time
	Configs     []string
}

type Proxy interface {
	NewConfig(Options) error
	LoadConfig(string) error
	ListConfig() ([]ConfigFileInfo, error)
	Start() error
	Ping() error
	Restart() error
	Stop() error
	ParseSubLink(string) error
	ListSubConfig() ([]ConfigFileInfo, error)
	LoadSubConfig(string) error
	ReloadConfigGroup(string) error
	CreateGroup(string) error
	ListGroups() ([]GroupInfo, error)
	AddToGroup(string, string, bool) error
	GetModes() ([]string, error)
	SwitchMode(string) error
}
