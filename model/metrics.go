package model

import (
	"log"
	"syscall"
)

type MetricModel struct {
	Uptime      int64
	FreeRam     uint64
	BufferedRam uint64
	TotalRam    uint64
	TotalSwap   uint64
	FreeSwap    uint64
}

func New() *MetricModel {
	sysInf := &syscall.Sysinfo_t{}

	err := syscall.Sysinfo(sysInf)
	if err != nil {
		log.Fatal(err)
	}
	m := new(MetricModel)

	m.Uptime = sysInf.Uptime
	m.TotalRam = sysInf.Totalram
	m.FreeRam = sysInf.Freeram
	m.TotalSwap = sysInf.Totalswap
	m.FreeSwap = sysInf.Freeswap
	m.BufferedRam = sysInf.Bufferram

	return m
}
