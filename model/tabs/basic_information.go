package tabs

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/cmds"
	"github.com/mjehanno/gtop/model/metrics/linux"
	"github.com/mjehanno/gtop/model/metrics/os"
	"github.com/mjehanno/gtop/model/network"
	"github.com/mjehanno/gtop/model/styles"
)

type BasicInformationModel struct {
	OS           os.Os
	ramProgress  progress.Model
	swapProgress progress.Model
	cpuProgress  []progress.Model
	interfaces   []network.Interface
}

func NewBasincInformationModel() *BasicInformationModel {
	return &BasicInformationModel{
		ramProgress:  progress.New(progress.WithDefaultGradient()),
		swapProgress: progress.New(progress.WithDefaultGradient()),
		interfaces:   network.GetInterfaces(),
		cpuProgress:  []progress.Model{},
	}
}

func (b *BasicInformationModel) Init() tea.Cmd {
	initOSData(b)
	if len(b.cpuProgress) > 0 {
		b.cpuProgress = b.cpuProgress[:0]
	}
	for range b.OS.Metrics.GetCpuLoad() {
		b.cpuProgress = append(b.cpuProgress, progress.New(progress.WithDefaultGradient()))
	}

	return tea.Batch(b.updateProgressesBars())
}

func (b *BasicInformationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.TickMsg:
		initOSData(b)
		return b, tea.Batch(b.updateProgressesBars(), SyncedTick)
	case progress.FrameMsg:
		cmds := []tea.Cmd{}
		progressRamModel, cmdRam := b.ramProgress.Update(msg)
		b.ramProgress = progressRamModel.(progress.Model)

		cmds = append(cmds, cmdRam)

		progressSwapModel, cmdSwap := b.swapProgress.Update(msg)
		b.swapProgress = progressSwapModel.(progress.Model)
		cmds = append(cmds, cmdSwap)

		for i, cpuBar := range b.cpuProgress {
			model, cmd := cpuBar.Update(msg)
			cmds = append(cmds, cmd)
			b.cpuProgress[i] = model.(progress.Model)
		}

		return b, tea.Batch(cmds...)
	}
	return b, nil
}

func (b *BasicInformationModel) View() string {
	netAddresses := make([]string, len(b.interfaces))

	for i, a := range b.interfaces {
		netAddresses[i] = a.String()
	}

	hostname, err := b.OS.Metrics.GetHostname()
	if err != nil {
		hostname = styles.ErrorStyle("- error while getting hostname")
	}

	var stringedUptime string

	uptime, err := b.OS.Metrics.GetUptime()
	if err != nil {
		stringedUptime = styles.ErrorStyle("- error while getting uptime")
	} else {
		stringedUptime = humanize.Time(time.Now().Add(-time.Second * time.Duration(uptime)))
	}

	userLine := styles.LabelStyleRender("Current User:") + styles.SpaceSep + b.OS.Metrics.GetCurrentUser().Uid + styles.SpaceSep + b.OS.Metrics.GetCurrentUser().Username + "@" + hostname + styles.TabSep + styles.LabelStyleRender("Groups:") + styles.SpaceSep + strings.Join(b.OS.Metrics.GetCurrentUser().Groups, ", ") + styles.Cr
	systemLine := styles.LabelStyleRender("Uptime:") + styles.SpaceSep + stringedUptime + styles.TabSep + styles.LabelStyleRender("Network:") + styles.SpaceSep + strings.Join(netAddresses, ", ") + styles.Cr

	ramLine := styles.LabelStyleRender("Ram usage:") + styles.TabSep + b.ramProgress.View() + styles.SpaceSep + humanize.Bytes(b.OS.Metrics.GetTotalRam()-b.OS.Metrics.GetAvailableRam()) + "/" + humanize.Bytes(b.OS.Metrics.GetTotalRam()) + styles.TabSep + styles.LabelStyleRender("Swap usage:") + styles.TabSep + b.swapProgress.View() + styles.SpaceSep + humanize.Bytes(b.OS.Metrics.GetTotalSwap()-b.OS.Metrics.GetAvailableSwap()) + "/" + humanize.Bytes(b.OS.Metrics.GetTotalSwap()) + styles.Cr

	cpuLines := styles.LabelStyleRender("CPUs:") + styles.Cr
	for i, c := range b.OS.Metrics.GetCpuLoad() {
		cpuinfoLine := styles.SubLabelStyle("core:") + styles.SpaceSep + strconv.Itoa(c.ProcessorId+1) + styles.TabSep + styles.SubLabelStyle("model:") + styles.SpaceSep + c.ModelName + styles.TabSep + styles.SubLabelStyle("freq:") + styles.SpaceSep + humanize.Ftoa(c.Mhz) + "MHz" + styles.TabSep + styles.SubLabelStyle("cache size:") + styles.SpaceSep + humanize.Bytes(uint64(c.CacheSize)) + styles.Cr
		cpustatLine := styles.SubLabelStyle("usage:") + styles.SpaceSep + b.cpuProgress[i].View() + styles.Cr
		cpuLines += cpuinfoLine + styles.TabSep + cpustatLine
	}

	textBlock := lipgloss.JoinVertical(lipgloss.Left, userLine, systemLine)
	s := lipgloss.JoinVertical(lipgloss.Left, textBlock, ramLine, cpuLines)

	return s
}

func initOSData(b *BasicInformationModel) {
	switch runtime.GOOS {
	case "linux":
		metrics, _ := linux.New()
		b.OS = os.Os{
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

func (b *BasicInformationModel) updateProgressesBars() tea.Cmd {
	cmds := []tea.Cmd{}
	cpus := b.OS.Metrics.GetCpuLoad()

	cmds = append(cmds, updateProgressBar(b.OS.Metrics.GetAvailableRam(), b.OS.Metrics.GetTotalRam(), &b.ramProgress))
	cmds = append(cmds, updateProgressBar(b.OS.Metrics.GetAvailableSwap(), b.OS.Metrics.GetTotalSwap(), &b.swapProgress))

	for i, bar := range b.cpuProgress {
		cmds = append(cmds, updateProgressBar(cpus[i].GetIdle(), cpus[i].GetTotal(), &bar))
	}

	return tea.Batch(cmds...)
}
