package main

import (
	"github.com/mkaaad/go-proxy-tui/internal/kernel/mihomo"
	"github.com/mkaaad/go-proxy-tui/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	page := tview.NewPages()
	st := tui.NewState(app, page, &mihomo.Client{})
	tui.ShowKernelConfigPage(st)

	if err := app.SetRoot(page, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
