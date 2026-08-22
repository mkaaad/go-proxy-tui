package tui

import (
	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

func GetControlFlex(app *tview.Application, pages *tview.Pages, kernelAPI kernel.Proxy) *tview.Flex {
	flex := tview.NewFlex()
	var startButton, endButton, restartButton *tview.Button

	showStart := func() {
		flex.RemoveItem(endButton).RemoveItem(restartButton).AddItem(startButton, 0, 1, true)
		app.SetFocus(startButton)
	}
	showRunning := func() {
		flex.RemoveItem(startButton).AddItem(endButton, 0, 1, true).AddItem(restartButton, 0, 1, true)
		app.SetFocus(endButton)
	}

	startButton = tview.NewButton("start").SetSelectedFunc(func() {
		if err := kernelAPI.Start(); err != nil {
			ShowError(pages, err)
			return
		}
		showRunning()
	})

	endButton = tview.NewButton("stop").SetSelectedFunc(func() {
		if err := kernelAPI.Stop(); err != nil {
			if kernelAPI.Ping() != nil {
				ShowError(pages, err, showStart)
			} else {
				ShowError(pages, err)
			}
			return
		}
		showStart()
	})
	restartButton = tview.NewButton("restart").SetSelectedFunc(func() {
		if err := kernelAPI.Restart(); err != nil {
			if kernelAPI.Ping() != nil {
				ShowError(pages, err, showStart)
			} else {
				ShowError(pages, err)
			}
			return
		}
	})
	flex.AddItem(startButton, 0, 1, true)
	return flex
}
