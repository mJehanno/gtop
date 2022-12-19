package memory

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

const memTag = "mem"

type Memory struct {
	MemTotal     uint64 `mem:"MemTotal"`
	MemAvailable uint64 `mem:"MemAvailable"`
	SwapTotal    uint64 `mem:"SwapTotal"`
	SwapFree     uint64 `mem:"SwapFree"`
}

func New() (*Memory, error) {
	buf, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	mem := &Memory{}

	err = mem.UnMarshal(buf)
	if err != nil {
		return nil, err
	}

	return mem, nil
}

func (m *Memory) UnMarshal(data []byte) error {
	unMarshaledMap := map[string]uint64{}

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
			unMarshaledMap[fields[0]] = value
		}
	}

	v := reflect.ValueOf(m)
	mutable := reflect.Indirect(v)
	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup(memTag); ok {
			mutable.FieldByName(value).SetUint(unMarshaledMap[value] * 1024)
		}
	}

	return nil
}
