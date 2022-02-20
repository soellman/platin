package main

import (
	"fmt"

	"github.com/soellman/platin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current status of Platin hub",
	Long: `Show current status of Platin hub, including power, volume, and sources.
	Currently selected source is indicated by an asterisk.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showPower(quiet); err != nil {
			fmt.Printf("error showing power: %v\n", err)
		}
		if err := showVolume(quiet); err != nil {
			fmt.Printf("error showing volume: %v\n", err)
		}
		if err := showSources(quiet); err != nil {
			fmt.Printf("error showing sources: %v\n", err)
		}
	},
}

func showPower(quiet bool) error {
	if !quiet {
		fmt.Printf("Power: ")
	}

	hub := platin.NewHub(address)
	on, err := hub.Power()
	if err != nil {
		return err
	}
	str := "off"
	if on {
		str = "on"
	}

	fmt.Println(str)
	return nil
}

func showVolume(quiet bool) error {
	if !quiet {
		fmt.Printf("Volume: ")
	}

	hub := platin.NewHub(address)
	vol, err := hub.Volume()
	if err != nil {
		return err
	}

	fmt.Println(vol)
	return nil
}

func showSources(quiet bool) error {
	if !quiet {
		fmt.Println("Sources:")
	}

	hub := platin.NewHub(address)
	sources, err := hub.Sources()
	if err != nil {
		return err
	}

	for _, source := range sources {
		if !quiet {
			listMarker := "-"
			if source.Active {
				listMarker = "*"
			}
			fmt.Printf(" %s ", listMarker)
		}
		fmt.Println(source.Name)
	}
	return nil
}
