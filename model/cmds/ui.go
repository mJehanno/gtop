package cmds

import tea "github.com/charmbracelet/bubbletea"

type LeaveOverlayMsg bool

func LeaveOverlay() tea.Cmd {
	return func() tea.Msg {
		return LeaveOverlayMsg(true)
	}
}
