// Package df run `df` command.
package df

import (
	"os/exec"
	"strconv"
	"strings"
)

// Result represents the result of `df` command execution.
type Result struct {
	Filesystem string `json:"filesystem"`
	Blocks     uint64 `json:"blocks"`
	Used       uint64 `json:"used"`
	Available  uint64 `json:"available"`
	Use        uint64 `json:"use"`
	MountedOn  string `json:"mounted_on"`
}

// Exec exec `df` command.
func Exec() ([]Result, error) {
	out, err := exec.Command("df").Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")

	var results []Result
	for _, line := range lines[1 : len(lines)-1] {
		values := strings.Fields(line)
		blocks, err := strconv.ParseUint(values[1], 10, 32)
		if err != nil {
			return nil, err
		}

		used, err := strconv.ParseUint(values[2], 10, 32)
		if err != nil {
			return nil, err
		}

		available, err := strconv.ParseUint(values[3], 10, 32)
		if err != nil {
			return nil, err
		}

		use, err := strconv.ParseUint(strings.Replace(values[4], "%", "", 1), 10, 32)
		if err != nil {
			return nil, err
		}

		result := Result{
			Filesystem: values[0],
			Blocks:     blocks,
			Used:       used,
			Available:  available,
			Use:        use,
			MountedOn:  values[5],
		}

		results = append(results, result)
	}

	return results, nil
}
