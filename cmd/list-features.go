package cmd

import (
	"fmt"
	"os"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func ListFeaturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all registered features",
		Run: func(c *cobra.Command, args []string) {
			runListFeatures()
		},
	}
	return cmd
}

func runListFeatures() {
	store := &feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	features, err := store.List()
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Please add feature repository before listing features")
		} else {
			fmt.Printf("Unable to load feature list: %q\n", err)
		}
		return
	}
	if len(features) == 0 {
		fmt.Println("No repositories with features added yet!")
	}
	for _, f := range features {
		fmt.Println("Repository:      ", f.Repository)
		fmt.Println("Name:            ", f.Name)
		fmt.Println("Gloo Directory:  ", f.GlooDir)
		fmt.Println("Envoy Directory: ", f.EnvoyDir)
		fmt.Println("Enabled:         ", f.Enabled)
		if len(f.Tags) != 0 {
			fmt.Println("Tags:            ", f.Tags)
		}
		fmt.Println("")
	}
}
