package metrics

import (
	"github.com/mjehanno/gtop/model/metrics/linux/cpu"
	"github.com/mjehanno/gtop/model/user"
)

type Metric interface {
	GetHostname() (string, error)
	GetCurrentUser() *user.User
	GetUptime() (float64, error)
	GetTotalRam() uint64
	GetAvailableRam() uint64
	GetTotalSwap() uint64
	GetAvailableSwap() uint64
	GetCpuLoad() []cpu.CPU
}
