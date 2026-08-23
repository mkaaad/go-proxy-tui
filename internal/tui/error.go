package tui

import "github.com/rivo/tview"

func showError(st *UIState, err error) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			st.Pages.RemovePage("error")
		})
	st.Pages.AddPage("error", modal, true, true)
}
