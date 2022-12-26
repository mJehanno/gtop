package cmds

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func TickCommand(dur time.Duration) tea.Cmd {
	return tea.Tick(dur, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
