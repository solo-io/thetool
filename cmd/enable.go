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
	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
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
	cmd.PersistentFlags().StringVarP(&featureName, "name", "n", "", "name of feature to remove")
	return cmd
}

func runChangeStatus(featureName string, status bool) error {
	if featureName == "" {
		return fmt.Errorf("name of the feature can't be empty")
	}
	existing, err := feature.LoadFromFile(dataFile)
	if err != nil {
		fmt.Printf("Unable to load feature list: %q\n", err)
		return nil
	}
	found := false
	for i, f := range existing {
		if featureName == f.Name {
			existing[i].Enabled = status
			found = true
		}
	}
	if !found {
		fmt.Printf("Unable to find feature %s\n", featureName)
		return nil
	}
	err = feature.SaveToFile(existing, dataFile)
	if err != nil {
		fmt.Printf("Unable to save feature list: %q\n", err)
	}
	return nil
}
