package model

import (
	"log"
	"syscall"
)

type MetricModel struct {
	Uptime   int64
	FreeRam  uint64
	TotalRam uint64
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

	return m
}
