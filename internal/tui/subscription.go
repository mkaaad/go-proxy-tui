package tui

import (
	"github.com/rivo/tview"
)

type subscription struct {
	st       *UIState
	list     *tview.List
	selected string
}

func (sub *subscription) showSubForm() {
	subscribeLink := tview.NewInputField().SetLabel("Subscribe Link").SetFieldWidth(20).SetMaskCharacter('*')
	form := tview.NewForm().
		AddFormItem(subscribeLink).
		AddButton("Parse And Save", func() {
			link := subscribeLink.GetText()
			runAsync(sub.st, func() error {
				return sub.st.Kernel.ParseSubLink(link)
			}, func() {
				sub.getNewSubList()
				sub.st.Pages.RemovePage("subscribe form")
			})
		}).
		AddButton("Clear", func() {
			subscribeLink.SetText("")
		}).
		AddButton("Cancel", func() {
			sub.st.Pages.RemovePage("subscribe form")
		})
	sub.st.Pages.AddPage("subscribe form", form, true, true)
}
func GetSubFlex(st *UIState) *tview.Flex {
	sub := &subscription{st: st, list: tview.NewList()}
	addSubBtn := tview.NewButton("Add By Sub Link").SetSelectedFunc(func() {
		sub.showSubForm()
	})
	sub.getNewSubList()
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(sub.list, 0, 9, true).
		AddItem(addSubBtn, 0, 1, false)
	return flex
}
func (sub *subscription) getNewSubList() {
	configs, err := sub.st.Kernel.ListConfigFiles()
	if err != nil {
		showError(sub.st, err)
		return
	}
	sub.list.Clear()
	for i, config := range configs {
		var pre string
		if config.Name == sub.selected {
			pre = "[x]"
		} else {
			pre = "[ ]"
		}
		var shortcut rune
		if i <= 9 {
			shortcut = rune('0' + i)
		} else {
			shortcut = 0
		}
		sub.list.AddItem(pre+config.Name, config.ModTime.Format("2006-01-02"), shortcut, func() {
			runAsync(sub.st, func() error {
				return sub.st.Kernel.ReloadConfigFile(config.Name)
			}, func() {
				sub.selected = config.Name
				sub.getNewSubList()
			})
		})
	}
}

/*
func GetGroupList() {
	groups, err := st.Kernel.ListGroups()
	if err != nil {
		showError(st, err)
		return
	}
	for i, group := range groups {
		list.AddItem("Group | "+group.Name, group.ModTime.Format("2006-01-02"), rune(i), func() {

		})
	}
}
*/
