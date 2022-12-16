package model

import (
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/gtop/model/metrics/linux"
	"github.com/mjehanno/gtop/model/metrics/os"
	"github.com/mjehanno/gtop/model/network"
)

type AppModel struct {
	OS           os.Os
	ramProgress  progress.Model
	swapProgress progress.Model
	interfaces   []network.Interface
}

func InitialModel() *AppModel {
	return &AppModel{
		ramProgress:  progress.New(progress.WithDefaultGradient()),
		swapProgress: progress.New(progress.WithDefaultGradient()),
		interfaces:   network.GetInterfaces(),
	}
}

func (a *AppModel) Init() tea.Cmd {
	initOSData(a)
	return tea.Batch(tickCommand(time.Second), a.updateProgressesBars())
}

func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	case tickMsg:
		return a, tea.Batch(a.updateProgressesBars(), tickCommand(time.Second))
	case progress.FrameMsg:
		cmds := []tea.Cmd{}
		progressRamModel, cmdRam := a.ramProgress.Update(msg)
		a.ramProgress = progressRamModel.(progress.Model)

		cmds = append(cmds, cmdRam)

		progressSwapModel, cmdSwap := a.swapProgress.Update(msg)
		a.swapProgress = progressSwapModel.(progress.Model)

		cmds = append(cmds, cmdSwap)

		return a, tea.Batch(cmds...)

	}
	return a, nil
}

func (a *AppModel) View() string {
	s := ""
	s += "Current User : " + a.OS.Metrics.GetCurrentUser().Uid + " " + a.OS.Metrics.GetCurrentUser().Username + "   " + strings.Join(a.OS.Metrics.GetCurrentUser().Groups, ",")
	return s
}

func initOSData(a *AppModel) {
	switch runtime.GOOS {
	case "linux":
		a.OS = os.Os{
			Metrics: &linux.LinuxMetric{},
		}
	}
}

func updateProgressBar(available, total uint64, bar *progress.Model) tea.Cmd {
	used := total - available
	return bar.SetPercent(float64(used) / float64(total))

}

func (a *AppModel) updateProgressesBars() tea.Cmd {
	return tea.Batch(updateProgressBar(a.OS.Metrics.GetAvailableRam(), a.OS.Metrics.GetTotalRam(), &a.ramProgress), updateProgressBar(a.OS.Metrics.GetAvailableSwap(), a.OS.Metrics.GetTotalSwap(), &a.swapProgress))
}

type tickMsg time.Time

func tickCommand(dur time.Duration) tea.Cmd {
	return tea.Tick(dur, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
