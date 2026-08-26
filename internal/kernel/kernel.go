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
	Enable() error
	Ping() error
	Restart() error
	Disable() error
	ParseSubLink(string) error
	ListSubConfig() ([]ConfigFileInfo, error)
	LoadSubConfig(string) error
	ReloadConfigGroup(string) error
	CreateGroup(string) error
	ListGroups() ([]GroupInfo, error)
	AddToGroup(string, string, bool) error
	ListModes() ([]string, error)
	SwitchMode(string) error
}
