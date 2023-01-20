package tabs

import (
	"strconv"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/cmds"
	"github.com/mjehanno/gtop/model/stats"
	"github.com/mjehanno/gtop/model/styles"
)

type UsageModel struct {
	stats        *stats.Stats
	ramProgress  progress.Model
	swapProgress progress.Model
	cpuProgress  []progress.Model
}

func NewUsageModel() (*UsageModel, error) {
	stat, _ := stats.New()
	model := &UsageModel{
		ramProgress:  progress.New(progress.WithDefaultGradient()),
		swapProgress: progress.New(progress.WithDefaultGradient()),
		cpuProgress:  []progress.Model{},
		stats:        stat,
	}
	return model, nil
}

func (u *UsageModel) Init() tea.Cmd {
	if len(u.cpuProgress) > 0 {
		u.cpuProgress = u.cpuProgress[:0]
	}

	for range u.stats.GetCpuLoad() {
		u.cpuProgress = append(u.cpuProgress, progress.New(progress.WithDefaultGradient()))
	}

	return tea.Batch(u.updateProgressesBars())
}

func (u *UsageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.TickMsg:
		return u, tea.Batch(u.updateProgressesBars(), SyncedTick)
	case progress.FrameMsg:
		cmds := []tea.Cmd{}
		progressRamModel, cmdRam := u.ramProgress.Update(msg)
		u.ramProgress = progressRamModel.(progress.Model)

		cmds = append(cmds, cmdRam)

		progressSwapModel, cmdSwap := u.swapProgress.Update(msg)
		u.swapProgress = progressSwapModel.(progress.Model)
		cmds = append(cmds, cmdSwap)

		for i, cpuBar := range u.cpuProgress {
			model, cmd := cpuBar.Update(msg)
			cmds = append(cmds, cmd)
			u.cpuProgress[i] = model.(progress.Model)
		}

		return u, tea.Batch(cmds...)
	}
	return u, nil
}

func (u *UsageModel) View() string {
	ramLine := styles.LabelStyleRender("Ram usage:") + styles.TabSep + u.ramProgress.View() + styles.SpaceSep + humanize.Bytes(u.stats.GetTotalRam()-u.stats.GetAvailableRam()) + "/" + humanize.Bytes(u.stats.GetTotalRam()) + styles.TabSep + styles.LabelStyleRender("Swap usage:") + styles.TabSep + u.swapProgress.View() + styles.SpaceSep + humanize.Bytes(u.stats.GetTotalSwap()-u.stats.GetAvailableSwap()) + "/" + humanize.Bytes(u.stats.GetTotalSwap()) + styles.Cr

	cpuLines := styles.LabelStyleRender("CPUs:") + styles.Cr
	for i, c := range u.stats.GetCpuLoad() {
		cpuinfoLine := styles.SubLabelStyle("core:") + styles.SpaceSep + strconv.Itoa(c.ProcessorId+1) + styles.TabSep + styles.SubLabelStyle("model:") + styles.SpaceSep + c.ModelName + styles.TabSep + styles.SubLabelStyle("freq:") + styles.SpaceSep + humanize.Ftoa(c.Mhz) + "MHz" + styles.TabSep + styles.SubLabelStyle("cache size:") + styles.SpaceSep + humanize.Bytes(uint64(c.CacheSize)) + styles.Cr
		cpustatLine := styles.SubLabelStyle("usage:") + styles.SpaceSep + u.cpuProgress[i].View() + styles.Cr
		cpuLines += cpuinfoLine + styles.TabSep + cpustatLine
	}

	s := lipgloss.JoinVertical(lipgloss.Left, ramLine, cpuLines)
	return s
}

func updateProgressBar(available, total uint64, bar *progress.Model) tea.Cmd {
	used := total - available
	return bar.SetPercent(float64(used) / float64(total))
}

func (u *UsageModel) updateProgressesBars() tea.Cmd {
	cmds := []tea.Cmd{}
	cpus := u.stats.GetCpuLoad()

	cmds = append(cmds, updateProgressBar(u.stats.GetAvailableRam(), u.stats.GetTotalRam(), &u.ramProgress))
	cmds = append(cmds, updateProgressBar(u.stats.GetAvailableSwap(), u.stats.GetTotalSwap(), &u.swapProgress))

	for i, bar := range u.cpuProgress {
		cmds = append(cmds, updateProgressBar(cpus[i].GetIdle(), cpus[i].GetTotal(), &bar))
	}

	return tea.Batch(cmds...)
}
