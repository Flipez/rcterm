package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/flipez/rcterm/ws"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
	room        ws.ChatRoom
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list            list.Model
	channelMessages list.Model
	ready           bool
	viewport        viewport.Model
	connection      *ws.Connection
	logs            []string
}

type newRoomActivity struct {
	room ws.ChatRoom
}

type newMessagesActivity struct {
	messages []ws.Message
}

type newLogActivity struct {
	message string
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
			m.list.InsertItem(0, item{title: msg.room.Name, room: msg.room})
		}
	case newMessagesActivity:
		for _, msg := range msg.messages {
			m.channelMessages.InsertItem(0, item{title: msg.Sender.Username, desc: msg.Message})
		}
	case newLogActivity:
		m.logs = append(m.logs, msg.message)
		//m.viewport.SetContent(strings.Join(m.logs, "\n"))
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if i, ok := m.list.SelectedItem().(item); ok {
				m.connection.GetHistory(i.room.Id)
				for i, _ := range m.channelMessages.Items() {
					m.channelMessages.RemoveItem(i)
				}
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.channelMessages.SetSize(msg.Width-h, msg.Height-v)

		if !m.ready {
			m.viewport = viewport.New(msg.Width-h, msg.Height-v)
			m.ready = true
		} else {
			m.viewport.Height = msg.Height - v
			m.viewport.Width = msg.Width - h
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.list.View(),
		m.channelMessages.View(),
		m.viewport.View())
}

func main() {
	newRoom := make(chan ws.ChatRoom)
	messages := make(chan []ws.Message)
	newLog := make(chan string, 1000)

	items := []list.Item{}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	channelMessagesDelegate := list.NewDefaultDelegate()

	m := model{list: list.New(items, delegate, 0, 0), channelMessages: list.New(items, channelMessagesDelegate, 0, 0)}
	m.list.Title = "Channels"
	m.channelMessages.Title = "Channel Name Here"
	m.channelMessages.Styles.StatusBar = lipgloss.NewStyle()

	m.connection = &ws.Connection{RoomChannel: newRoom, MessagesChannel: messages, LogsChannel: newLog}
	m.connection.Connect()

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			select {
			case room := <-newRoom:
				p.Send(newRoomActivity{room})
			case newMessages := <-messages:
				p.Send(newMessagesActivity{newMessages})
			case log := <-newLog:
				p.Send(newLogActivity{log})
			}
		}
	}()

	fmt.Println("start program")
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
