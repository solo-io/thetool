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
			fmt.Println("not implemetned deploying locally")
		},
	}
	return cmd
}

func runDeployLocal() {
	/// storage option is local files?
	// run glue locally
	// run envoy in docker
}
