package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	var repoURL string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "remove a Gloo feature repository",
		RunE: func(c *cobra.Command, args []string) error {
			return runDelete(repoURL)
		},
	}
	cmd.Flags().StringVarP(&repoURL, "repository", "r", "", "URL of the repository to remove")
	cmd.MarkFlagRequired("repository")
	return cmd
}

func runDelete(repoURL string) error {
	// remove features for the repo
	featureStore := &feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	if err := featureStore.RemoveForRepo(repoURL); err != nil {
		fmt.Printf("unable to remove features for repository %s: %q\n", repoURL, err)
		return nil
	}
	// remove the repo
	repoStore := &feature.FileRepoStore{Filename: feature.ReposFileName}
	if err := repoStore.Remove(repoURL); err != nil {
		fmt.Printf("Uable to remove repository %s: %q\n", repoURL, err)
	}

	return nil
}
