package tui

import (
	"github.com/rivo/tview"
)

type ControlPage struct {
	state      *UIState
	enableBtn  *tview.Button
	disableBtn *tview.Button
}

func GetControlFlex(st *UIState) *tview.Flex {
	p := &ControlPage{state: st}
	p.enableBtn = tview.NewButton("Enable").SetSelectedFunc(func() {
		runAsync(p.state, p.state.Kernel.Start, nil)
	})
	p.disableBtn = tview.NewButton("Disable").SetSelectedFunc(func() {
		runAsync(p.state, p.state.Kernel.Stop, nil)
	})
	flex := tview.NewFlex().
		AddItem(p.enableBtn, 0, 1, true).
		AddItem(p.disableBtn, 0, 1, false)
	return flex
}

/*func (p *ControlPage) Refresh() {
	online := p.state.Status == StatusRunning
	p.startBtn.SetDisabled(online)
	p.stopBtn.SetDisabled(!online)
}
func (p *ControlPage) refreshStatus() {
	go func() {
		err := p.state.Kernel.Ping()
		p.state.App.QueueUpdateDraw(func() {
			if err != nil {
				p.state.Status = StatusStopped
			} else {
				p.state.Status = StatusRunning
			}
			p.Refresh()
		})
	}()
}*/
