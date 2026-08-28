package tui

import (
	"github.com/rivo/tview"
)

const (
	PageSubcription  = "subscription"
	PageKernelConfig = "kernel-config"
)

func (st *UIState) initPages() {
	st.Pages.
		AddPage(PageSubcription, GetSubTabFlex(st), true, false)
}

func GetTabFlex(st *UIState) *tview.Flex {
	subBtn := tview.NewButton("Sub").SetSelectedFunc(func() {
		if st.HasConfig {
			st.Pages.SwitchToPage(PageSubcription)
		}
	})
	proxiesBtn := tview.NewButton("Proxies")
	configBtn := tview.NewButton("Config").SetSelectedFunc(func() {
		if st.HasConfig {
			st.Pages.SwitchToPage(PageKernelConfig)
		}
	})
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(subBtn, 0, 1, true).
		AddItem(proxiesBtn, 0, 1, false).
		AddItem(configBtn, 0, 1, false)
	return flex
}

func GetSubTabFlex(st *UIState) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(GetSubFlex(st), 0, 9, true)
	controlFlex := tview.NewFlex().
		AddItem(GetControlFlex(st), 0, 1, true).
		AddItem(GetKernelFlex(st), 15, 1, false)
	modeFlex := GetModeFlex(st)
	if modeFlex != nil {
		flex.AddItem(modeFlex, 0, 1, false)
	}
	flex.AddItem(controlFlex, 0, 1, false)
	subTabFlex := tview.NewFlex().
		AddItem(flex, 0, 9, true).
		AddItem(GetTabFlex(st), 0, 1, false)
	return subTabFlex

}
func GetProxiesTabFlex(st *UIState) *tview.Flex {
	//TODO
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 9, true)
	controlFlex := tview.NewFlex().
		AddItem(GetControlFlex(st), 0, 1, true).
		AddItem(GetKernelFlex(st), 15, 1, false)
	modeFlex := GetModeFlex(st)
	if modeFlex != nil {
		flex.AddItem(modeFlex, 0, 1, false)
	}
	flex.AddItem(controlFlex, 0, 1, false)
	subTabFlex := tview.NewFlex().
		AddItem(flex, 0, 9, false).
		AddItem(GetTabFlex(st), 0, 1, true)
	return subTabFlex

}
func GetConfigTabFlex(st *UIState) *tview.Flex {
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(GetKernelConfigFlex(st), 0, 9, true)
	//controlFlex := tview.NewFlex().
	//	AddItem(GetControlFlex(st), 0, 1, true).
	//	AddItem(GetKernelFlex(st), 15, 1, false)
	modeFlex := GetModeFlex(st)
	if modeFlex != nil {
		flex.AddItem(modeFlex, 0, 1, false)
	}
	//flex.AddItem(controlFlex, 0, 1, false)
	subTabFlex := tview.NewFlex().
		AddItem(flex, 0, 9, true).
		AddItem(GetTabFlex(st), 0, 1, false)
	return subTabFlex
}
