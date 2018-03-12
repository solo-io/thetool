package addon

import (
	"fmt"

	"github.com/spf13/cobra"
)

func markInstallCmd() *cobra.Command {
	var addonName string
	cmd := &cobra.Command{
		Use:   "install",
		Short: "mark an addon for install",
		Run: func(c *cobra.Command, args []string) {
			runMarkInstall(addonName, false)
		},
	}
	cmd.Flags().StringVarP(&addonName, "name", "n", "", "name of an addon to mark as install")
	cmd.MarkFlagRequired("name")
	return cmd
}

func markConfigOnlyCmd() *cobra.Command {
	var addonName string
	cmd := &cobra.Command{
		Use:   "config-only",
		Short: "mark an addon as config-only",
		Run: func(c *cobra.Command, args []string) {
			runMarkInstall(addonName, true)
		},
	}
	cmd.Flags().StringVarP(&addonName, "name", "n", "", "name of an addon to mark as install only")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runMarkInstall(addonName string, configOnly bool) {
	addons, err := load(addonFilename)
	if err != nil {
		fmt.Println("Unable to load list of addons", err)
		return
	}
	for _, a := range addons {
		if a.Name == addonName {
			a.ConfigOnly = &configOnly
		}
	}

	if err := save(addonFilename, addons); err != nil {
		fmt.Println("Unable to update list of addons", err)
	}
}
