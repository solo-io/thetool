package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func markInstallCmd() *cobra.Command {
	var serviceName string
	cmd := &cobra.Command{
		Use:   "install",
		Short: "mark a service for install",
		Run: func(c *cobra.Command, args []string) {
			runMarkInstall(serviceName, true)
		},
	}
	cmd.Flags().StringVarP(&serviceName, "name", "n", "", "name of service to disable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func markUninstallCmd() *cobra.Command {
	var serviceName string
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "mark a service for uninstall",
		Run: func(c *cobra.Command, args []string) {
			runMarkInstall(serviceName, false)
		},
	}
	cmd.Flags().StringVarP(&serviceName, "name", "n", "", "name of service to disable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runMarkInstall(serviceName string, install bool) {
	services, err := load(serviceFilename)
	if err != nil {
		fmt.Println("Unable to load list of services", err)
		return
	}
	for _, s := range services {
		if s.Name == serviceName {
			s.Install = install
		}
	}

	if err := save(serviceFilename, services); err != nil {
		fmt.Println("Unable to update list of services", err)
	}
}
