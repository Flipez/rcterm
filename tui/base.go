package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/flipez/rcterm/ws"
)

type NewRoomActivity struct {
	Room ws.ChatRoom
}

type NewMessagesActivity struct {
	Messages []ws.Message
}

type NewMessageSubActivity struct {
	Message ws.Message
}

type Model struct {
	ChannelList list.Model
	MessageList list.Model
	Connection  *ws.Connection
	ActiveRoom  ws.ChatRoom

	ActiveList int
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case NewRoomActivity:
		if msg.Room.Name != "" {
			m.ChannelList.InsertItem(0, NewChannelListItem(msg.Room.Name, msg.Room))
		}
	case NewMessagesActivity:
		for _, msg := range msg.Messages {
			title := fmt.Sprintf("%s @ %s", msg.Sender.Username, time.Unix(0, int64(msg.Date.Timestamp)*int64(time.Millisecond)))
			m.MessageList.InsertItem(0, NewMessageListItem(title, msg.Message))
		}
	case NewMessageSubActivity:
		if m.ActiveRoom.Id == msg.Message.Rid {
			m.MessageList.InsertItem(-1, NewMessageListItem(msg.Message.Sender.Username, msg.Message.Message))
		} else {
			m.MessageList.NewStatusMessage(fmt.Sprintf("New message in %s from %s!", msg.Message.Rid, msg.Message.Sender.Username))
		}

	case tea.KeyMsg:
		if msg.String() == "tab" {
			m.ActiveList = (m.ActiveList + 1) % 2
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if i, ok := m.ChannelList.SelectedItem().(ChannelListItem); ok {
				m.ActiveRoom = i.Room
				m.MessageList.Title = i.Room.Name
				m.Connection.OpenRoom(i.Room.Id)
				m.Connection.GetHistory(i.Room.Id)
				for i := range m.MessageList.Items() {
					m.MessageList.RemoveItem(i)
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.ChannelList.SetSize(msg.Width-h, msg.Height-v)
		m.MessageList.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.ActiveList == 0 {
		m.ChannelList, cmd = m.ChannelList.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.MessageList, cmd = m.MessageList.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.ChannelList.View(),
		m.MessageList.View(),
	)
}
