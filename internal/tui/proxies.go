package tui

import (
	"fmt"

	"github.com/mkaaad/go-proxy-tui/internal/kernel"
	"github.com/rivo/tview"
)

const (
	triangleCollapsed = "▸ "
	triangleExpanded  = "▾ "
)

type proxiesPage struct {
	st   *UIState
	tree *tview.TreeView
}

func GetProxiesFlex(st *UIState) *tview.Flex {
	p := &proxiesPage{st: st, tree: tview.NewTreeView()}
	p.rebuild()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(p.tree, 0, 9, true)
	flex.SetBorder(true).SetTitle("Proxies")
	return flex
}

func (p *proxiesPage) rebuild() {
	root := tview.NewTreeNode("").SetExpanded(true)
	proxies, err := p.st.Kernel.ListProxies()
	if err != nil {
		showError(p.st, err)
		root.AddChild(tview.NewTreeNode("Failed to load proxies"))
		p.tree.SetRoot(root).SetTopLevel(1)
		p.bindToggle()
		return
	}
	if len(proxies) == 0 {
		root.AddChild(tview.NewTreeNode("No proxies found"))
	}

	byName := make(map[string]kernel.ProxyInfo, len(proxies))
	members := make(map[string]bool)
	globalNow := ""
	for _, proxy := range proxies {
		byName[proxy.Name] = proxy
		if proxy.Name == "GLOBAL" && isProxyGroup(proxy.Type) {
			globalNow = proxy.Now
		}
		if isProxyGroup(proxy.Type) {
			for _, m := range proxy.All {
				members[m] = true
			}
		}
	}

	for _, proxy := range proxies {
		if isProxyGroup(proxy.Type) {
			label := fmt.Sprintf("%s (%s) [%d]", proxy.Name, proxy.Type, len(proxy.All))
			if proxy.Now != "" {
				label += " → " + proxy.Now
			}
			node := tview.NewTreeNode(triangleCollapsed + label)
			node.SetReference(label)
			node.SetExpanded(false)
			for _, member := range proxy.All {
				label := memberLabel(member, proxy.Now)
				if info, ok := byName[member]; ok {
					label += delaySuffix(info)
				}
				node.AddChild(tview.NewTreeNode(label))
			}
			root.AddChild(node)
		} else {
			label := memberLabel(proxy.Name+" ("+proxy.Type+")", globalNow)
			label += delaySuffix(proxy)
			root.AddChild(tview.NewTreeNode(label))
		}
	}
	p.tree.SetRoot(root).SetTopLevel(1)
	p.bindToggle()
}

// bindToggle makes activating a node with children toggle its expansion and
// updates the triangle marker: mouse click, Enter and Space all route through
// TreeView's selected callback.
func (p *proxiesPage) bindToggle() {
	p.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if len(node.GetChildren()) == 0 {
			return
		}
		expanded := !node.IsExpanded()
		node.SetExpanded(expanded)
		if ref, ok := node.GetReference().(string); ok && ref != "" {
			if expanded {
				node.SetText(triangleExpanded + ref)
			} else {
				node.SetText(triangleCollapsed + ref)
			}
		}
	})
}

func isProxyGroup(typ string) bool {
	switch typ {
	case "Selector", "URLTest", "Fallback", "LoadBalance":
		return true
	}
	return false
}

func memberLabel(name, now string) string {
	label := "[ ] " + name
	if name == now {
		label = "[x] " + name
	}
	return label
}

func delaySuffix(p kernel.ProxyInfo) string {
	if !p.Alive {
		return " [dead]"
	}
	if p.Delay > 0 {
		return fmt.Sprintf(" [%dms]", p.Delay)
	}
	return " [-]"
}
