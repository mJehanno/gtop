package tabs

import tea "github.com/charmbracelet/bubbletea"

type NetworkTabModel struct{}

func NewNetworkTabModel() (*NetworkTabModel, error) {
	return &NetworkTabModel{}, nil
}

func (n *NetworkTabModel) Init() tea.Cmd {
	return nil
}
func (n *NetworkTabModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return n, nil
}
func (n *NetworkTabModel) View() string {
	s := ""
	return s
}
