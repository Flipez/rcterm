package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/flipez/rcterm/tui"
	"github.com/flipez/rcterm/ws"
)

func main() {
	newRoom := make(chan ws.ChatRoom)
	messages := make(chan []ws.Message)
	messageSubs := make(chan ws.Message)

	m := tui.Model{ChannelList: tui.NewChannelList(), MessageList: tui.NewMessageList()}
	m.ChannelList.Title = "Channels"

	m.MessageList.Title = ""

	m.Connection = &ws.Connection{RoomChannel: newRoom, MessagesChannel: messages, MessageSubChannel: messageSubs}
	m.Connection.Connect()

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			select {
			case room := <-newRoom:
				p.Send(tui.NewRoomActivity{Room: room})
			case newMessages := <-messages:
				p.Send(tui.NewMessagesActivity{Messages: newMessages})
			case msg := <-messageSubs:
				p.Send(tui.NewMessageSubActivity{Message: msg})
			}
		}
	}()

	fmt.Println("start program")
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
