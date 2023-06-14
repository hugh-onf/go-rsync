package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hugh-onf/go-rsync/internal/os_ext"
	"github.com/schollz/progressbar/v3"
)

// Hard code to rsync, for now
var copySubCmd string = "rsync"
var defaultProgressInterval int64 = 1000

// Exec executes the copy sub command using the passed in args
func Exec(args ...string) (*int64, string, error) {
	// Check if copy command available
	if !os_ext.IsCommandAvailable(copySubCmd) {
		return nil, copySubCmd, fmt.Errorf("'%s' command not found, need to install it first", copySubCmd)
	}

	// Hard code to rsync
	// so 2nd arg is the source folder
	// and 3rd arg is the dest folder
	if len(args) != 3 {
		return nil, copySubCmd, fmt.Errorf("invalid args, should be `-ah SOURCE_DIR TARGET_DIR`, we expect at least one option flag")
	}
	copyFrom := args[1]
	copyTo := args[2]

	// Start the progress bar
	totalSize, err := os_ext.GetFolderSize(copyFrom)
	if err != nil {
		return nil, copySubCmd, fmt.Errorf("cannot get directory size of source path '%s'", copyFrom)
	}
	// Partial copy may happen
	partialSize, err := os_ext.GetFolderSize(copyTo)
	copySize := totalSize
	if err == nil && partialSize > 0 {
		copySize = totalSize - partialSize
		if copySize < 0 {
			copySize = 0
		}
	}
	bar := progressbar.NewOptions64(copySize,
		progressbar.OptionShowBytes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(true),
	)
	defer bar.Exit()

	// Pass all the args to copy sub command
	cmd := exec.Command(copySubCmd, args...)

	// Queue it
	err = cmd.Start()
	if err != nil {
		return &copySize, copySubCmd, err
	}

	go func() {
		// Progess checking in background
		progressIntervalMs, err := strconv.ParseInt(os.Getenv("GO_RSYNC_PROGRESS_INTERVAL"), 10, 0)
		if err != nil {
			progressIntervalMs = defaultProgressInterval
		}
		// Loop until the end, i.e. until the program exit
		for 1 > 0 {
			time.Sleep(time.Duration(progressIntervalMs) * time.Millisecond)
			targetSize, _ := os_ext.GetFolderSize(copyTo)
			// Substract partial copy
			targetSize = targetSize - partialSize
			if targetSize > copySize {
				targetSize = copySize
			}
			bar.Set64(targetSize)
		}
	}()

	// Run it and wait for it to finish
	err = cmd.Wait()

	// Fill the bar
	if err == nil {
		bar.Set64(copySize)
	}

	return &copySize, copySubCmd, err
}
