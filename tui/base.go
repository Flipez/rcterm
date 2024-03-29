package tui

import (
	"fmt"

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

	ActiveList     int
	StatusBarWidth int
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
			m.MessageList.InsertItem(0, CreateMessageItem(msg))
		}
	case NewMessageSubActivity:
		if m.ActiveRoom.Id == msg.Message.Rid {
			item := CreateMessageItem(msg.Message)
			index := len(m.MessageList.Items())
			m.MessageList.InsertItem(index, item)
			m.MessageList.Select(index)
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
		m.ChannelList.SetSize(msg.Width-h-2, msg.Height-v-2)
		m.MessageList.SetSize(msg.Width-h-2, msg.Height-v-2)
		m.StatusBarWidth = msg.Width - h
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
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					m.ChannelList.View(),
					m.MessageList.View(),
				),
				lipgloss.NewStyle().Width(m.StatusBarWidth).Background(lipgloss.Color("#3C3C3C")).Foreground(lipgloss.Color("202")).Height(1).Align(lipgloss.Right).Render("test"),
			),
		)
}

func NewBaseModel() Model {
	newRoom := make(chan ws.ChatRoom)
	messages := make(chan []ws.Message)
	messageSubs := make(chan ws.Message)

	m := Model{ChannelList: NewChannelList(), MessageList: NewMessageList()}
	m.ChannelList.Title = "Channels"

	m.MessageList.Title = ""

	m.Connection = &ws.Connection{RoomChannel: newRoom, MessagesChannel: messages, MessageSubChannel: messageSubs}

	return m
}
