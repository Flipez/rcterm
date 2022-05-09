package tui

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/flipez/rcterm/tui/channel_list"
	"github.com/flipez/rcterm/tui/message_list"
)

func emtpyItems() []list.Item {
	items := []list.Item{}

	return items
}

func NewChannelList() list.Model {
	return channel_list.New(emtpyItems())
}

func NewMessageList() list.Model {
	return message_list.New(emtpyItems())
}
