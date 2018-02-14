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
	cmd.Flags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runDelete(featureToRemove string) error {
	existing, err := feature.LoadFromFile(dataFile)
	if err != nil {
		fmt.Printf("Unable to load existing features: %q\n", err)
		return nil
	}
	var updated []feature.Feature
	for _, f := range existing {
		if featureToRemove != f.Name {
			updated = append(updated, f)
		}
	}

	if len(updated) == len(existing) {
		fmt.Printf("Didn't find feature %s\n", featureToRemove)
		return nil
	}

	err = feature.SaveToFile(updated, dataFile)
	if err != nil {
		fmt.Printf("Unable to update feature list: %q\n", err)
	}
	return nil
}
