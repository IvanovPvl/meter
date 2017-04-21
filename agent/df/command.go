package df

import "os/exec"

type DfResult struct {
	Filesystem string
	Blocks     uint
	Used       uint
	Available  uint
}

func Exec() (string, error) {
	out, err := exec.Command("df").Output()
	return string(out), err
}
