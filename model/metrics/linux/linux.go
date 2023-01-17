package linux

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mjehanno/gtop/model/metrics/linux/cpu"
	"github.com/mjehanno/gtop/model/metrics/linux/memory"
	"github.com/mjehanno/gtop/model/metrics/linux/system"
	"github.com/mjehanno/gtop/model/user"
)

type LinuxMetric struct {
	*memory.Memory
	CPUs    []cpu.CPU
	SysInfo *system.SystemInfo
}

func New() (*LinuxMetric, error) {
	lm := new(LinuxMetric)
	mem, err := memory.New()
	if err != nil {
		return nil, err
	}
	lm.Memory = mem

	c, err := cpu.New()
	if err != nil {
		return nil, err
	}
	lm.CPUs = c

	s, err := system.New()
	if err != nil {
		return nil, err
	}
	lm.SysInfo = s

	return lm, nil
}

func (l *LinuxMetric) GetDistribution() string {
	return fmt.Sprintf("%s %s", l.SysInfo.DistribName, l.SysInfo.DistribVersion)

}

// GetHostname return the current host hostname.
func (l *LinuxMetric) GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// GetCurrentUser return the current user with it's UID, Username and Groups.
func (l *LinuxMetric) GetCurrentUser() *user.User {
	currentUser := user.New()
	return currentUser
}

func (l *LinuxMetric) GetUptime() (float64, error) {
	buf, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(buf))

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return uptime, nil
}

func (l *LinuxMetric) GetTotalRam() uint64 {
	return l.MemTotal
}

func (l *LinuxMetric) GetAvailableRam() uint64 {
	return l.MemAvailable
}

func (l *LinuxMetric) GetTotalSwap() uint64 {
	return l.SwapTotal
}

func (l *LinuxMetric) GetAvailableSwap() uint64 {
	return l.SwapFree
}

func (l *LinuxMetric) GetCpuLoad() []cpu.CPU {
	return l.CPUs
}
