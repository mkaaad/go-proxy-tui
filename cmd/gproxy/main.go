package main

import (
	"github.com/mkaaad/go-proxy-tui/internal/kernel/mihomo"
	"github.com/mkaaad/go-proxy-tui/internal/tui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	kernel := &mihomo.Client{}
	page := tview.NewPages()
	st := tui.NewState(app, page, kernel)
	mainFlex := getMainFlex(st)
	done := make(chan struct{})
	go func() {
		<-done
		tui.ShowKernelConfigPage(st, done)
		page.AddPage("mainflex", mainFlex, true, true)
	}()

	if err := app.SetRoot(page, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
func getMainFlex(st *tui.UIState) *tview.Flex {
	controlFlex := tview.NewFlex().
		AddItem(tui.GetControlFlex(st), 0, 8, true).
		AddItem(tui.GetKernelFlex(st), 0, 2, false)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tui.GetSubFlex(st), 0, 9, true).
		AddItem(controlFlex, 3, 1, false)
	return flex
}
