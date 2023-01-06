package tabs

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/cmds"
	"github.com/mjehanno/gtop/model/metrics/linux/process"
	"github.com/mjehanno/gtop/model/styles"
)

type ProcessManagerModel struct {
	processes []process.Process
	table     table.Model
	err       error
}

func NewProcessManagerModel() *ProcessManagerModel {
	ps, err := process.GetAllProcess()
	model := ProcessManagerModel{
		processes: ps,
		err:       err,
	}

	columns := []table.Column{
		{Title: "PID", Width: 6},
		{Title: "User", Width: 10},
		{Title: "Process Name", Width: 15},
		{Title: "Process", Width: 50},
		{Title: "Ram Usage", Width: 20},
	}

	style := table.DefaultStyles()
	style.Header = style.Header.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Bold(true)
	style.Cell = style.Cell.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Margin(0)

	tab := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithStyles(style),
	)

	model.updateTable(&tab)

	model.table = tab

	return &model
}

func (p *ProcessManagerModel) Init() tea.Cmd {
	return nil
}

func (p *ProcessManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case cmds.TickMsg:
		var cmd tea.Cmd
		ps, err := process.GetAllProcess()
		if err != nil {
			p.err = err
		}
		p.processes = ps
		p.updateTable(&p.table)
		p.table, cmd = p.table.Update(msg)
		return p, tea.Batch(cmd, SyncedTick)
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

func (p *ProcessManagerModel) updateTable(tab *table.Model) {
	rows := []table.Row{}

	for _, p := range p.processes {
		rows = append(rows, table.Row{strconv.Itoa(p.PID), p.User, p.Name, p.Path, humanize.Bytes(p.Usage * 1024)})
	}

	tab.SetRows(rows)
}
