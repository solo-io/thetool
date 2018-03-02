package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func DeployLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "deploy the universe locally",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("Deploying locally has not been implemented yet.")
		},
	}
	return cmd
}

func runDeployLocal() {
	/// storage option is local files
	// run gloo locally?
	// run envoy in docker
}
