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
			AddItem(GetControlFlex(st), 0, 8, true).
			AddItem(GetKernelFlex(st), 0, 2, false)
		st.mainFlex = tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(GetSubFlex(st), 0, 9, true).
			AddItem(controlFlex, 3, 1, false)
	}
	st.Pages.RemovePage(pageMain)
	st.Pages.AddPage(pageMain, st.mainFlex, true, true)
}

func (st *UIState) openConfig() {
	st.Pages.SwitchToPage(pageKernelConfig)
}
