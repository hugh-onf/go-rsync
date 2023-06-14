package os_ext

import (
	"os/exec"
)

func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
