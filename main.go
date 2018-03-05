package main

import (
	"github.com/solo-io/thetool/cmd"
	"github.com/spf13/cobra"
)

var Version = "DEV"

func main() {
	rootCmd := &cobra.Command{
		Use:     "thetool",
		Short:   "Build Tool",
		Long:    "Build the Universe and gloo things together",
		Version: Version,
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

	rootCmd.Execute()
}
