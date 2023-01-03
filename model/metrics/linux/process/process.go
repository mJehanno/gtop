package process

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const pstag = "pidstatus"

type Process struct {
	PID   int
	Path  string
	Name  string
	User  string
	Usage uint64
}

type ProcessStatus struct {
	Name  string `pidstatus:"Name"`
	Uid   uint64 `pidstatus:"Uid"`
	VmRes uint64 `pidstatus:"VmRSS"`
}

func GetAllProcess() ([]Process, error) {
	result := []Process{}

	matches, err := filepath.Glob("/proc/*/exe")
	if err != nil {
		return nil, err
	}

	for _, file := range matches {
		target, _ := os.Readlink(file)
		if len(target) > 0 {
			splittedPath := strings.FieldsFunc(file, func(r rune) bool {
				return r == '/'
			})

			if splittedPath[1] == "self" || splittedPath[1] == "thread-self" {
				continue
			}

			pid, err := strconv.Atoi(splittedPath[1])
			if err != nil {
				return nil, err
			}

			statusFile, err := os.ReadFile(getStatusFilePath(splittedPath[0], splittedPath[1]))
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

			p := Process{
				PID:   pid,
				Path:  fmt.Sprintf("%+v", target),
				Name:  ps.Name,
				User:  processUser.Username,
				Usage: ps.VmRes,
			}

			result = append(result, p)
		}
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
			fields := strings.FieldsFunc(line, func(r rune) bool {
				return r == ':'
			})
			tmp[fields[0]] = strings.Trim(fields[1], " 	kB")
		}
	}

	reflectedValue := reflect.ValueOf(p)
	reflectedType := reflect.TypeOf(*p)

	for i := 0; i < reflectedType.NumField(); i++ {
		if tagName, ok := reflectedType.Field(i).Tag.Lookup(pstag); ok {
			switch reflectedType.Field(i).Type.Name() {
			case "uint64":
				splitString := strings.Fields(tmp[tagName])

				value, err := strconv.ParseUint(splitString[0], 10, 64)
				if err != nil {
					return err
				}
				reflectedValue.Elem().Field(i).SetUint(value)
			default:
				reflectedValue.Elem().Field(i).Set(reflect.ValueOf(tmp[tagName]))
			}
		}
	}

	return nil
}
