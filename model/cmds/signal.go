package cmds

import (
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

type SignalSend struct {
	pid uint64
	sig syscall.Signal
	err error
}

func SendSignal(pid uint64, sig syscall.Signal) tea.Msg {
	err := syscall.Kill(int(pid), sig)
	return SignalSend{
		pid,
		sig,
		err,
	}
}
