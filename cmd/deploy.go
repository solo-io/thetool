package cmd

import "github.com/spf13/cobra"

func DeployCmd() *cobra.Command {
	var verbose bool
	var dryRun bool
	var dockerUser string
	var imageTag string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the universe",
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show verbose build log")
	cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run; only generate build file")
	cmd.PersistentFlags().StringVarP(&dockerUser, "docker-user", "u", "", "Docker user for publishing images")
	cmd.PersistentFlags().StringVarP(&imageTag, "image-tag", "t", "", "tag for Docker images; uses auto-generated hash if empty")

	cmd.AddCommand(DeployLocalCmd())
	cmd.AddCommand(DeployK8SCmd())
	cmd.AddCommand(DeployK8SOutCmd())
	return cmd
}
