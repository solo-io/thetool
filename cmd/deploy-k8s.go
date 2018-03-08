package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/solo-io/thetool/cmd/service"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/spf13/cobra"
)

const (
	glooChartYaml = "gloo-chart.yaml"
	bootstrapFile = "gloo-bootstrap.yaml"
)

func DeployK8SCmd() *cobra.Command {
	var resume bool
	var namespace string
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "deploy the universe in Kubernetes",
		Long: `
You can use dry-run to just generate gloo-chart.yaml file with parameters for Helm.
After it you can edit the file to make any changes and continue with --resume flag.`,
		Run: func(c *cobra.Command, args []string) {
			f := c.InheritedFlags()
			verbose, _ := f.GetBool("verbose")
			dryRun, _ := f.GetBool("dry-run")
			dockerUser, _ := f.GetString("docker-user")
			imageTag, _ := f.GetString("image-tag")
			if err := runDeployK8S(verbose, dryRun, dockerUser, imageTag, namespace, resume); err != nil {
				fmt.Printf("Unable to deploy Gloo: %q\n", err)
			}
		},
	}
	cmd.Flags().BoolVarP(&resume, "resume", "r", false, "resume deployment with existing "+glooChartYaml)
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "gloo-system", "namespace to deploy gloo and its components")
	return cmd
}

func runDeployK8S(verbose, dryRun bool, dockerUser, imageTag, namespace string, resume bool) error {
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

	if !dryRun {
		fmt.Printf("Downloading Gloo chart from %s\n", conf.GlooChartRepo)
		if err := downloader.Download(conf.GlooChartRepo, conf.GlooChartHash, conf.WorkDir, verbose); err != nil {
			return errors.Wrap(err, "unable to download Gloo chart")
		}
	}

	if !resume {
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

		templateFile := filepath.Join(conf.WorkDir, downloader.RepoDir(conf.GlooChartRepo), "helm", "bootstrap.yaml")
		if !dryRun { // we can only generate the bootstrap if Gloo charts are downloaded
			if err := generateBootstrap(templateFile, namespace); err != nil {
				return errors.Wrap(err, "unable to generate pre Helm YAML")
			}
		}
	}

	bootstrapArgs := []string{"apply", "-f", bootstrapFile}
	err = util.RunCmd(verbose, dryRun, "kubectl", bootstrapArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to run bootstrap")
	}

	// install Gloo using Helm
	helmArgs := []string{"install", filepath.Join(
		conf.WorkDir, downloader.RepoDir(conf.GlooChartRepo), "helm", "gloo"),
		"-f", glooChartYaml}

	if namespace != "" {
		helmArgs = append(helmArgs, "--namespace", namespace)
	}
	return util.RunCmd(verbose, dryRun, "helm", helmArgs...)
}

func generateHelmValues(verbose bool, imageTag, user string) error {
	fmt.Println("Generating Helm Chart values...")
	filename := glooChartYaml
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()

	services, err := service.List()
	if err != nil {
		return errors.Wrap(err, "unable to load supporting services")
	}
	err = helmValuesTemplate.Execute(f, map[string]interface{}{
		"EnvoyImage": user + "/envoy",
		"EnvoyTag":   imageTag,
		"GlooImage":  user + "/gloo",
		"GlooTag":    imageTag,
		"DockerUser": user,
		"Services":   services,
	})
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}

func loadFeatures() ([]feature.Feature, error) {
	store := &feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	return store.List()
}

func generateBootstrap(templateFile, namespace string) error {
	f, err := os.OpenFile(bootstrapFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if namespace == "" {
		namespace = "gloo-system"
	}
	t := template.Must(template.ParseFiles(templateFile))
	t.Execute(f, map[string]string{
		"Namespace": namespace,
	})
	return nil
}
