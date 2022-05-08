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
	list     list.Model
	ready    bool
	viewport viewport.Model
}

type newRoomActivity struct {
	room ws.ChatRoom
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
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if i, ok := m.list.SelectedItem().(item); ok {
				m.viewport.SetContent(i.room.Id)
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

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
		m.viewport.View())
}

func main() {
	newRoom := make(chan ws.ChatRoom)

	items := []list.Item{}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	m := model{list: list.New(items, delegate, 0, 0)}
	m.list.Title = "Channels"

	go ws.ConnectServer(newRoom)
	p := tea.NewProgram(m, tea.WithAltScreen())

	go func() {
		for {
			select {
			case room := <-newRoom:
				p.Send(newRoomActivity{room})
			}
		}
	}()

	fmt.Println("start program")
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
