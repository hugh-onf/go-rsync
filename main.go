package main

import (
	"fmt"

	"github.com/hugh-onf/go-rsync/pkg/cmd"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	err, size, cmd := cmd.Exec()
	if err != nil {
		fmt.Printf("\ncommand finished with error, %v. Check '%s' manual for help.\n", err, cmd)
	} else {
		p := message.NewPrinter(language.English)
		p.Printf("\ncommand finished successfully, total copy size: %d bytes.\n", *size)
	}
}
