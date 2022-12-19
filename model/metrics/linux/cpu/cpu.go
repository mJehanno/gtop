package cpu

import "reflect"

const cpuTag = "cpu"

type CPUInfo struct {
	ModelName   string  `cpu:"model name"`
	Mhz         float64 `cpu:"cpu MHz"`
	CacheSize   int     `cpu:"cache size"`
	PhysicalId  int     `cpu:"physical id"`
	CoreId      int     `cpu:"core id"`
	ProcessorId int     `cpu:"processor"`
}

type CPUStat struct {
	user      uint64
	nice      uint64
	system    uint64
	idle      uint64
	iowait    uint64
	irq       uint64
	softirq   uint64
	steal     uint64
	guest     uint64
	guestNice uint64
}

func (c *CPUStat) GetTotal() uint64 {
	t := reflect.TypeOf(*c)
	val := reflect.ValueOf(*c)
	sum := uint64(0)

	for i := 0; i < t.NumField(); i++ {
		sum += val.Field(i).Uint()
	}

	return sum
}

func (c *CPUStat) GetIdle() uint64 {
	return c.idle
}

type CPU struct {
	CPUInfo
	CPUStat
}
