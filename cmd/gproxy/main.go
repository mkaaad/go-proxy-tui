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
	page.AddPage("controlflex", tui.GetControlFlex(app, page, kernel), true, true)
	if err := app.SetRoot(page, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
