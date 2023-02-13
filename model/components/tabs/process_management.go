package tabs

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/evertras/bubble-table/table"
	"github.com/mjehanno/gtop/model/cmds"
	"github.com/mjehanno/gtop/model/components/modal"
	"github.com/mjehanno/gtop/model/stats/metrics/process"
	"github.com/mjehanno/gtop/model/styles"
)

type ProcessManagerModel struct {
	processes        []process.Process
	table            table.Model
	err              error
	treeMode         bool
	isSignalListOpen bool
	signalModal      *modal.SignalListModel
}

const (
	pid         = "pid"
	user        = "user"
	processName = "processName"
	processPath = "process"
	ramUsage    = "ramUsage"
)

var selectedStyle = lipgloss.NewStyle().Background(styles.RetroBlue).Bold(true).Foreground(styles.Yellow)
var headerStyle = lipgloss.NewStyle().Foreground(styles.White)
var baseStyle = lipgloss.NewStyle().Foreground(styles.BasicGrey)

func NewProcessManagerModel() *ProcessManagerModel {
	ps, err := process.GetAllProcess()
	model := ProcessManagerModel{
		processes:        ps,
		err:              err,
		treeMode:         false,
		isSignalListOpen: false,
	}

	columns := []table.Column{
		table.NewColumn(pid, "PID", 6).WithFiltered(true),
		table.NewColumn(user, "User", 10).WithFiltered(true),
		table.NewColumn(processName, "Process Name", 15).WithFiltered(true),
		table.NewColumn(processPath, "Process", 50).WithFiltered(true),
		table.NewColumn(ramUsage, "Ram Usage", 20),
	}

	tab := table.New(columns)
	tab = model.updateTable(&tab)

	model.table = tab.WithPageSize(25).Focused(true).BorderRounded().HighlightStyle(selectedStyle).WithPaginationWrapping(false).HeaderStyle(headerStyle).WithBaseStyle(baseStyle).Filtered(true)

	return &model
}

func (p *ProcessManagerModel) Init() tea.Cmd {
	return nil
}

func (p *ProcessManagerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.table = p.table.WithMaxTotalWidth(2 * msg.Width / 3)
		if p.isSignalListOpen {
			p.signalModal.SetSize(msg.Width/3, msg.Height)
			m, cmd := p.signalModal.List.Update(msg)
			p.signalModal.List = m
			return p, cmd
		}
		return p, nil
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
	case cmds.LeaveOverlayMsg:
		if p.isSignalListOpen {
			p.isSignalListOpen = false
		}
	case tea.KeyMsg:
		if p.isSignalListOpen {
			m, cmd := p.signalModal.Update(msg)
			p.signalModal = m.(*modal.SignalListModel)
			return p, cmd
		}
		switch msg.String() {
		case "f5":
			var cmd tea.Cmd
			p.treeMode = !p.treeMode
			p.table = p.updateTable(&p.table)
			_, cmd = p.table.Update(msg)
			return p, cmd
		case "f9":
			p.isSignalListOpen = true
			p.signalModal = modal.NewSignalListModel()
			//pid := p.table.HighlightedRow().Data["pid"]
			//realpid, _ := strconv.ParseUint(pid.(string), 10, 64)
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	p.table, cmd = p.table.Update(msg)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *ProcessManagerModel) View() string {
	s := ""

	if p.isSignalListOpen {
		s += lipgloss.JoinHorizontal(lipgloss.Left, styles.ModalStyle.Render(p.signalModal.View()), p.table.View())
	} else {
		if p.err != nil {
			s = styles.ErrorStyle("can't display process list : " + styles.Cr + p.err.Error())
		} else {
			s += lipgloss.JoinHorizontal(lipgloss.Center, p.table.View())
		}
	}

	return s
}

func (p *ProcessManagerModel) updateTable(tab *table.Model) table.Model {
	rows := []table.Row{}
	if !p.treeMode {
		for _, p := range p.processes {
			rows = append(rows, table.NewRow(table.RowData{pid: strconv.FormatUint(p.PID, 10), user: p.User, processName: p.Name, processPath: p.Path, ramUsage: humanize.Bytes(p.Usage * 1024)}))
		}
	} else {
		for k, v := range process.TreeMode(p.processes) {
			processesPath := k.Path + styles.Cr
			for i, p := range v {
				processesPath += styles.TabSep + "|-" + p.Path
				if i < len(v)-1 {
					processesPath += styles.Cr
				}
			}
			rows = append(rows, table.NewRow(table.RowData{pid: strconv.FormatUint(k.PID, 10), user: k.User, processName: k.Name, processPath: processesPath, ramUsage: humanize.Bytes(k.Usage * 1024)}))
		}
	}

	return tab.WithRows(rows)
}
