package tui

import (
	"sync/atomic"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

type KernelStatus int

const (
	StatusUnknown KernelStatus = iota
	StatusRunning
	StatusStopped
)

func Status2String(code KernelStatus) string {
	switch code {
	case StatusRunning:
		return "Running"
	case StatusStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

type UIState struct {
	App      *tview.Application
	Pages    *tview.Pages
	Kernel   kernel.Proxy
	Status   KernelStatus
	Ready    bool
	mainFlex *tview.Flex
	//Version string
	Busy atomic.Bool
}

func NewState(app *tview.Application, pages *tview.Pages, kernelAPI kernel.Proxy) *UIState {
	return &UIState{App: app, Pages: pages, Kernel: kernelAPI}
}
