package tui

import (
	"github.com/rivo/tview"
)

type ControlPage struct {
	state      *UIState
	statusText *tview.TextView
	startBtn   *tview.Button
	stopBtn    *tview.Button
	restartBtn *tview.Button
}

func (p *ControlPage) Refresh() {
	online := p.state.Status == StatusRuning
	busy := p.state.Busy.Load()
	p.startBtn.SetDisabled(busy || online)
	p.stopBtn.SetDisabled(busy || !online)
	p.restartBtn.SetDisabled(busy || !online)
	p.statusText.SetText("Status: " + Status2String(p.state.Status))
}
func NewControlPage(st *UIState) *tview.Flex {
	p := &ControlPage{state: st}
	p.statusText = tview.NewTextView()
	p.startBtn = tview.NewButton("Start").SetSelectedFunc(func() {
		runAsync(p.state, p.state.Kernel.Start, p.refreshStatus)
	})
	p.stopBtn = tview.NewButton("Stop").SetSelectedFunc(func() {
		runAsync(p.state, p.state.Kernel.Stop, p.refreshStatus)
	})
	p.restartBtn = tview.NewButton("Restart").SetSelectedFunc(func() {
		runAsync(p.state, p.state.Kernel.Restart, p.refreshStatus)
	})
	btnFlex := tview.NewFlex().
		AddItem(p.startBtn, 0, 1, true).
		AddItem(p.stopBtn, 0, 1, false).
		AddItem(p.restartBtn, 0, 1, false)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnFlex, 0, 1, true).
		AddItem(p.statusText, 1, 0, false)
	p.refreshStatus()
	return flex
}
func (p *ControlPage) refreshStatus() {
	go func() {
		err := p.state.Kernel.Ping()
		p.state.App.QueueUpdateDraw(func() {
			if err != nil {
				p.state.Status = StatusStopped
			} else {
				p.state.Status = StatusRuning
			}
			p.Refresh()
		})
	}()
}
