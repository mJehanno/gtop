package model

import (
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/metrics/linux"
	"github.com/mjehanno/gtop/model/metrics/os"
	"github.com/mjehanno/gtop/model/network"
)

var labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Underline(true).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ab2727")).Render
var titleStyle = lipgloss.NewStyle().Margin(1).Padding(0, 2).Align(lipgloss.Center).Foreground(lipgloss.Color("228")).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Render

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
	tabSep := "    "
	spaceSep := " "
	cr := "\n"

	netAddresses := make([]string, len(a.interfaces))

	for i, a := range a.interfaces {
		netAddresses[i] = a.String()
	}

	hostname, err := a.OS.Metrics.GetHostname()
	if err != nil {
		hostname = errorStyle("- error while getting hostname")
	}

	var stringedUptime string

	uptime, err := a.OS.Metrics.GetUptime()
	if err != nil {
		stringedUptime = errorStyle("- error while getting uptime")
	} else {
		stringedUptime = humanize.Time(time.Now().Add(-time.Second * time.Duration(uptime)))
	}

	userLine := labelStyle("Current User:") + spaceSep + a.OS.Metrics.GetCurrentUser().Uid + spaceSep + a.OS.Metrics.GetCurrentUser().Username + "@" + hostname + tabSep + labelStyle("Groups:") + spaceSep + strings.Join(a.OS.Metrics.GetCurrentUser().Groups, ", ") + cr
	systemLine := labelStyle("Uptime:") + spaceSep + stringedUptime + tabSep + labelStyle("Network:") + spaceSep + strings.Join(netAddresses, ", ") + cr

	ramLine := labelStyle("Ram usage:") + tabSep + a.ramProgress.View() + spaceSep + humanize.Bytes(a.OS.Metrics.GetTotalRam()-a.OS.Metrics.GetAvailableRam()) + "/" + humanize.Bytes(a.OS.Metrics.GetTotalRam()) + tabSep + labelStyle("Swap usage:") + tabSep + a.swapProgress.View() + spaceSep + humanize.Bytes(a.OS.Metrics.GetTotalSwap()-a.OS.Metrics.GetAvailableSwap()) + "/" + humanize.Bytes(a.OS.Metrics.GetTotalSwap()) + cr

	textBlock := lipgloss.JoinVertical(0.3, userLine, systemLine)
	ramBlock := lipgloss.JoinVertical(lipgloss.Left, ramLine)

	i := lipgloss.JoinVertical(lipgloss.Left, textBlock, ramBlock)

	s := lipgloss.PlaceHorizontal(120, lipgloss.Center, titleStyle("GTop")) + cr
	s += lipgloss.PlaceHorizontal(240, lipgloss.Left, i)

	return s
}

func initOSData(a *AppModel) {
	switch runtime.GOOS {
	case "linux":
		metrics, _ := linux.New()
		a.OS = os.Os{
			Metrics: metrics,
		}
	case "windows":
	case "darwin":
	case "bsd":
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
