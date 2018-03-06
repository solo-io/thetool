package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

// AddCmd adds a Gloo feature repository at specific commit hash
// It downloads and parses features.json and adds the features
// listed in the file
func AddCmd() *cobra.Command {
	var repoURL string
	var commitHash string
	var verbose bool

	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"update"},
		Short:   "add or update a Gloo feature repository",
		Long:    "add or update a Gloo feature repository and all the features in the repository",
		Run: func(c *cobra.Command, args []string) {
			err := runAdd(verbose, repoURL, commitHash)
			if err != nil {
				fmt.Println("unable to add/update the repository", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&repoURL, "repository", "r", "", "repository URL")
	flags.StringVarP(&commitHash, "commit", "c", "", "commit hash")
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose logging")

	cmd.MarkFlagRequired("repository")
	cmd.MarkFlagRequired("commit")

	return cmd
}

func runAdd(verbose bool, repo, hash string) error {
	repoStore := &feature.FileRepoStore{Filename: feature.ReposFileName}
	if !downloader.SupportedURL(repo) {
		return fmt.Errorf("unsupported repository URL %s\nShould either end in '.git' or be HTTP/HTTPS", repo)
	}

	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		return errors.Wrapf(err, "unable to load configuration from %s", config.ConfigFile)
	}

	err = downloader.Download(repo, hash, conf.WorkDir, verbose)
	if err != nil {
		return errors.Wrapf(err, "unable to download repository %s", repo, err)
	}

	mf, err := feature.LoadManifest(filepath.Join(conf.WorkDir, downloader.RepoDir(repo)))
	if err != nil {
		return errors.Wrapf(err, "unable to load features manifest for repository %s", repo)
	}
	if len(mf) == 0 {
		return fmt.Errorf("not adding repository %s as it does not contain any Gloo features", repo)
	}

	features := feature.ToFeatures(repo, hash, mf)
	featureStore := &feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	err = featureStore.AddOrUpdateAll(features)
	if err != nil {
		return errors.Wrapf(err, "unable to add features found in repo %s", repo)
	}
	updated, err := repoStore.AddOrUpdate(feature.Repository{URL: repo, Commit: hash})
	if err != nil {
		return errors.Wrapf(err, "unable to save repo %s", repo)
	}
	if updated {
		fmt.Printf("Updated repository %s to commit hash %s\n", repo, hash)
	} else {
		fmt.Printf("Added repository %s with commit hash %s\n", repo, hash)
	}
	return nil
}
