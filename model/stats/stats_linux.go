package stats

import (
	"github.com/mjehanno/gtop/model/stats/metrics"
	"github.com/mjehanno/gtop/model/stats/system"
	"github.com/mjehanno/gtop/model/stats/system/network"
)

type Stats struct {
	*metrics.Metrics
	*system.SystemInfo
	NetInterfaces []network.Interface
}

func New() (*Stats, error) {
	s := &Stats{}

	s.Metrics, _ = metrics.New()
	s.SystemInfo, _ = system.New()
	s.NetInterfaces = network.GetInterfaces()

	return s, nil
}
