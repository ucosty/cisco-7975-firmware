package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: os.Args[0],
}

func Execute() error {
	return rootCmd.Execute()
}
