package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all registered features",
		RunE: func(c *cobra.Command, args []string) error {
			return runList()
		},
	}
	return cmd
}

func runList() error {
	features, err := feature.LoadFromFile("features.json")
	if err != nil {
		return err
	}

	for _, f := range features {
		fmt.Println("Name:       ", f.Name)
		fmt.Println("Repository: ", f.Repository)
		fmt.Println("Commit:     ", f.Version)
		fmt.Println("Enabled:    ", f.Enabled)
		fmt.Println("")
	}
	return nil
}
