package main

import (
	"errors"
	"fmt"

	"github.com/soellman/platin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sourceCmd)
	sourceCmd.AddCommand(sourceActiveCmd)
	sourceCmd.AddCommand(sourceSelectCmd)
}

var sourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Show and select sources",
	Long: `Show list of sources and select source to be active.
	Currently selected source indicated by an asterisk.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showSources(quiet); err != nil {
			fmt.Printf("error showing sources: %v\n", err)
		}
	},
}

var sourceActiveCmd = &cobra.Command{
	Use:   "active",
	Short: "Show active source",
	Long:  `Show source that is currently active.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showActiveSource(quiet); err != nil {
			fmt.Printf("error showing active source: %v\n", err)
		}
	},
}

var sourceSelectCmd = &cobra.Command{
	Use:   "select",
	Short: "Select source",
	Long:  `Select a source to become active. Requires source name argument.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please specify source name")
		} else if len(args) > 1 {
			return errors.New("provide only a single source argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := selectSource(args[0]); err != nil {
			fmt.Printf("error selecting source %q: %v\n", args[0], err)
			return
		}
		if err := showActiveSource(quiet); err != nil {
			fmt.Printf("error showing active source: %v\n", err)
		}
	},
}

func showActiveSource(quiet bool) error {
	hub := platin.NewHub(address)
	source, err := hub.ActiveSource()
	if err != nil {
		return err
	}

	fmt.Print(source.Name)
	if !quiet {
		fmt.Print(": active")
	}
	fmt.Print("\n")
	return nil
}

func selectSource(name string) error {
	hub := platin.NewHub(address)
	return hub.SetSource(name)
}
