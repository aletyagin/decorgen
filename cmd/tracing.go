package cmd

import (
	"github.com/aletyagin/decorgen/command/tracing"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// tracingCmd represents the tracing command
var tracingCmd = &cobra.Command{
	Use:   "tracing",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := os.Create(destination)
		defer out.Close()

		err := tracing.Run(source, out, decorName, packageName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tracingCmd)
}
