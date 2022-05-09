package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/flipez/rcterm/tui"
	"github.com/flipez/rcterm/ws"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	channelList     list.Model
	channelMessages list.Model
	connection      *ws.Connection
	activeRoom      ws.ChatRoom

	activeList int
}

type newRoomActivity struct {
	room ws.ChatRoom
}

type newMessagesActivity struct {
	messages []ws.Message
}

type newMessageSubActivity struct {
	message ws.Message
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case newRoomActivity:
		if msg.room.Name != "" {
			m.channelList.InsertItem(0, tui.NewChannelListItem(msg.room.Name, msg.room))
		}
	case newMessagesActivity:
		for _, msg := range msg.messages {
			title := fmt.Sprintf("%s @ %s", msg.Sender.Username, time.Unix(0, int64(msg.Date.Timestamp)*int64(time.Millisecond)))
			m.channelMessages.InsertItem(0, tui.NewMessageListItem(title, msg.Message))
		}
	case newMessageSubActivity:
		if m.activeRoom.Id == msg.message.Rid {
			m.channelMessages.InsertItem(-1, tui.NewMessageListItem(msg.message.Sender.Username, msg.message.Message))
		} else {
			m.channelMessages.NewStatusMessage(fmt.Sprintf("New message in %s from %s!", msg.message.Rid, msg.message.Sender.Username))
		}

	case tea.KeyMsg:
		if msg.String() == "tab" {
			m.activeList = (m.activeList + 1) % 2
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if i, ok := m.channelList.SelectedItem().(tui.ChannelListItem); ok {
				m.activeRoom = i.Room
				m.channelMessages.Title = i.Room.Name
				m.connection.OpenRoom(i.Room.Id)
				m.connection.GetHistory(i.Room.Id)
				for i := range m.channelMessages.Items() {
					m.channelMessages.RemoveItem(i)
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.channelList.SetSize(msg.Width-h, msg.Height-v)
		m.channelMessages.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.activeList == 0 {
		m.channelList, cmd = m.channelList.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.channelMessages, cmd = m.channelMessages.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.channelList.View(),
		m.channelMessages.View(),
	)
}

func main() {
	newRoom := make(chan ws.ChatRoom)
	messages := make(chan []ws.Message)
	messageSubs := make(chan ws.Message)

	m := model{channelList: tui.NewChannelList(), channelMessages: tui.NewMessageList()}
	m.channelList.Title = "Channels"

	m.channelMessages.Title = ""

	m.connection = &ws.Connection{RoomChannel: newRoom, MessagesChannel: messages, MessageSubChannel: messageSubs}
	m.connection.Connect()

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			select {
			case room := <-newRoom:
				p.Send(newRoomActivity{room})
			case newMessages := <-messages:
				p.Send(newMessagesActivity{newMessages})
			case msg := <-messageSubs:
				p.Send(newMessageSubActivity{msg})
			}
		}
	}()

	fmt.Println("start program")
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
