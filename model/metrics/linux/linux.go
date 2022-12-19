package linux

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/mjehanno/gtop/model/user"
)

const tagName = "mem"

type LinuxMetric struct {
	MemTotal     uint64 `mem:"MemTotal"`
	MemAvailable uint64 `mem:"MemAvailable"`
	SwapTotal    uint64 `mem:"SwapTotal"`
	SwapFree     uint64 `mem:"SwapFree"`
}

func (l *LinuxMetric) UnMarshal(data []byte) error {
	m := map[string]uint64{}

	stringed := string(data)

	pured := strings.ReplaceAll(stringed, ":", "")

	lines := strings.Split(pured, "\n")

	for _, v := range lines {
		fields := strings.Fields(v)
		if len(fields) > 1 {

			value, err := strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return err
			}
			m[fields[0]] = value
		}
	}

	v := reflect.ValueOf(l)
	mutable := reflect.Indirect(v)
	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup(tagName); ok {
			mutable.FieldByName(value).SetUint(m[value] * 1024)
		}
	}

	return nil
}

func New() (*LinuxMetric, error) {
	lm := new(LinuxMetric)
	err := lm.getMetrics()
	if err != nil {
		return nil, err
	}

	return lm, nil
}

func (l *LinuxMetric) getMetrics() error {
	buf, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return err
	}

	err = l.UnMarshal(buf)
	if err != nil {
		return err
	}

	return nil
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

func (l *LinuxMetric) GetCpuLoad() {}
