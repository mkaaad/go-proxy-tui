package tui

import (
	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

type kernelConfig struct {
	st   *UIState
	list *tview.List
}

func (k *kernelConfig) showKernelConfigForm() {
	kernelConfigNameField := tview.NewInputField().SetLabel("Config Name").SetFieldWidth(30)
	urlField := tview.NewInputField().SetLabel("URL").SetFieldWidth(30)
	secretField := tview.NewInputField().SetLabel("Secret").SetFieldWidth(30).SetMaskCharacter('*')
	form := tview.NewForm().
		AddFormItem(kernelConfigNameField).
		AddFormItem(urlField).
		AddFormItem(secretField).
		AddButton("Create", func() {
			err := k.st.Kernel.NewConfig(kernel.Options{
				Name:   kernelConfigNameField.GetText(),
				URL:    urlField.GetText(),
				Secret: secretField.GetText(),
			})
			if err != nil {
				showError(k.st, err)
				return
			}
			k.st.Pages.RemovePage("kernel config form")
			k.st.Ready = true
			k.getNewConfigList()
			k.st.enterMain()
		}).
		AddButton("Clear", func() {
			kernelConfigNameField.SetText("")
			urlField.SetText("")
			secretField.SetText("")
		}).
		AddButton("Cancel", func() {
			k.st.Pages.RemovePage("kernel config form")
		})
	form.SetBorder(true).SetTitle("Add New Kernel Config")
	k.st.Pages.AddPage("kernel config form", form, true, true)
}

func (k *kernelConfig) getNewConfigList() {
	configs, err := k.st.Kernel.ListConfig()
	if err != nil {
		showError(k.st, err)
		return
	}
	k.list.Clear()
	if len(configs) == 0 {
		k.list.AddItem("No Config File Found", "", 0, func() {})
	}
	for i, config := range configs {
		var shortcut rune
		if i >= 10 {
			shortcut = 0
		} else {
			shortcut = rune('0' + i)
		}
		k.list.AddItem(tview.Escape(config.Name), config.ModTime.Format("2006-01-02"), shortcut, func() {
			err := k.st.Kernel.LoadConfig(config.Name)
			if err != nil {
				showError(k.st, err)
				return
			}
			k.st.Ready = true
			k.st.enterMain()
		})
	}
}

func ShowKernelConfigPage(st *UIState) {
	kc := kernelConfig{
		st:   st,
		list: tview.NewList(),
	}
	kc.getNewConfigList()
	newBtn := tview.NewButton("New").SetSelectedFunc(func() {
		kc.showKernelConfigForm()
	})
	backBtn := tview.NewButton("Back").SetSelectedFunc(func() {
		st.enterMain()
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(kc.list, 0, 9, true).
		AddItem(newBtn, 0, 1, false).
		AddItem(backBtn, 0, 1, false)
	st.Pages.AddPage(pageKernelConfig, flex, true, true)
}
