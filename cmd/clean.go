package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "clean up files and directory",
		RunE: func(c *cobra.Command, args []string) error {
			return runClean()
		},
	}
	return cmd
}

func runClean() error {
	toDelete := []string{
		"BUILD", "WORKSPACE", "Dockerfile.envoy", "envoy",
		"build.sh", "bazel-bin", "bazel-genfiles", "bazel-out",
		"bazel-source", "bazel-testlogs", "glue-chart.yaml",
	}

	for _, f := range toDelete {
		if err := os.Remove(f); err != nil {
			return errors.Wrap(err, "unable to delete "+f)
		}
	}
	return nil
}
