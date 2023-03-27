package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	source, destination string
	decorName           string
	packageName         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "decorgen",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&source, "source", "", "interface file path")
	rootCmd.MarkFlagRequired("source")

	rootCmd.PersistentFlags().StringVar(&destination, "destination", "", "decorator file path")
	rootCmd.MarkFlagRequired("destination")

	rootCmd.PersistentFlags().StringVar(&decorName, "decor-name", "", "decorator struct name")
	rootCmd.MarkFlagRequired("decor-name")

	rootCmd.PersistentFlags().StringVar(&packageName, "package", "", "package name")
	rootCmd.MarkFlagRequired("package")
}
