package tabs

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/evertras/bubble-table/table"
	"github.com/mjehanno/gtop/model/cmds"
	"github.com/mjehanno/gtop/model/metrics/linux/process"
	"github.com/mjehanno/gtop/model/styles"
)

type ProcessManagerModel struct {
	processes []process.Process
	table     table.Model
	err       error
}

const (
	pid         = "pid"
	user        = "user"
	processName = "processName"
	processPath = "process"
	ramUsage    = "ramUsage"
)

func NewProcessManagerModel() *ProcessManagerModel {
	ps, err := process.GetAllProcess()
	model := ProcessManagerModel{
		processes: ps,
		err:       err,
	}

	columns := []table.Column{
		table.NewColumn(pid, "PID", 6),
		table.NewColumn(user, "User", 10),
		table.NewColumn(processName, "Process Name", 15),
		table.NewColumn(processPath, "Process", 50),
		table.NewColumn(ramUsage, "Ram Usage", 20),
	}

	tab := table.New(columns)
	tab = model.updateTable(&tab)
	selectedStyle := lipgloss.NewStyle().Background(styles.RetroBlue).Bold(true).Foreground(styles.Yellow)
	headerStyle := lipgloss.NewStyle().Foreground(styles.White)
	baseStyle := lipgloss.NewStyle().Foreground(styles.BasicGrey)

	model.table = tab.WithPageSize(25).Focused(true).BorderRounded().HighlightStyle(selectedStyle).WithPaginationWrapping(false).HeaderStyle(headerStyle).WithBaseStyle(baseStyle)

	return &model
}

func (p *ProcessManagerModel) Init() tea.Cmd {
	return nil
}

func (p *ProcessManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.TickMsg:
		var cmd tea.Cmd
		ps, err := process.GetAllProcess()
		if err != nil {
			p.err = err
		}
		p.processes = ps
		p.table = p.updateTable(&p.table)
		_, cmd = p.table.Update(msg)
		return p, tea.Batch(cmd, SyncedTick)
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	var cmd tea.Cmd
	p.table, cmd = p.table.Update(msg)
	return p, cmd
}

func (p *ProcessManagerModel) View() string {
	s := ""
	if p.err != nil {
		s = styles.ErrorStyle("can't display process list : " + styles.Cr + p.err.Error())
	} else {
		s += lipgloss.JoinVertical(lipgloss.Center, p.table.View())
	}

	return s
}

func (p *ProcessManagerModel) updateTable(tab *table.Model) table.Model {
	rows := []table.Row{}

	for _, p := range p.processes {
		rows = append(rows, table.NewRow(table.RowData{pid: strconv.Itoa(p.PID), user: p.User, processName: p.Name, processPath: p.Path, ramUsage: humanize.Bytes(p.Usage * 1024)}))
	}

	return tab.WithRows(rows)
}
