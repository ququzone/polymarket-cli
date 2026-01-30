package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "An example subcommand",
	Long:  `A longer description of the example command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Example command called!")
		if len(args) > 0 {
			fmt.Printf("Arguments: %v\n", args)
		}
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)

	exampleCmd.Flags().StringP("name", "n", "", "Name to greet")
	exampleCmd.Flags().CountP("verbose", "v", "Verbose output")
}
