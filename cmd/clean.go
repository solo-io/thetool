package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "clean up files and directory",
		Run: func(c *cobra.Command, args []string) {
			runClean()
		},
	}
	return cmd
}

func runClean() {
	toDelete := []string{
		"BUILD", "WORKSPACE", "Dockerfile.envoy", "envoy",
		"build-envoy.sh", "bazel-bin", "bazel-genfiles", "bazel-out",
		"bazel-source", "bazel-testlogs", "gloo-chart.yaml",
		"build-gloo.sh", "gloo", "cache",
	}

	for _, f := range toDelete {
		if err := os.RemoveAll(f); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Unable to delete %v: %q\n", f, err)
			}
		}
	}
}
