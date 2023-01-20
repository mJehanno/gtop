package metrics

import (
	"github.com/mjehanno/gtop/model/stats/metrics/cpu"
	"github.com/mjehanno/gtop/model/stats/metrics/memory"
)

type Metrics struct {
	memory *memory.Memory
	cpus   []cpu.CPU
}

func New() (*Metrics, error) {
	m := &Metrics{}
	mem, _ := memory.New()
	m.memory = mem

	cpus, _ := cpu.New()
	m.cpus = cpus

	return m, nil
}

func (m *Metrics) GetAvailableRam() uint64 {
	return m.memory.MemAvailable
}
func (m *Metrics) GetTotalRam() uint64 {
	return m.memory.MemTotal
}
func (m *Metrics) GetAvailableSwap() uint64 {
	return m.memory.SwapFree
}
func (m *Metrics) GetTotalSwap() uint64 {
	return m.memory.SwapTotal
}

func (m *Metrics) GetAllAvailableCpu() []uint64 {
	r := make([]uint64, len(m.cpus))

	for i, cpu := range m.cpus {
		r[i] = cpu.GetIdle()
	}

	return r
}

func (m *Metrics) GetAllTotalCpu() []uint64 {
	r := make([]uint64, len(m.cpus))

	for i, cpu := range m.cpus {
		r[i] = cpu.GetTotal()
	}

	return r
}

func (m *Metrics) GetCpuLoad() []cpu.CPU {
	return m.cpus
}
