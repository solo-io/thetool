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
		Run: func(c *cobra.Command, args []string) {
			runList()
		},
	}
	return cmd
}

func runList() {
	features, err := feature.LoadFromFile(dataFile)
	if err != nil {
		fmt.Printf("Unable to load feature list: %q\n", err)
		return
	}
	if len(features) == 0 {
		fmt.Println("No features added yet!")
	}
	for _, f := range features {
		fmt.Println("Name:       ", f.Name)
		fmt.Println("Repository: ", f.Repository)
		fmt.Println("Commit:     ", f.Version)
		fmt.Println("Enabled:    ", f.Enabled)
		fmt.Println("")
	}
}
