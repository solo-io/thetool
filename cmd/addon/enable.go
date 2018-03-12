package addon

import (
	"fmt"

	"github.com/spf13/cobra"
)

func enableCmd() *cobra.Command {
	var addonName string
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable an addon",
		Run: func(c *cobra.Command, args []string) {
			runChangeStatus(addonName, true)
		},
	}
	cmd.Flags().StringVarP(&addonName, "name", "n", "", "name of addon to enable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func disableCmd() *cobra.Command {
	var addonName string
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable an addon",
		Run: func(c *cobra.Command, args []string) {
			runChangeStatus(addonName, false)
		},
	}
	cmd.Flags().StringVarP(&addonName, "name", "n", "", "name of addon to disable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runChangeStatus(addonName string, status bool) {
	addons, err := load(addonFilename)
	if err != nil {
		fmt.Println("Unable to load list of addons", err)
		return
	}
	for _, a := range addons {
		if a.Name == addonName {
			a.Enable = status
		}
	}

	if err := save(addonFilename, addons); err != nil {
		fmt.Println("Unable to update list of addons", err)
	}
}
