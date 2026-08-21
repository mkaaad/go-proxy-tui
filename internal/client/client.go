package client

type Client interface {
	Init(Config) error
	Start() error
	Restart() error
	End() error
	GetModes() ([]string, error)
	SwitchMode(string) error
}
type Config struct {
	URL string
}
