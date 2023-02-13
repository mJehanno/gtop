package process

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const pstag = "pidstatus"

type Process struct {
	PID   uint64
	Ppid  uint64
	Path  string
	Name  string
	User  string
	Usage uint64
}

type ProcessStatus struct {
	Name  string `pidstatus:"Name"`
	Uid   uint64 `pidstatus:"Uid"`
	VmRes uint64 `pidstatus:"VmRSS"`
	Ppid  uint64 `pidstatus:"PPid"`
}

func GetAllProcess() ([]Process, error) {
	result := []Process{}

	matches, err := filepath.Glob("/proc/*/exe")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, file := range matches {
		var p Process
		target, _ := os.Readlink(file)
		splittedPath := strings.FieldsFunc(file, func(r rune) bool {
			return r == '/'
		})

		if splittedPath[1] == "self" || splittedPath[1] == "thread-self" {
			continue
		}

		pid, err := strconv.ParseUint(splittedPath[1], 10, 64)
		if err != nil {
			return nil, err
		}

		filePath := getStatusFilePath(splittedPath[0], splittedPath[1])
		statusFile, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		var ps ProcessStatus

		err = ps.Unmarshal(statusFile)
		if err != nil {
			return nil, err

		}
		processUser, err := user.LookupId(strconv.FormatUint(ps.Uid, 10))
		if err != nil {
			return nil, err
		}
		p = Process{
			PID:   pid,
			Name:  ps.Name,
			User:  processUser.Username,
			Usage: ps.VmRes,
			Ppid:  ps.Ppid,
		}

		if len(target) > 0 {
			p.Path = fmt.Sprintf("%+v", target)
		} else {
			p.Path = ""
		}

		result = append(result, p)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Usage > result[j].Usage
	})

	return result, nil
}

func getStatusFilePath(root, pid string) string {
	return path.Join("/", root, pid, "status")
}

func (p *ProcessStatus) Unmarshal(data []byte) error {
	stringed := string(data)
	lines := strings.Split(stringed, "\n")
	tmp := map[string]string{}

	for _, line := range lines {
		if len(line) > 0 {
			fields := strings.FieldsFunc(line, func(r rune) bool { return r == ':' })
			tmp[fields[0]] = strings.Trim(fields[1], " 	kB")
		}
	}

	reflectedValue := reflect.ValueOf(p)
	reflectedType := reflect.TypeOf(*p)

	for i := 0; i < reflectedType.NumField(); i++ {
		if tagName, ok := reflectedType.Field(i).Tag.Lookup(pstag); ok {
			if v, ok := tmp[tagName]; ok {
				switch reflectedType.Field(i).Type.Name() {
				case "uint64":
					splitString := strings.Fields(v)
					value, err := strconv.ParseUint(splitString[0], 10, 64)
					if err != nil {
						return err
					}
					reflectedValue.Elem().Field(i).SetUint(value)
				default:
					reflectedValue.Elem().Field(i).Set(reflect.ValueOf(v))
				}
			}
		}
	}

	return nil
}

func TreeMode(processes []Process) map[Process][]Process {
	m := map[Process][]Process{}
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].PID > processes[j].PID
	})
	for _, p := range processes {
		if p.Ppid != 0 {
			index := slices.IndexFunc(processes, func(pr Process) bool {
				return pr.PID == p.Ppid
			})

			if index != -1 {
				m[processes[index]] = append(m[processes[index]], p)
			}
		}
	}

	return m
}
