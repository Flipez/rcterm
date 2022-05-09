package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func NewProgramm() *tea.Program {
	m := NewBaseModel()
	m.Connection.Connect()

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			select {
			case room := <-m.Connection.RoomChannel:
				p.Send(NewRoomActivity{Room: room})
			case newMessages := <-m.Connection.MessagesChannel:
				p.Send(NewMessagesActivity{Messages: newMessages})
			case msg := <-m.Connection.MessageSubChannel:
				p.Send(NewMessageSubActivity{Message: msg})
			}
		}
	}()

	return p
}

func StartProgram() {
	p := NewProgramm()

	fmt.Println("start program")
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
