package system

import (
	"os"
	"reflect"
	"strings"
	"unicode"
)

type SystemInfo struct {
	DistribName    string
	DistribVersion string
}

func New() (*SystemInfo, error) {
	var s SystemInfo

	o, err := readOSRelease()
	if err != nil {
		return nil, err
	}

	s.DistribName = o.Name
	s.DistribVersion = o.Version

	return &s, nil
}

const osreleasetag = "osr"

type OSRelease struct {
	Name    string `osr:"NAME"`
	Version string `osr:"VERSION"`
}

func readOSRelease() (*OSRelease, error) {
	file, err := os.ReadFile("/etc/os-release")
	if err != nil && os.IsNotExist(err) {
		file, err = os.ReadFile("/usr/lib/os-release")
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	osRelease := &OSRelease{}

	err = osRelease.UnMarshal(file)
	if err != nil {
		return nil, err
	}

	return osRelease, nil
}

func (o *OSRelease) UnMarshal(data []byte) error {
	stringed := string(data)
	lines := strings.Split(stringed, "\n")
	lines = lines[:len(lines)-1]

	m := map[string]string{}

	for _, line := range lines {
		line = strings.ReplaceAll(line, "\"", "")
		kv := strings.Split(line, "=")
		m[kv[0]] = kv[1]
	}

	t := reflect.TypeOf(*o)
	mutable := reflect.ValueOf(o)

	for i := 0; i < t.NumField(); i++ {
		if tag, ok := t.Field(i).Tag.Lookup(osreleasetag); ok {
			fieldName := toUpperLower(tag)
			value := reflect.ValueOf(m[tag])

			mutable.Elem().FieldByName(fieldName).Set(value)
		}
	}

	return nil
}

func toUpperLower(s string) string {
	result := []rune{}
	for i, r := range s {
		if i != 0 {
			result = append(result, unicode.ToLower(r))
			continue
		}
		result = append(result, unicode.ToUpper(r))
	}
	return string(result)
}
