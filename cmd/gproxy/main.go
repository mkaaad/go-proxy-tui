package main

import (
	"github.com/mkaaad/go-proxy-tui/internal/kernel/mihomo"
	"github.com/mkaaad/go-proxy-tui/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	kernel := mihomo.New()
	page := tview.NewPages()
	st := tui.NewState(app, page, kernel)
	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.GetSubFlex(st), 0, 9, true).
		AddItem(tui.GetControlFlex(st), 0, 1, false)
	page.AddPage("mainflex", mainFlex, true, true)
	if err := app.SetRoot(page, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
