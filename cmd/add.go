package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

// AddCmd adds the provided feature to feature list and enables it.
// Ading a feature checkes out the code from given repository for the
// given commit hash
func AddCmd() *cobra.Command {
	var featureName string
	var featureRepository string
	var featureHash string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add a feature",
		RunE: func(c *cobra.Command, args []string) error {
			return runAdd(featureName, featureRepository, featureHash, verbose)
		},
	}

	cmd.Flags().StringVarP(&featureName, "name", "n", "", "Feature name")
	cmd.Flags().StringVarP(&featureRepository, "repository", "r", "", "Repository URL")
	cmd.Flags().StringVarP(&featureHash, "commit", "c", "", "Commit hash")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("commit")

	return cmd
}

func runAdd(name, repo, hash string, verbose bool) error {
	f := feature.Feature{
		Name:       name,
		Version:    hash,
		Repository: repo,
		Enabled:    true,
	}
	existing, err := feature.LoadFromFile(dataFile)
	if err != nil {
		fmt.Printf("Unable to load feature list: %q\n", err)
		return nil
	}

	// check it isn't already existing feature
	for _, e := range existing {
		if e.Name == name {
			fmt.Printf("Feature %s already exists\n", name)
			return nil
		}
	}
	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s: %q\n", config.ConfigFile, err)
		return nil
	}
	// let's get the external dependency
	err = downloader.Download(f, conf.WorkDir, verbose)
	if err != nil {
		fmt.Printf("Unable to download repository %s: %q\n", f.Repository, err)
		return nil
	}

	err = feature.SaveToFile(append(existing, f), dataFile)
	if err != nil {
		fmt.Printf("Unable to save feature %s: %q\n", f.Name, err)
	}
	return nil
}
