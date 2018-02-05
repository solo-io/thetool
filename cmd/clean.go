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
		"build.sh", "bazel-bin", "bazel-genfiles", "bazel-out",
		"bazel-source", "bazel-testlogs", "glue-chart.yaml",
		"Dockerfile.glue", "build-glue.sh", "glue",
	}

	for _, f := range toDelete {
		if err := os.Remove(f); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Unable to delete %v: %q\n", f, err)
			}
		}
	}
}
