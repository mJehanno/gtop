package tabs

import (
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/mjehanno/gtop/model/metrics/linux"
	"github.com/mjehanno/gtop/model/metrics/os"
	"github.com/mjehanno/gtop/model/network"
	"github.com/mjehanno/gtop/model/styles"
)

type BasicInformationModel struct {
	OS         os.Os
	interfaces []network.Interface
}

func NewBasincInformationModel() *BasicInformationModel {
	return &BasicInformationModel{
		interfaces: network.GetInterfaces(),
	}
}

func (b *BasicInformationModel) Init() tea.Cmd {
	b.OS = *InitOsData()
	return nil
}

func (b *BasicInformationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

func (b *BasicInformationModel) View() string {
	netAddresses := make([]string, len(b.interfaces))

	for i, a := range b.interfaces {
		netAddresses[i] = a.String()
	}

	hostname, err := b.OS.Metrics.GetHostname()
	if err != nil {
		hostname = styles.ErrorStyle("- error while getting hostname")
	}

	var stringedUptime string

	uptime, err := b.OS.Metrics.GetUptime()
	if err != nil {
		stringedUptime = styles.ErrorStyle("- error while getting uptime")
	} else {
		stringedUptime = humanize.Time(time.Now().Add(-time.Second * time.Duration(uptime)))
	}

	distrib := styles.LabelStyleRender("Distribution:") + styles.SpaceSep + b.OS.Metrics.GetDistribution()

	userLine := styles.LabelStyleRender("Current User:") + styles.SpaceSep + b.OS.Metrics.GetCurrentUser().Uid + styles.SpaceSep + b.OS.Metrics.GetCurrentUser().Username + "@" + hostname + styles.TabSep + styles.LabelStyleRender("Groups:") + styles.SpaceSep + strings.Join(b.OS.Metrics.GetCurrentUser().Groups, ", ") + styles.Cr
	systemLine := styles.LabelStyleRender("Uptime:") + styles.SpaceSep + stringedUptime + styles.TabSep + styles.LabelStyleRender("Network:") + styles.SpaceSep + strings.Join(netAddresses, ", ") + styles.Cr

	textBlock := lipgloss.JoinVertical(lipgloss.Left, distrib, userLine, systemLine)

	return textBlock
}

func InitOsData() *os.Os {
	switch runtime.GOOS {
	case "linux":
		metrics, _ := linux.New()
		return &os.Os{
			Metrics: metrics,
		}

	case "windows":
	case "darwin":
	case "bsd":
	}
	return nil
}
