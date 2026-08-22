package kernel

import "time"

type Options struct {
	URL    string
	Secret string
}

type ConfigFileInfo struct {
	Name    string
	Size    int64
	ModTime time.Time
}

type GroupInfo struct {
	Name        string
	MemberCount int
	ModTime     time.Time
}

type Proxy interface {
	Start() error
	Ping() error
	Restart() error
	Stop() error
	ParseSubLink(string) error
	ListConfigFiles() ([]ConfigFileInfo, error)
	ReloadConfigFile(string) error
	ReloadConfigGroup(string) error
	CreateGroup(string) error
	ListGroups() ([]GroupInfo, error)
	AddToGroup(string, string, bool) error
	GetModes() ([]string, error)
	SwitchMode(string) error
}
