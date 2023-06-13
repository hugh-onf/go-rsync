package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hugh-onf/go-rsync/os_xtras"
	"github.com/schollz/progressbar/v3"
)

// Hard code to rsync, for now
var COPY_SUB_COMMAND string = "rsync"
var DEFAULT_PROGRESS_INTERVAL int64 = 1000

func main() {
	// Check if copy command available
	if !os_xtras.IsCommandAvailable(COPY_SUB_COMMAND) {
		panic(fmt.Sprintf("'%s' command not found, need to install it first.", COPY_SUB_COMMAND))
	}

	// Hard code to rsync
	// so 2nd arg is the source folder
	// and 3rd arg is the dest folder
	if len(os.Args) != 4 {
		panic("invalid args, should be `-ah SOURCE_DIR TARGET_DIR`, we expect at least one option flag.")
	}
	copyFrom := os.Args[2]
	copyTo := os.Args[3]

	// Start the progress bar
	totalSize, err := os_xtras.GetFolderSize(copyFrom)
	if err != nil {
		panic(fmt.Sprintf("cannot get directory size of source path '%s'.", copyFrom))
	}
	// Partial copy may happen
	partialSize, err := os_xtras.GetFolderSize(copyTo)
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
	cmd := exec.Command(COPY_SUB_COMMAND, os.Args[1:]...)

	// Queue it
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go func() {
		// Progess checking in background
		progressIntervalMs, err := strconv.ParseInt(os.Getenv("GO_RSYNC_PROGRESS_INTERVAL"), 10, 0)
		if err != nil {
			progressIntervalMs = DEFAULT_PROGRESS_INTERVAL
		}
		// Loop until the end, i.e. until the program exit
		for 1 > 0 {
			time.Sleep(time.Duration(progressIntervalMs) * time.Millisecond)
			targetSize, _ := os_xtras.GetFolderSize(copyTo)
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

	if err != nil {
		fmt.Printf("\ncommand finished with error, %v. Check '%s' manual for help.\n", err, COPY_SUB_COMMAND)
	} else {
		bar.Set64(copySize)
		fmt.Printf("\ncommand finished successfully.\n")
	}
}
