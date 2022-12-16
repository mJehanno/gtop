package linux

import (
	"os"

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

func (l *LinuxMetric) GetUptime() int64 {
	return 0
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
