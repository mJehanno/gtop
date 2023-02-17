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
	Name          string
	width, height int
}

func (t *Tab) SetSize(w, h int) {
	t.width, t.height = w, h

	switch t.Model.(type) {
	case *ProcessManagerModel:
		m := t.Model.(*ProcessManagerModel)
		m.SetSize(w, h)
	}
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

	um, _ := NewUsageModel()
	ut := Tab{
		Model: um,
		Name:  "Usage",
	}

	// nm, _ := NewNetworkTabModel()
	// nt := Tab{
	// 	Model: nm,
	// 	Name:  "Network",
	// }

	tabs := []Tab{
		bt,
		ut,
		prt,
		//nt,
	}

	return tabs
}
