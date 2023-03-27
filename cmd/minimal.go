package cmd

import (
	"decorgen/command/minimal"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// minimalCmd represents the minimal command
var minimalCmd = &cobra.Command{
	Use:   "minimal",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := os.Create(destination)
		defer out.Close()

		err := minimal.Run(source, out, decorName, packageName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(minimalCmd)
}
