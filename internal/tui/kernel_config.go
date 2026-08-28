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
			k.getNewConfigList()
			k.st.HasConfig = true
			k.st.initPages()
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
		k.list.AddItem(config.Name, config.ModTime.Format("2006-01-02"), shortcut, func() {
			k.loadConfig(config.Name)
		})
	}
	k.list.SetBorder(true).SetTitle("Kernel Config")
}

func (k *kernelConfig) loadConfig(configName string) {
	err := k.st.Kernel.LoadConfig(configName)
	if err != nil {
		showError(k.st, err)
		return
	}
	k.st.HasConfig = true
	k.st.initPages()
	k.st.Pages.SwitchToPage(PageSubcription)

}
func GetKernelConfigFlex(st *UIState) *tview.Flex {
	kc := kernelConfig{
		st:   st,
		list: tview.NewList(),
	}
	kc.getNewConfigList()
	newBtn := tview.NewButton("New Kernel Config").SetSelectedFunc(func() {
		kc.showKernelConfigForm()
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(kc.list, 0, 9, true).
		AddItem(newBtn, 0, 1, false)
	return flex
}
