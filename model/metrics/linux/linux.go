package linux

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mjehanno/gtop/model/user"
)

type LinuxMetric struct{}

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
	buf, err := ioutil.ReadFile("/proc/uptime")
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
	return 0
}

func (l *LinuxMetric) GetAvailableRam() uint64 {
	return 0
}

func (l *LinuxMetric) GetTotalSwap() uint64 {
	return 0
}

func (l *LinuxMetric) GetAvailableSwap() uint64 {
	return 0
}

func (l *LinuxMetric) GetCpuLoad() {}
