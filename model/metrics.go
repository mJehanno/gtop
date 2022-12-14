package model

import (
	"log"
	"syscall"
)

type MetricModel struct {
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

	m.TotalRam = sysInf.Totalram
	m.FreeRam = sysInf.Freeram

	return m
}
