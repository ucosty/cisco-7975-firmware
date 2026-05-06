package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ucosty/cisco-7975-firmware/archive-tools/pkg/cnu"
)

var listCmd = &cobra.Command{
	Use:   "list <archive filename>",
	Args:  cobra.ExactArgs(1),
	Short: "List contents of firmware file",
	Run: func(cmd *cobra.Command, args []string) {
		err := cnu.List(args[0])
		if err != nil {
			fmt.Printf("Failed to list firmware contents: %v\n", err)
			return
		}
	},
}

var unpackCmd = &cobra.Command{
	Use:   "unpack <archive filename> <output directory>",
	Args:  cobra.ExactArgs(2),
	Short: "Unpack firmware file",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Printf("Failed to get verbose flag: %v\n", err)
			return
		}

		err = cnu.Unpack(args[0], args[1], !verbose)
		if err != nil {
			fmt.Printf("Failed to unpack firmware: %v\n", err)
			return
		}
	},
}

var packCmd = &cobra.Command{
	Use:   "pack <archive filename> <input directory>",
	Args:  cobra.ExactArgs(2),
	Short: "Pack firmware file",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Printf("Failed to get verbose flag: %v\n", err)
			return
		}

		err = cnu.Pack(args[1], args[0], !verbose)
		if err != nil {
			fmt.Printf("Failed to pack firmware: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(unpackCmd)
	unpackCmd.Flags().BoolP("verbose", "v", false, "verbose output")

	rootCmd.AddCommand(packCmd)
	packCmd.Flags().BoolP("verbose", "v", false, "verbose output")
}
