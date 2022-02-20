package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/soellman/platin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(volumeCmd)
	volumeCmd.AddCommand(volumeSetCmd)
}

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Show and set the volume",
	Long:  `Show current volume and set the volume level. Valid volumes are 0-100`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showVolume(quiet); err != nil {
			fmt.Printf("error showing volume: %v\n", err)
		}
	},
}

var volumeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the volume",
	Long:  `Set the volume level. Valid volumes are 0-100`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please specify volume level (0-100)")
		} else if len(args) > 1 {
			return errors.New("provide only a single volume argument")
		}
		vol, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("volume must be an integer")
		}
		if vol < 0 || vol > 100 {
			return errors.New("volume must be from 0-100")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		vol, _ := strconv.Atoi(args[0])
		if err := setVolume(vol); err != nil {
			fmt.Printf("error setting volume: %v\n", err)
			return
		}
		if err := showVolume(quiet); err != nil {
			fmt.Printf("error showing volume: %v\n", err)
		}
	},
}

func setVolume(volume int) error {
	hub := platin.NewHub(address)
	return hub.SetVolume(volume)
}
