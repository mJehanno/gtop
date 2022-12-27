package tabs

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	columns := []table.Column{
		{Title: "PID", Width: 6},
		{Title: "User", Width: 10},
		{Title: "Process Name", Width: 15},
		{Title: "Process", Width: 50},
	}
	rows := []table.Row{}

	for _, p := range ps {
		rows = append(rows, table.Row{strconv.Itoa(p.PID), p.User, p.Name, p.Path})
	}

	style := table.DefaultStyles()
	style.Header = style.Header.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Bold(true)
	style.Cell = style.Cell.BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")).Margin(0)

	tab := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithStyles(style),
	)

	return &ProcessManagerModel{
		processes: ps,
		err:       err,
		table:     tab,
	}
}

func (p *ProcessManagerModel) Init() tea.Cmd {
	return nil
}

func (p *ProcessManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
