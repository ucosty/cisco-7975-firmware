package main

import (
	"fmt"
	"os"

	"github.com/ucosty/cisco-7975-firmware/archive-tools/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
