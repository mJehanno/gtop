package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mjehanno/gtop/model/styles"
	"github.com/mjehanno/gtop/model/tabs"
)

type AppModel struct {
	currentTab tabs.TabsEnum
	tabs       []tabs.Tab
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "tab":
			a.currentTab++
			if int(a.currentTab) >= len(a.tabs) {
				a.currentTab = 0
			}
			return a, a.tabs[a.currentTab].Init()
		case "shift+tab":
			a.currentTab--
			if int(a.currentTab) < 0 {
				a.currentTab = tabs.TabsEnum(len(a.tabs) - 1)
			}
			return a, a.tabs[a.currentTab].Init()
		}
	}

	m, cmd := a.tabs[a.currentTab].Update(msg)
	a.tabs[a.currentTab].Model = m
	return a, cmd
}

func (a *AppModel) View() string {
	s := ""

	tabLine := []string{}

	for i, tab := range a.tabs {
		if i == int(a.currentTab) {
			tabLine = append(tabLine, styles.ActivatedTabStyle(tab.Name))
			continue
		}
		tabLine = append(tabLine, styles.DeactivatedTabStyle(tab.Name))
	}

	s += lipgloss.PlaceHorizontal(120, lipgloss.Center, styles.TitleStyle("GTop")) + styles.Cr
	s += lipgloss.JoinHorizontal(lipgloss.Bottom, tabLine...) + styles.Cr
	s += lipgloss.PlaceHorizontal(240, lipgloss.Left, a.tabs[a.currentTab].View())

	return s
}
