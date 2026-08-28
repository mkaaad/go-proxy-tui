package kernel

import "errors"

var ErrNotConfigured = errors.New("[Not Configured]: load a kernel config first")

type guardedProxy struct {
	Proxy
}

func Guard(p Proxy) Proxy { return &guardedProxy{Proxy: p} }

func (g *guardedProxy) require() error {
	if !g.Ready() {
		return ErrNotConfigured
	}
	return nil
}

func (g *guardedProxy) Enable() error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.Enable()
}

func (g *guardedProxy) Ping() error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.Ping()
}

func (g *guardedProxy) Restart() error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.Restart()
}

func (g *guardedProxy) Disable() error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.Disable()
}

func (g *guardedProxy) SwitchMode(mode string) error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.SwitchMode(mode)
}

func (g *guardedProxy) ReloadConfigGroup(group string) error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.ReloadConfigGroup(group)
}

func (g *guardedProxy) LoadSubConfig(name string) error {
	if err := g.require(); err != nil {
		return err
	}
	return g.Proxy.LoadSubConfig(name)
}

func (g *guardedProxy) ListProxies() ([]ProxyInfo, error) {
	if err := g.require(); err != nil {
		return nil, err
	}
	return g.Proxy.ListProxies()
}
