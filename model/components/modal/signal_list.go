package modal

import (
	"syscall"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/gtop/model/cmds"
	"golang.org/x/sys/unix"
)

type signalItem struct {
	name        string
	description string
	hiddenValue syscall.Signal
}

func (i signalItem) Title() string       { return i.name }
func (i signalItem) Description() string { return i.description }
func (i signalItem) FilterValue() string { return i.name }

type SignalListModel struct {
	Pid     int
	List    list.Model
	padding int
	width   int
	height  int
}

func NewSignalListModel() *SignalListModel {
	s := &SignalListModel{
		List:    list.New(nil, list.NewDefaultDelegate(), 0, 0),
		padding: 1,
	}
	s.List.Title = "Signals"
	s.List.SetShowHelp(false)

	return s
}

func (s *SignalListModel) Init() tea.Cmd {
	items := []list.Item{}
	for i := syscall.Signal(0); i < syscall.Signal(255); i++ {
		name := unix.SignalName(i)
		if name != "" {
			items = append(items, signalItem{name: name, description: i.String(), hiddenValue: i})
		}
	}

	return s.List.SetItems(items)
}

func (s *SignalListModel) Width() int {
	return s.width
}

func (s *SignalListModel) Height() int {
	return s.height
}

func (s *SignalListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, cmds.LeaveOverlay()
		case "enter":
			selected := s.List.Items()[s.List.Index()].(signalItem)
			cmds.SendSignal(uint64(s.Pid), selected.hiddenValue)
		}
	}
	var cmd tea.Cmd
	s.List, cmd = s.List.Update(msg)
	return s, cmd
}

func (s *SignalListModel) SetSize(x, y int) {
	s.width, s.height = x, y
	s.List.SetSize(x, y)
}

func (s *SignalListModel) View() string {
	sr := ""
	sr += s.List.View()
	return sr
}
