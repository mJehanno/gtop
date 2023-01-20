package cpu

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const cpuTag = "cpu"

type CPUInfo struct {
	ModelName   string  `cpu:"model name"`
	Mhz         float64 `cpu:"cpu MHz"`
	CacheSize   int     `cpu:"cache size"`
	PhysicalId  int     `cpu:"physical id"`
	CoreId      int     `cpu:"core id"`
	ProcessorId int     `cpu:"processor"`
}

func (c *CPUInfo) UnMarshal(data []byte) error {
	unMarshaledMap := map[string]any{}
	stringed := string(data)

	lines := strings.Split(stringed, "\n")

	v := reflect.ValueOf(c)
	t := reflect.TypeOf(*c)

	for _, line := range lines {
		fields := strings.Split(line, ":")

		for i, f := range fields {
			fields[i] = strings.Trim(f, "	 KB")
		}

		if len(fields) > 1 {
			unMarshaledMap[fields[0]] = fields[1]
		}
	}

	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup(cpuTag); ok {
			var err error
			switch t.Field(i).Type.Name() {
			case "int":
				s, ok := unMarshaledMap[value].(string)
				if !ok {
					log.Fatalf("name: %s, value: %s, err : %s \n", value, unMarshaledMap[value], err)
				}
				unMarshaledMap[value], err = strconv.Atoi(s)
				if err != nil {
					log.Fatalf("name: %s, value: %s, err : %s \n", value, unMarshaledMap[value], err)
				}
				if value == "cache size" {
					unMarshaledMap[value] = unMarshaledMap[value].(int) * 1024
				}
			case "float64":
				unMarshaledMap[value], err = strconv.ParseFloat(unMarshaledMap[value].(string), 64)
				if err != nil {
					log.Fatal(err)
				}
			}

			newValue := reflect.ValueOf(unMarshaledMap[value])
			v.Elem().Field(i).Set(newValue)
		}
	}

	return nil
}

func getCpuInfos() ([]CPUInfo, error) {
	file, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}

	matrix := [][]string{}

	lines := strings.Split(string(file), "\n")

	tmp := []string{}

	for _, line := range lines {
		if len(line) == 0 {
			if len(tmp) > 1 {
				matrix = append(matrix, tmp)
			}
			tmp = []string{}
		}
		tmp = append(tmp, line)
	}

	result := make([]CPUInfo, len(matrix))
	for i, arr := range matrix {
		block := strings.Join(arr, "\n")
		var c CPUInfo
		err = c.UnMarshal([]byte(block))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result[i] = c
	}
	return result, nil
}

func getCpuStats() ([]CPUStat, error) {
	file, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	result := make([]CPUStat, 0)
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)

		if !strings.Contains(fields[0], "cpu") {
			break
		}

		values := make([]uint64, len(fields)-1)

		for j, field := range fields {
			if j > 0 {
				value, err := strconv.ParseUint(field, 10, 64)
				if err != nil {
					return nil, err
				}
				values[j-1] = value
			}
		}

		var stat CPUStat
		reflectedValue := reflect.ValueOf(&stat)

		for i, v := range values {
			reflectedValue.Elem().Field(i).SetUint(v)
		}

		result = append(result, stat)
	}

	return result, nil
}

func getCpu() ([]CPU, error) {
	cpuInfos, err := getCpuInfos()
	if err != nil {
		return nil, err
	}

	cpuStats, err := getCpuStats()
	if err != nil {
		return nil, err
	}

	result := make([]CPU, len(cpuInfos))

	for i := range cpuInfos {
		result[i].CPUInfo = &cpuInfos[i]
	}

	for i := range cpuStats {
		result[i].CPUStat = &cpuStats[i]
	}

	return result, nil
}

type CPUStat struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
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
	return c.Idle
}

type CPU struct {
	*CPUInfo
	*CPUStat
}

func New() ([]CPU, error) {
	cpu, err := getCpu()
	if err != nil {
		return nil, fmt.Errorf("error while getting cpu data in constructor : %w", err)
	}

	return cpu, nil
}
