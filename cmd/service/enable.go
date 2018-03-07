package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func enableCmd() *cobra.Command {
	var serviceName string
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable a supporting service",
		Run: func(c *cobra.Command, args []string) {
			runChangeStatus(serviceName, true)
		},
	}
	cmd.Flags().StringVarP(&serviceName, "name", "n", "", "name of service to enable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func disableCmd() *cobra.Command {
	var serviceName string
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable a supporting service",
		Run: func(c *cobra.Command, args []string) {
			runChangeStatus(serviceName, false)
		},
	}
	cmd.Flags().StringVarP(&serviceName, "name", "n", "", "name of service to disable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runChangeStatus(serviceName string, status bool) {
	services, err := load(serviceFilename)
	if err != nil {
		fmt.Println("Unable to load list of services", err)
		return
	}
	for _, s := range services {
		if s.Name == serviceName {
			s.Enable = status
		}
	}

	if err := save(serviceFilename, services); err != nil {
		fmt.Println("Unable to update list of services", err)
	}
}
