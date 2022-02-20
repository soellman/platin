package main

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	address string
	quiet   bool

	rootCmd = &cobra.Command{
		Use:   "platin",
		Short: "Control your Platin hub",
		Long:  `platin is a simple commandline utility to control your Platin hub.`,
	}
)

func main() {
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "minimal output")
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "Platin hub address (required)")
	rootCmd.MarkPersistentFlagRequired("address")
	rootCmd.Execute()
}
