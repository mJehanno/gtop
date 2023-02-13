package tabs

import (
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/stats"
	"github.com/mjehanno/gtop/model/styles"
)

type BasicInformationModel struct {
	stats stats.Stats
}

func NewBasincInformationModel() *BasicInformationModel {
	stat, _ := stats.New()
	return &BasicInformationModel{
		stats: *stat,
	}
}

func (b *BasicInformationModel) Init() tea.Cmd {
	return nil
}

func (b *BasicInformationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *BasicInformationModel) View() string {
	netAddresses := []string{}

	pattern := regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
	for _, a := range b.stats.NetInterfaces {
		if pattern.MatchString(a.IpAddress) && a.IpAddress != "127.0.0.1" {
			netAddresses = append(netAddresses, a.String())
		}
	}

	hostname := b.stats.Hostname

	var stringedUptime string

	uptime := b.stats.Uptime
	stringedUptime = humanize.Time(time.Now().Add(-time.Second * time.Duration(uptime)))

	distributionFullName := b.stats.DistribName + styles.SpaceSep + b.stats.DistribVersion
	distrib := styles.LabelStyleRender("Distribution:") + distributionFullName + styles.TabSep + styles.LabelStyleRender("Kernel:") + styles.SpaceSep + b.stats.KernelVersion

	userLine := styles.LabelStyleRender("Current User:") + styles.SpaceSep + b.stats.User.Uid + styles.SpaceSep + b.stats.User.Username + "@" + hostname + styles.TabSep + styles.LabelStyleRender("Groups:") + styles.SpaceSep + strings.Join(b.stats.User.Groups, ", ") + styles.Cr
	systemLine := styles.LabelStyleRender("Uptime:") + styles.SpaceSep + stringedUptime + styles.TabSep + styles.LabelStyleRender("Network:") + styles.SpaceSep + strings.Join(netAddresses, ", ") + styles.Cr

	ramTitleLine := styles.LabelStyleRender("RAM:") + styles.Cr
	ramLine := "Free:" + styles.SpaceSep + humanize.Bytes(b.stats.GetAvailableRam()) + styles.TabSep + "Total:" + styles.SpaceSep + humanize.Bytes(b.stats.GetTotalRam())
	ramBloc := lipgloss.JoinVertical(lipgloss.Left, ramTitleLine, ramLine)

	swapTitleLine := styles.LabelStyleRender("Swap:") + styles.Cr
	swapLine := "Free:" + styles.SpaceSep + humanize.Bytes(b.stats.GetAvailableSwap()) + styles.TabSep + "Total:" + styles.SpaceSep + humanize.Bytes(b.stats.GetTotalSwap())
	swapBloc := lipgloss.JoinVertical(lipgloss.Left, swapTitleLine, swapLine)

	textBlock := lipgloss.JoinVertical(lipgloss.Left, distrib, userLine, systemLine, ramBloc, swapBloc)

	return textBlock
}
