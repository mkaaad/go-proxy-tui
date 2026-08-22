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

type Proxy interface {
	Start() error
	Ping() error
	Restart() error
	Stop() error
	ParseSubLink(string) error
	ReadConfigFile(string) (string, error)
	ListConfigFiles() ([]ConfigFileInfo, error)
	GetModes() ([]string, error)
	SwitchMode(string) error
}
