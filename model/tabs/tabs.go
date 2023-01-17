package tabs

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/gtop/model/cmds"
)

type TabsEnum int

const (
	BasicInformation TabsEnum = iota
	Usage
	ProcessManager
	Network
)

var SyncedTick tea.Cmd = cmds.TickCommand(400 * time.Millisecond)

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
