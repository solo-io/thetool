package main

import (
	"github.com/solo-io/thetool/cmd"
	"github.com/spf13/cobra"
)

//go:generate bash -c "sed \"s/VERSION_TEMPLATE/`git describe --tags`/\" version.go.template > cmd/version_xx_autogen.go"
func main() {
	rootCmd := &cobra.Command{
		Use:   "thetool",
		Short: "Build Tool",
		Long:  "Build the Universe and gloo things together",
	}

	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ListReposCmd())
	rootCmd.AddCommand(cmd.AddCmd())
	rootCmd.AddCommand(cmd.DeleteCmd())
	rootCmd.AddCommand(cmd.EnableCmd())
	rootCmd.AddCommand(cmd.DisableCmd())
	rootCmd.AddCommand(cmd.ListFeaturesCmd())
	rootCmd.AddCommand(cmd.BuildCmd())
	rootCmd.AddCommand(cmd.CleanCmd())
	rootCmd.AddCommand(cmd.DeployCmd())
	rootCmd.AddCommand(cmd.VersionCmd())

	rootCmd.Execute()
}
