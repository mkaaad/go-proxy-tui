package tui

import (
	"errors"
	"strings"
	"time"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

const (
	refreshPerTime = 2 * time.Second
)

type kernelStatus struct {
	state                 *UIState
	statusText            *tview.TextView
	restartBtn            *tview.Button
	lastManualRefreshTime time.Time
}

func GetKernelFlex(st *UIState) *tview.Flex {
	k := &kernelStatus{state: st}
	k.statusText = tview.NewTextView()
	k.state.Status = StatusUnknown
	k.statusText.SetText("Status: " + Status2String(StatusUnknown)).SetTextAlign(tview.AlignCenter)
	k.restartBtn = tview.NewButton("Restart").SetSelectedFunc(func() {
		err := st.Kernel.Restart()
		if err != nil {
			showError(st, err)
		}
		k.refreshStatus()
		k.lastManualRefreshTime = time.Now()
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(k.restartBtn, 0, 1, true).
		AddItem(k.statusText, 1, 1, false)
	k.refreshStatus()
	go k.refreshCron()
	return flex
}

func (k *kernelStatus) refreshCron() {

	ticker := time.NewTicker(refreshPerTime)
	defer ticker.Stop()
	for range ticker.C {
		if time.Now().Before(k.lastManualRefreshTime.Add(refreshPerTime)) {
			continue
		}
		k.refreshStatus()
	}

}

func (k *kernelStatus) refreshStatus() {
	pingAsync(k.state, func() {
		err := k.state.Kernel.Ping()
		switch {
		case errors.Is(err, kernel.ErrNotConfigured):
			k.state.Status = StatusUnknown
		case err != nil:
			k.state.Status = StatusStopped
		default:
			k.state.Status = StatusRunning
		}
	}, func() {
		isOnline := k.state.Status == StatusRunning
		k.restartBtn.SetDisabled(!isOnline)
		k.statusText.SetText("Status: " + Status2String(k.state.Status)).SetTextAlign(tview.AlignCenter)
	})
}
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	lower := strings.ToLower(s)
	return strings.ToUpper(s[0:1]) + lower[1:]
}
