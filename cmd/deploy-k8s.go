package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/util"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/spf13/cobra"
)

const (
	glooChartYaml = "gloo-chart.yaml"
)

func DeployK8SCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "deploy the universe in Kubernetes",
		Run: func(c *cobra.Command, args []string) {
			f := c.InheritedFlags()
			verbose, _ := f.GetBool("verbose")
			dryRun, _ := f.GetBool("dry-run")
			dockerUser, _ := f.GetString("docker-user")
			imageTag, _ := f.GetString("image-tag")
			if err := runDeployK8S(verbose, dryRun, dockerUser, imageTag); err != nil {
				fmt.Printf("Unable to deploy Gloo: %q\n", err)
			}
		},
	}
	return cmd
}

func runDeployK8S(verbose, dryRun bool, dockerUser, imageTag string) error {
	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		return errors.Wrapf(err, "unable to load configuration from %s", config.ConfigFile)
	}
	if dockerUser == "" {
		dockerUser = conf.DockerUser
	}
	if dockerUser == "" {
		return fmt.Errorf("need Docker user for referencing Docker images")
	}
	enabled, err := loadEnabledFeatures()
	if err != nil {
		return errors.Wrap(err, "unable to get enabled features")
	}
	fmt.Printf("Building with %d features\n", len(enabled))
	if imageTag == "" {
		imageTag = featuresHash(enabled)
	}
	if err := generateHelmValues(false, imageTag, dockerUser); err != nil {
		return errors.Wrap(err, "unable to generate Helm chart values")
	}

	if !dryRun {
		fmt.Printf("Downloading Gloo chart from %s", conf.GlooChartRepo)
		if err := downloader.Download(conf.GlooChartRepo, conf.GlooChartHash, conf.WorkDir, verbose); err != nil {
			return errors.Wrap(err, "unable to download Gloo chart")
		}
	}
	// install Gloo using Helm
	helmArgs := []string{"install", filepath.Join(conf.WorkDir, downloader.RepoDir(conf.GlooChartRepo)),
		"-f", glooChartYaml}
	return util.RunCmd(verbose, dryRun, "helm", helmArgs...)
}

func generateHelmValues(verbose bool, featureHash, user string) error {
	fmt.Println("Generating Helm Chart values...")
	filename := glooChartYaml
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()
	err = helmValuesTemplate.Execute(f, map[string]string{
		"EnvoyImage": user + "/envoy",
		"EnvoyTag":   featureHash,
		"GlooImage":  user + "/gloo",
		"GlooTag":    featureHash,
	})
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}
