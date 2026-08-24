package tui

import (
	"time"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

const (
	refreshPerTime = 2 * time.Second
)

type kernelStatus struct {
	state                 *UIState
	statusText            *tview.TextView
	restartBtn            *tview.Button
	lastManualRefreshTime time.Time
}

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
		AddItem(newBtn, 0, 1, true).
		AddItem(backBtn, 0, 1, true)
	st.Pages.AddPage(pageKernelConfig, flex, true, true)
}
func GetKernelFlex(st *UIState) *tview.Flex {
	k := &kernelStatus{state: st}
	k.statusText = tview.NewTextView()
	k.restartBtn = tview.NewButton("Restart").SetSelectedFunc(func() {
		runAsync(k.state, k.state.Kernel.Restart, func() {
			k.refreshStatus()
			k.lastManualRefreshTime = time.Now()
		})
	})
	configBtn := tview.NewButton("Config").SetSelectedFunc(func() {
		k.state.openConfig()
	})
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(k.restartBtn, 0, 4, true).
		AddItem(configBtn, 0, 4, true).
		AddItem(k.statusText, 0, 1, false)
	k.refreshStatus()
	go k.refreshCron()
	return flex
}

func (k *kernelStatus) refreshCron() {

	ticker := time.NewTicker(refreshPerTime)
	defer ticker.Stop()
	for range ticker.C {
		if time.Now().Before(k.lastManualRefreshTime.Add(refreshPerTime)) {
			continue
		}
		k.refreshStatus()
	}

}

func (k *kernelStatus) refreshStatus() {
	pingAsync(k.state, func() {
		err := k.state.Kernel.Ping()
		if err != nil {
			k.state.Status = StatusStopped
		} else {
			k.state.Status = StatusRunning
		}
	}, func() {
		isOnline := k.state.Status == StatusRunning
		k.restartBtn.SetDisabled(!isOnline)
		k.statusText.SetText("Status: " + Status2String(k.state.Status))
	})
}
