package tui

import "github.com/rivo/tview"

const (
	pageMain         = "main"
	pageKernelConfig = "kernel-config"
)

func (st *UIState) enterMain() {
	if !st.Ready {
		return
	}
	if st.mainFlex == nil {
		controlFlex := tview.NewFlex().
			AddItem(GetControlFlex(st), 0, 1, true).
			AddItem(GetKernelFlex(st), 15, 1, false)
		st.mainFlex = tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(GetSubFlex(st), 0, 9, true)
		modeFlex := GetModeFlex(st)
		if modeFlex != nil {
			st.mainFlex.AddItem(modeFlex, 0, 1, false)
		}
		st.mainFlex.AddItem(controlFlex, 0, 1, false)
	}
	st.Pages.RemovePage(pageMain)
	st.Pages.AddPage(pageMain, st.mainFlex, true, true)
}

func (st *UIState) openConfig() {
	st.Pages.SwitchToPage(pageKernelConfig)
}
