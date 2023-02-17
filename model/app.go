package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mjehanno/gtop/model/components/tabs"
	"github.com/mjehanno/gtop/model/styles"
)

type AppModel struct {
	currentTab    tabs.TabsEnum
	tabs          []tabs.Tab
	width, height int
}

func InitialModel() *AppModel {
	return &AppModel{
		tabs: tabs.GetAllTabs(),
	}
}

func (a *AppModel) Init() tea.Cmd {
	return tea.Batch(a.tabs[a.currentTab].Init(), tabs.SyncedTick)
}

func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.tabs[a.currentTab].Update(msg)
		a.width, a.height = msg.Width, msg.Height
		a.tabs[a.currentTab].SetSize(a.width, a.height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "tab":
			a.currentTab++
			if int(a.currentTab) >= len(a.tabs) {
				a.currentTab = 0
			}
			a.tabs[a.currentTab].SetSize(a.width, a.height)
			return a, a.tabs[a.currentTab].Init()
		case "shift+tab":
			a.currentTab--
			if int(a.currentTab) < 0 {
				a.currentTab = tabs.TabsEnum(len(a.tabs) - 1)
			}
			a.tabs[a.currentTab].SetSize(a.width, a.height)
			return a, a.tabs[a.currentTab].Init()
		}
	}

	m, cmd := a.tabs[a.currentTab].Update(msg)
	a.tabs[a.currentTab].Model = m
	return a, cmd
}

func (a *AppModel) View() string {
	s := ""

	tabs := []string{}

	for i, tab := range a.tabs {
		if i == int(a.currentTab) {
			tabs = append(tabs, styles.ActivatedTabStyle(tab.Name))
			continue
		}
		tabs = append(tabs, styles.DeactivatedTabStyle(tab.Name))
	}

	s += lipgloss.PlaceHorizontal(120, lipgloss.Center, styles.TitleStyle("GTop")) + styles.Cr
	tabline := lipgloss.JoinHorizontal(lipgloss.Bottom, tabs...)
	s += lipgloss.JoinVertical(lipgloss.Left, tabline, a.tabs[a.currentTab].View())

	return s
}
