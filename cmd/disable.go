package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func DisableCmd() *cobra.Command {
	var featureName string

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable a feature from feature list",
		RunE: func(c *cobra.Command, args []string) error {
			return runDisable(featureName)
		},
	}
	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
	return cmd
}

func runDisable(featureToDisable string) error {
	if featureToDisable == "" {
		return fmt.Errorf("name of the feature to disable can't be empty")
	}
	existing, err := feature.LoadFromFile("features.json")
	if err != nil {
		return err
	}
	found := false
	for i, f := range existing {
		if featureToDisable == f.Name {
			existing[i].Enabled = false
			found = true
		}
	}
	if !found {
		return fmt.Errorf("unable to find feature %s", featureToDisable)
	}
	return feature.SaveToFile(existing, "features.json")
}
