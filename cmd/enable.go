package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

// EnableCmd enables a given feature
func EnableCmd() *cobra.Command {
	var featureName string

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable a feature from feature list",
		RunE: func(c *cobra.Command, args []string) error {
			return runChangeStatus(featureName, true)
		},
	}
	cmd.Flags().StringVarP(&featureName, "name", "n", "", "name of feature to enable")
	cmd.MarkFlagRequired("name")
	return cmd
}

// DisableCmd disables a give feature
func DisableCmd() *cobra.Command {
	var featureName string

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable a feature from feature list",
		RunE: func(c *cobra.Command, args []string) error {
			return runChangeStatus(featureName, false)
		},
	}
	cmd.Flags().StringVarP(&featureName, "name", "n", "", "name of feature to disable")
	cmd.MarkFlagRequired("name")
	return cmd
}

func runChangeStatus(featureName string, status bool) error {
	store := &feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	existing, err := store.List()
	if err != nil {
		fmt.Printf("Unable to load feature list: %q\n", err)
		return nil
	}
	for i, f := range existing {
		if featureName == f.Name {
			existing[i].Enabled = status
			if err := store.Update(existing[i]); err != nil {
				fmt.Printf("Unable to update feature %s\n", featureName)
			}
			return nil
		}
	}
	fmt.Printf("Unable to find feature %s\n", featureName)
	return nil
}
