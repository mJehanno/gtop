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
		bi,
		"Basic Information",
	}
	tabs := []Tab{
		bt,
	}

	return tabs
}
