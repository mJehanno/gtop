package metrics

import "github.com/mjehanno/gtop/model/user"

type Metric interface {
	GetHostname() (string, error)
	GetCurrentUser() *user.User
	GetUptime() int64
	GetTotalRam() uint64
	GetAvailableRam() uint64
	GetTotalSwap() uint64
	GetAvailableSwap() uint64
	GetCpuLoad()
}
