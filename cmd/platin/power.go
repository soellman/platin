package main

import (
	"fmt"

	"github.com/soellman/platin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(powerCmd)
	powerCmd.AddCommand(powerOffCmd)
	powerCmd.AddCommand(powerOnCmd)
	powerCmd.AddCommand(powerToggleCmd)
}

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Turn speakers on/off",
	Long:  `Show current power status, toggle the power, or turn power on or off.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showPower(quiet); err != nil {
			fmt.Printf("error showing power: %v\n", err)
		}
	},
}

var powerOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn speakers off",
	Long:  `Turn off power to the speakers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setPower(false); err != nil {
			fmt.Printf("error setting power: %v\n", err)
			return
		}
		if err := showPower(quiet); err != nil {
			fmt.Printf("error showing power: %v\n", err)
		}
	},
}

var powerOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn speakers on",
	Long:  `Turn on power to the speakers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := setPower(true); err != nil {
			fmt.Printf("error setting power: %v\n", err)
			return
		}
		if err := showPower(quiet); err != nil {
			fmt.Printf("error showing power: %v\n", err)
		}
	},
}

var powerToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle power to speakers",
	Long:  `Toggle power to the speakers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := togglePower(); err != nil {
			fmt.Printf("error setting power: %v\n", err)
			return
		}
		if err := showPower(quiet); err != nil {
			fmt.Printf("error showing power: %v\n", err)
		}
	},
}

func setPower(on bool) error {
	hub := platin.NewHub(address)
	return hub.SetPower(on)
}

func togglePower() error {
	hub := platin.NewHub(address)
	return hub.TogglePower()
}
