package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type modeStruct struct {
	modeBtns []*tview.Button
}

func GetModeFlex(st *UIState) *tview.Flex {
	modes, err := st.Kernel.ListModes()
	if err != nil {
		showError(st, err)
	}
	flex := tview.NewFlex()
	m := &modeStruct{}
	for _, mode := range modes {
		modeBtn := tview.NewButton(capitalize(mode))
		m.modeBtns = append(m.modeBtns, modeBtn)
		flex.AddItem(
			modeBtn.SetSelectedFunc(func() {
				err := st.Kernel.SwitchMode(mode)
				st.App.SetFocus(nil)
				if err != nil {
					showError(st, err)
					return
				}
				m.clearBtnColor()
				modeBtn.SetLabelColor(tcell.Color100)
			}), 0, 1, true)
	}
	return flex
}

func (m *modeStruct) clearBtnColor() {
	for _, btn := range m.modeBtns {
		btn.SetLabelColor(tcell.Color255)
	}
}
