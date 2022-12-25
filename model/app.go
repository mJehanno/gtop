package model

import (
	"runtime"
	"strconv"
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

var labelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Underline(true)
var labelStyleRender = labelStyle.Render
var subLabelStyle = labelStyle.Copy().Bold(false).Italic(true).Underline(false).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ab2727")).Render
var titleStyle = lipgloss.NewStyle().Margin(1).Padding(0, 2).Align(lipgloss.Center).Foreground(lipgloss.Color("228")).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("63")).Render

type AppModel struct {
	OS           os.Os
	ramProgress  progress.Model
	swapProgress progress.Model
	cpuProgress  []progress.Model
	interfaces   []network.Interface
}

func InitialModel() *AppModel {
	return &AppModel{
		ramProgress:  progress.New(progress.WithDefaultGradient()),
		swapProgress: progress.New(progress.WithDefaultGradient()),
		interfaces:   network.GetInterfaces(),
		cpuProgress:  []progress.Model{},
	}
}

func (a *AppModel) Init() tea.Cmd {
	initOSData(a)
	for range a.OS.Metrics.GetCpuLoad() {
		a.cpuProgress = append(a.cpuProgress, progress.New(progress.WithDefaultGradient()))
	}

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
		initOSData(a)
		return a, tea.Batch(a.updateProgressesBars(), tickCommand(time.Second))
	case progress.FrameMsg:
		cmds := []tea.Cmd{}
		progressRamModel, cmdRam := a.ramProgress.Update(msg)
		a.ramProgress = progressRamModel.(progress.Model)

		cmds = append(cmds, cmdRam)

		progressSwapModel, cmdSwap := a.swapProgress.Update(msg)
		a.swapProgress = progressSwapModel.(progress.Model)
		cmds = append(cmds, cmdSwap)

		for i, cpuBar := range a.cpuProgress {
			model, cmd := cpuBar.Update(msg)
			cmds = append(cmds, cmd)
			a.cpuProgress[i] = model.(progress.Model)
		}

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

	userLine := labelStyleRender("Current User:") + spaceSep + a.OS.Metrics.GetCurrentUser().Uid + spaceSep + a.OS.Metrics.GetCurrentUser().Username + "@" + hostname + tabSep + labelStyleRender("Groups:") + spaceSep + strings.Join(a.OS.Metrics.GetCurrentUser().Groups, ", ") + cr
	systemLine := labelStyleRender("Uptime:") + spaceSep + stringedUptime + tabSep + labelStyleRender("Network:") + spaceSep + strings.Join(netAddresses, ", ") + cr

	ramLine := labelStyleRender("Ram usage:") + tabSep + a.ramProgress.View() + spaceSep + humanize.Bytes(a.OS.Metrics.GetTotalRam()-a.OS.Metrics.GetAvailableRam()) + "/" + humanize.Bytes(a.OS.Metrics.GetTotalRam()) + tabSep + labelStyleRender("Swap usage:") + tabSep + a.swapProgress.View() + spaceSep + humanize.Bytes(a.OS.Metrics.GetTotalSwap()-a.OS.Metrics.GetAvailableSwap()) + "/" + humanize.Bytes(a.OS.Metrics.GetTotalSwap()) + cr

	cpuLines := labelStyleRender("CPUs:") + cr
	for i, c := range a.OS.Metrics.GetCpuLoad() {
		cpuinfoLine := subLabelStyle("core:") + spaceSep + strconv.Itoa(c.ProcessorId+1) + tabSep + subLabelStyle("model:") + spaceSep + c.ModelName + tabSep + subLabelStyle("freq:") + spaceSep + humanize.Ftoa(c.Mhz) + "MHz" + tabSep + subLabelStyle("cache size:") + spaceSep + humanize.Bytes(uint64(c.CacheSize)) + cr
		cpustatLine := subLabelStyle("usage:") + spaceSep + a.cpuProgress[i].View() + cr
		cpuLines += cpuinfoLine + tabSep + cpustatLine
	}

	textBlock := lipgloss.JoinVertical(0.3, userLine, systemLine)
	i := lipgloss.JoinVertical(lipgloss.Left, textBlock, ramLine, cpuLines)

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

		//fmt.Println(a.OS.Metrics.GetCpuLoad())
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
	cmds := []tea.Cmd{}
	cpus := a.OS.Metrics.GetCpuLoad()

	cmds = append(cmds, updateProgressBar(a.OS.Metrics.GetAvailableRam(), a.OS.Metrics.GetTotalRam(), &a.ramProgress))
	cmds = append(cmds, updateProgressBar(a.OS.Metrics.GetAvailableSwap(), a.OS.Metrics.GetTotalSwap(), &a.swapProgress))

	for i, bar := range a.cpuProgress {
		cmds = append(cmds, updateProgressBar(cpus[i].GetIdle(), cpus[i].GetTotal(), &bar))
	}

	return tea.Batch(cmds...)
}

type tickMsg time.Time

func tickCommand(dur time.Duration) tea.Cmd {
	return tea.Tick(dur, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
