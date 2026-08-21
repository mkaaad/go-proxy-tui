// Package kernel 定义内核适配器的公共接口与连接选项。
package kernel

// Options 描述如何连接一个内核的控制接口。
type Options struct {
	URL    string
	Secret string
}

type Proxy interface {
	Start() error
	Ping() error
	Restart() error
	Stop() error
	ParseSubLink(string) error
	GetModes() ([]string, error)
	SwitchMode(string) error
}
