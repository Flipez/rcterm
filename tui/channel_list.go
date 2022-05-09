package tui

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/flipez/rcterm/ws"
)

type ChannelListItem struct {
	title       string
	description string
	Room        ws.ChatRoom
}

func (cl ChannelListItem) Title() string       { return cl.title }
func (cl ChannelListItem) Description() string { return cl.description }
func (cl ChannelListItem) FilterValue() string { return cl.title }

func ChannelListDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	return delegate
}

func NewChannelList() list.Model {
	return list.New(EmptyItems(), ChannelListDelegate(), 0, 0)
}

func NewChannelListItem(title string, room ws.ChatRoom) ChannelListItem {
	return ChannelListItem{title: title, Room: room}
}
