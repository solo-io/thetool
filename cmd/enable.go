package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func EnableCmd() *cobra.Command {
	var featureName string

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable a feature from feature list",
		RunE: func(c *cobra.Command, args []string) error {
			return runEnable(featureName)
		},
	}
	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
	return cmd
}

func runEnable(featureToEnable string) error {
	if featureToEnable == "" {
		return fmt.Errorf("name of the feature to enable can't be empty")
	}
	existing, err := feature.LoadFromFile("features.json")
	if err != nil {
		return err
	}
	found := false
	for i, f := range existing {
		if featureToEnable == f.Name {
			existing[i].Enabled = true
			found = true
		}
	}
	if !found {
		return fmt.Errorf("unable to find feature %s", featureToEnable)
	}
	return feature.SaveToFile(existing, "features.json")
}
