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
		err := p.state.Kernel.Enable()
		if err != nil {
			showError(st, err)
		}
		st.App.SetFocus(nil)

	})
	p.disableBtn = tview.NewButton("Disable").SetSelectedFunc(func() {
		err := p.state.Kernel.Disable()
		if err != nil {
			showError(st, err)
		}
		st.App.SetFocus(nil)

	})
	flex := tview.NewFlex().
		AddItem(p.enableBtn, 0, 1, true).
		AddItem(p.disableBtn, 0, 1, false)
	return flex
}
