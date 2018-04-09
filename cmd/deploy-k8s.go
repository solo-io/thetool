package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/solo-io/thetool/cmd/addon"
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
	installFile   = "install.yaml"
)

var (
	releaseRegex = regexp.MustCompile("\\s+release: RELEASE-NAME")
)

type k8sDeployOptions struct {
	resume            bool
	generateInstall   bool
	installPrometheus bool
	namespace         string
	releaseName       string
}

func DeployK8SCmd() *cobra.Command {
	options := k8sDeployOptions{}
	cmd := &cobra.Command{
		Use:   "k8s",
		Short: "deploy the universe in Kubernetes",
		Long: `
You can use dry-run to just generate gloo-chart.yaml file with parameters for Helm.
After it you can edit the file to make any changes and continue with --resume flag.

Use generate-install to generate a single install.yaml file that can be used with
kubectl`,
		Run: func(c *cobra.Command, args []string) {
			f := c.InheritedFlags()
			verbose, _ := f.GetBool("verbose")
			dryRun, _ := f.GetBool("dry-run")
			dockerUser, _ := f.GetString("docker-user")
			imageTag, _ := f.GetString("image-tag")
			options.installPrometheus = addon.InstallPrometheus()
			if err := runDeployK8S(verbose, dryRun, dockerUser, imageTag, options); err != nil {
				fmt.Printf("Unable to deploy Gloo: %q\n", err)
			}
		},
	}
	cmd.Flags().BoolVarP(&options.resume, "resume", "r", false, "resume deployment with existing "+glooChartYaml)
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "gloo-system", "namespace to deploy gloo and its components")
	cmd.Flags().StringVar(&options.releaseName, "release-name", "", "release name for Helm")
	cmd.Flags().BoolVarP(&options.generateInstall, "generate-install", "g", false, "generate install.yaml")
	return cmd
}

func runDeployK8S(verbose, dryRun bool, dockerUser, imageTag string, options k8sDeployOptions) error {
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

	if options.generateInstall && options.resume {
		return fmt.Errorf("can't use generate-install and resume at the same time")
	}

	if !glooDownloaded(conf.GlooRepo) {
		fmt.Printf("Downloading Gloo from %s\n", conf.GlooRepo)
		if err := downloader.Download(conf.GlooRepo, conf.GlooHash, config.WorkDir, verbose); err != nil {
			return errors.Wrap(err, "unable to download Gloo")
		}
	}
	helmDir := filepath.Join(config.WorkDir, downloader.RepoDir(conf.GlooRepo), "install", "helm")
	templateFile := filepath.Join(helmDir, "bootstrap.yaml")
	if !options.resume {
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

		if err := saveBootstrap(templateFile, options.namespace); err != nil {
			return errors.Wrap(err, "unable to generate boostrap YAML")
		}
	}

	glooChartDir := filepath.Join(helmDir, "gloo")
	if options.generateInstall {
		// use helm to generate final yaml
		helmArgs := []string{"template", glooChartDir, "-f", glooChartYaml}
		if options.namespace != "" {
			helmArgs = append(helmArgs, "--namespace", options.namespace)
		}
		b, err := util.RunCmdCaptureOut(true, false, "helm", helmArgs...)
		if err != nil {
			return errors.Wrap(err, "unable to run Helm template")

		}
		content := b.Bytes()
		if options.releaseName == "" {
			content = bytes.Replace(content, []byte("RELEASE-NAME-"), []byte(""), -1)
			content = releaseRegex.ReplaceAll(content, []byte{})
		}

		buffer := &bytes.Buffer{}
		if err := generateBootstrap(buffer, templateFile, options.namespace); err != nil {
			return errors.Wrap(err, "unable to generate bootstrap for install")
		}
		buffer.WriteByte('\n')
		buffer.Write(content)
		if err := ioutil.WriteFile(installFile, buffer.Bytes(), 0644); err != nil {
			return errors.Wrapf(err, "unable to save %s", installFile)
		}
	} else {
		// do we need to install prometheus?
		if options.installPrometheus {
			err = util.RunCmd(verbose, dryRun, "helm",
				"repo", "add", "coreos", "https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/",
			)
			if err != nil {
				return errors.Wrap(err, "unable to add repository for prometheus operator")
			}
			err = util.RunCmd(verbose, dryRun, "helm",
				"install", "coreos/prometheus-operator", "--name", "prometheus-operator",
				"--namespace", options.namespace)
			if err != nil {
				return errors.Wrap(err, "unable to install prometheus operator")
			}
		}
		// bootstrap using kubectl
		bootstrapArgs := []string{"apply", "-f", bootstrapFile}
		err = util.RunCmd(verbose, dryRun, "kubectl", bootstrapArgs...)
		if err != nil {
			return errors.Wrap(err, "unable to run bootstrap")
		}

		// install Gloo using Helm
		helmArgs := []string{"install", glooChartDir, "-f", glooChartYaml}

		if options.namespace != "" {
			helmArgs = append(helmArgs, "--namespace", options.namespace)
		}
		return util.RunCmd(verbose, dryRun, "helm", helmArgs...)
	}

	return nil
}

func glooDownloaded(glooRepo string) bool {
	glooDir := filepath.Join(config.WorkDir, downloader.RepoDir(glooRepo))
	_, err := os.Stat(glooDir)
	return err == nil
}

func generateHelmValues(verbose bool, imageTag, user string) error {
	fmt.Println("Generating Helm Chart values...")
	filename := glooChartYaml
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()

	addons, err := addon.List()
	if err != nil {
		return errors.Wrap(err, "unable to load addons")
	}
	err = helmValuesTemplate.Execute(f, map[string]interface{}{
		"EnvoyImage": user + "/envoy",
		"EnvoyTag":   imageTag,
		"GlooImage":  user + "/gloo",
		"GlooTag":    imageTag,
		"DockerUser": user,
		"Addons":     addons,
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

func generateBootstrap(w io.Writer, templateFile, namespace string) error {
	if namespace == "" {
		namespace = "gloo-system"
	}
	t := template.Must(template.ParseFiles(templateFile))
	return t.Execute(w, map[string]string{
		"Namespace": namespace,
	})
}

func saveBootstrap(templateFile, namespace string) error {
	f, err := os.OpenFile(bootstrapFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return generateBootstrap(f, templateFile, namespace)
}
