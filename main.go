package main

import (
	"context"
	"time"

	checkpoint "github.com/solo-io/go-checkpoint"
	"github.com/solo-io/thetool/cmd"
	"github.com/solo-io/thetool/cmd/addon"
	"github.com/spf13/cobra"
)

var Version = "DEV"

func main() {
	start := time.Now()
	defer telemetry(start)
	rootCmd := &cobra.Command{
		Use:     "thetool",
		Short:   "Build Tool",
		Long:    "Build the Universe and gloo things together",
		Version: Version,
	}

	rootCmd.AddCommand(cmd.InitCmd())
	rootCmd.AddCommand(cmd.ConfigureCmd())
	rootCmd.AddCommand(cmd.ListReposCmd())
	rootCmd.AddCommand(cmd.AddCmd())
	rootCmd.AddCommand(cmd.DeleteCmd())
	rootCmd.AddCommand(cmd.EnableCmd())
	rootCmd.AddCommand(cmd.DisableCmd())
	rootCmd.AddCommand(cmd.ListFeaturesCmd())
	rootCmd.AddCommand(cmd.BuildCmd())
	rootCmd.AddCommand(cmd.CleanCmd())
	rootCmd.AddCommand(cmd.DeployCmd())
	rootCmd.AddCommand(addon.AddonCmd())

	rootCmd.Execute()
}

func telemetry(t time.Time) {
	ctx := context.Background()
	report := &checkpoint.ReportParams{
		Product:       "thetool",
		Version:       Version,
		StartTime:     t,
		EndTime:       time.Now(),
		SignatureFile: ".sig",
	}
	checkpoint.Report(ctx, report)
}
