package tui

import (
	"github.com/rivo/tview"
)

type subscription struct {
	st       *UIState
	list     *tview.List
	selected string
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
		AddItem(addSubBtn, 3, 1, false)

	flex.SetBorder(true).SetTitle("Subscriptions")
	return flex
}

func (sub *subscription) showSubForm() {
	subscribeLink := tview.NewInputField().SetLabel("Sub Link").SetFieldWidth(30).SetMaskCharacter('*')
	form := tview.NewForm().
		AddFormItem(subscribeLink).
		AddButton("Parse And Save", func() {
			link := subscribeLink.GetText()
			err := sub.st.Kernel.ParseSubLink(link)
			if err != nil {
				showError(sub.st, err)
				return
			}
			sub.getNewSubList()
			sub.st.Pages.RemovePage("subscribe form")
		}).
		AddButton("Clear", func() {
			subscribeLink.SetText("")
		}).
		AddButton("Cancel", func() {
			sub.st.Pages.RemovePage("subscribe form")
		})
	sub.st.Pages.AddPage("subscribe form", form, true, true)
}

func (sub *subscription) getNewSubList() {
	configs, err := sub.st.Kernel.ListSubConfig()
	if err != nil {
		showError(sub.st, err)
		return
	}
	sub.list.Clear()
	if len(configs) == 0 {
		sub.list.AddItem("No Config File Found", "", 0, func() {})
	}
	for i, config := range configs {
		var pre string
		isSlected := config.Name == sub.selected
		if isSlected {
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
		sub.list.AddItem(tview.Escape(pre+config.Name), config.ModTime.Format("2006-01-02"), shortcut, func() {
			err := sub.st.Kernel.LoadSubConfig(config.Name)
			if err != nil {
				showError(sub.st, err)
			}
			sub.selected = config.Name
			sub.getNewSubList()
		})
		if isSlected {
			sub.list.SetCurrentItem(i)
		}
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
