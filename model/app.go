package model

import (
	"os/user"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/metrics"
)

type AppModel struct {
	user         *user.User
	metrics      *metrics.MetricModel
	ramProgress  progress.Model
	swapProgress progress.Model
}

func InitialModel() *AppModel {
	user, _ := user.Current()
	return &AppModel{
		user:         user,
		metrics:      metrics.New(),
		ramProgress:  progress.New(progress.WithDefaultGradient()),
		swapProgress: progress.New(progress.WithDefaultGradient()),
	}
}

func (a *AppModel) Init() tea.Cmd {
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
		a.metrics = metrics.New()
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
	usedRam := a.metrics.TotalRam - a.metrics.BufferedRam - a.metrics.FreeRam
	usedSwap := a.metrics.TotalSwap - a.metrics.FreeSwap
	groupIds, _ := a.user.GroupIds()
	s := "Current user : " + a.user.Uid + " " + a.user.Username + "   Groups : "
	for i, id := range groupIds {
		group, _ := user.LookupGroupId(id)
		s += group.Name
		if i != len(groupIds)-1 {
			s += ","
		}
	}
	s += "\n"
	s += "Uptime : " + humanize.Time(time.Now().Add(-time.Duration(a.metrics.Uptime)*time.Second)) + "\n"
	s += "Memory usage : "
	s += a.ramProgress.View() + "   "
	s += humanize.Bytes(usedRam) + "/" + humanize.Bytes(a.metrics.TotalRam) + "\n"
	s += "Swap usage : " + a.swapProgress.View() + "   " + humanize.Bytes(usedSwap) + "/" + humanize.Bytes(a.metrics.TotalSwap)

	return s
}

func updateProgressBar(free, buffered, total uint64, bar *progress.Model) tea.Cmd {
	used := total - buffered - free
	return bar.SetPercent(float64(used) / float64(total))

}

func (a *AppModel) updateProgressesBars() tea.Cmd {
	return tea.Batch(updateProgressBar(a.metrics.FreeRam, a.metrics.BufferedRam, a.metrics.TotalRam, &a.ramProgress), updateProgressBar(a.metrics.FreeSwap, 0, a.metrics.TotalSwap, &a.swapProgress))
}

type tickMsg time.Time

func tickCommand(dur time.Duration) tea.Cmd {
	return tea.Tick(dur, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
