package cmd

import "os"

// ExecCurrentContext executes the copy sub command using the same args as the main function
func ExecCurrentContext() (*int64, string, error) {
	return Exec(os.Args[1:]...)
}
