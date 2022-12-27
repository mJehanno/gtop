package tabs

import tea "github.com/charmbracelet/bubbletea"

type TabsEnum int

const (
	BasicInformation TabsEnum = iota
	ProcessManager
)

type Tab struct {
	tea.Model
	Name string
}

func GetAllTabs() []Tab {
	bi := NewBasincInformationModel()
	bt := Tab{
		Model: bi,
		Name:  "Basic Information",
	}

	pmm := NewProcessManagerModel()
	prt := Tab{
		Model: pmm,
		Name:  "Process Management",
	}

	tabs := []Tab{
		bt,
		prt,
	}

	return tabs
}
