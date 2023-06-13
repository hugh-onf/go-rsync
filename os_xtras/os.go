// Abstractions on OS calls, so we can swap implementations if need
package os_xtras

import (
	"os/exec"

	getFolderSize "github.com/markthree/go-get-folder-size/src"
)

// GetFolderSize returns the size of the target folder in bytes
func GetFolderSize(path string) (int64, error) {
	size := getFolderSize.LooseParallel(path)
	return size, nil
}

// IsCommandAvailable checks if the command is available in the system
func IsCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
