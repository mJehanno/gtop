package model

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
)

type AppModel struct {
	metrics     *MetricModel
	ramProgress progress.Model
}

func InitialModel() *AppModel {
	return &AppModel{
		metrics:     New(),
		ramProgress: progress.New(progress.WithDefaultGradient()),
	}
}

func (a *AppModel) Init() tea.Cmd {
	return tea.Batch(tickCommand(time.Second), a.updateRamProgress())
}

func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	case tickMsg:
		a.metrics = New()
		return a, tea.Batch(a.updateRamProgress(), tickCommand(time.Second))
	case progress.FrameMsg:
		progressModel, cmd := a.ramProgress.Update(msg)
		a.ramProgress = progressModel.(progress.Model)
		return a, cmd

	}
	return a, nil
}

func (a *AppModel) View() string {
	s := "Memory usage : "
	s += a.ramProgress.View() + "   "
	s += humanize.Bytes(a.metrics.FreeRam) + "/" + humanize.Bytes(a.metrics.TotalRam)
	return s
}

func (a *AppModel) updateRamProgress() tea.Cmd {
	return a.ramProgress.SetPercent(float64(a.metrics.FreeRam) / float64(a.metrics.TotalRam))
}

type tickMsg time.Time

func tickCommand(dur time.Duration) tea.Cmd {
	return tea.Tick(dur, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
