package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ucosty/cisco-7975-firmware/archive-tools/pkg/cnu"
)

var unsignCmd = &cobra.Command{
	Use:   "unsign <signed firmware filename> <output filename>",
	Args:  cobra.ExactArgs(2),
	Short: "Remove digital signature from firmware",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			panic(err)
		}

		parsedTokens, err := cnu.ParseSBN(data)
		if err != nil {
			fmt.Println("Error parsing SBN:", err)
			os.Exit(1)
		}

		signatureLength, err := cnu.GetSignatureLength(parsedTokens)
		if err != nil {
			fmt.Println("Error getting signature length:", err)
			os.Exit(1)
		}

		isArchive := cnu.HasArchiveHeader(data[signatureLength:])

		if !isArchive {
			fmt.Printf("File is not a CNU archive")
			os.Exit(1)
		}

		err = os.WriteFile(args[1], data[signatureLength:], 0644)
		if err != nil {
			panic(err)
		}
	},
}

var parseCmd = &cobra.Command{
	Use:   "parse-signature <signed firmware filename>",
	Args:  cobra.ExactArgs(1),
	Short: "Parse firmware digital signature",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			panic(err)
		}

		parsedTokens, err := cnu.ParseSBN(data)
		if err != nil {
			fmt.Println("Error parsing SBN:", err)
			return
		}

		cnu.PrintSBN(parsedTokens, 0)
	},
}

func init() {
	rootCmd.AddCommand(unsignCmd)
	rootCmd.AddCommand(parseCmd)
}
