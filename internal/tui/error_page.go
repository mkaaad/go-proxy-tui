package tui

import "github.com/rivo/tview"

func ShowError(pages *tview.Pages, err error, onClose ...func()) {
	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"continue"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("error")
			for _, fn := range onClose {
				fn()
			}
		})
	pages.AddPage("error", modal, true, true)
}
