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
	var s *SignalListModel
	items := []list.Item{}

	for i := syscall.SIGABRT; i < syscall.SIGXFSZ; i++ {
		name := unix.SignalName(i)
		items = append(items, signalItem{name: name, description: i.String(), hiddenValue: i})
	}

	s = &SignalListModel{
		List:    list.New(nil, list.NewDefaultDelegate(), 0, 0),
		padding: 1,
	}
	s.List.Title = "Signals"
	s.List.SetShowHelp(false)
	for _, item := range items {
		s.List.InsertItem(len(s.List.Items()), item)
	}

	return s
}

func (s *SignalListModel) Init() tea.Cmd {
	return nil
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
	s.List.SetSize(x-s.padding, y-s.padding)
}

func (s *SignalListModel) View() string {
	sr := ""
	sr += s.List.View()
	return sr
}
