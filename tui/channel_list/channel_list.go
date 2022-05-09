package channel_list

import (
	"github.com/charmbracelet/bubbles/list"
)

func Delegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	return delegate
}

func New(items []list.Item) list.Model {
	return list.New(items, Delegate(), 0, 0)
}
