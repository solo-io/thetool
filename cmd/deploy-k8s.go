package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/spf13/cobra"
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
			runDeployK8S(verbose, dryRun, dockerUser)
		},
	}
	return cmd
}

func runDeployK8S(verbose, dryRun bool, dockerUser string) error {
	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s: %q\n", config.ConfigFile, err)
		return nil
	}
	if dockerUser == "" {
		dockerUser = conf.DockerUser
	}
	if dockerUser == "" {
		return fmt.Errorf("need Docker user for referencing Docker images")
	}
	enabled, err := loadEnabledFeatures()
	if err != nil {
		fmt.Println("Unable to get enabled features")
		return nil
	}
	fmt.Printf("Building with %d features\n", len(enabled))
	featuresHash := featuresHash(enabled)
	err = generateHelmValues(false, featuresHash, dockerUser)
	if err != nil {
		fmt.Println("Unable to generate Helm chart values")
	}

	// install Glue using Helm
	return nil
}

func generateHelmValues(verbose bool, featureHash, user string) error {
	fmt.Println("Generating Helm Chart values...")
	filename := "glue-chart.yaml"
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()
	err = helmValuesTemplate.Execute(f, map[string]string{
		"EnvoyImage": user + "/envoy",
		"EnvoyTag":   featureHash,
		"GlueImage":  user + "/glue",
		"GlueTag":    featureHash,
	})
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}
