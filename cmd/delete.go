package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

var (
	featureToRemove string
)

func DeleteCmd() *cobra.Command {
	var featureName string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "remove a feature from feature list",
		RunE: func(c *cobra.Command, args []string) error {
			return runDelete(featureName)
		},
	}
	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
	return cmd
}

func runDelete(featureToRemove string) error {
	if featureToRemove == "" {
		return fmt.Errorf("name of the feature to remove can't be empty")
	}
	existing, err := feature.LoadFromFile("features.json")
	if err != nil {
		return err
	}
	var updated []feature.Feature
	for _, f := range existing {
		if featureToRemove != f.Name {
			updated = append(updated, f)
		}
	}

	if len(updated) == len(existing) {
		return fmt.Errorf("didn't find feature %s", featureToRemove)
	}

	return feature.SaveToFile(updated, "features.json")
}
