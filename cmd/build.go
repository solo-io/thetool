package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/solo-io/thetool/pkg/component"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var jobs int
	config := component.BuilderConfig{}
	cmd := &cobra.Command{
		Use:   "build [target to build]",
		Short: "build the universe",
		Long: `
Build gloo and its components.
Supported components are:
  envoy, gloo, function-discovery, ingress, k8s-discovery or all`,
		ValidArgs: component.Components(),
		Args:      cobra.OnlyValidArgs,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("please specify a build target")
			}
			target := strings.ToLower(args[0])
			return runBuild(jobs, config, target)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&config.Verbose, "verbose", "v", false, "show verbose build log")
	flags.BoolVarP(&config.DryRun, "dry-run", "d", false, "dry run; only generate build file")
	flags.BoolVar(&config.UseCache, "cache", true, "use cache for builds")
	flags.BoolVarP(&config.PublishImage, "publish", "p", true, "publish Docker images to registry")
	flags.StringVarP(&config.ImageTag, "image-tag", "t", "", "tag for Docker images; uses auto-generated hash if empty")
	flags.StringVarP(&config.DockerUser, "docker-user", "u", "", "Docker user for publishing images")
	flags.StringVar(&config.SSHKeyFile, "ssh-key", "", "file containg SSH key for git to use with private repositories")
	flags.IntVarP(&jobs, "jobs", "j", 1, "number of jobs to run simultaneously")
	return cmd
}

func runBuild(jobs int, buildConfig component.BuilderConfig, target string) error {
	var err error
	buildConfig.Config, err = config.Load(config.ConfigFile)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s: %q\n", config.ConfigFile, err)
		return nil
	}
	if buildConfig.DockerUser == "" {
		buildConfig.DockerUser = buildConfig.Config.DockerUser
	}
	if buildConfig.DockerUser == "" && buildConfig.PublishImage {
		return fmt.Errorf("need Docker user ID to publish images")
	}
	buildConfig.Enabled, err = loadEnabledFeatures()
	if err != nil {
		return err
	}
	if buildConfig.PublishImage {
		fmt.Printf("Building and publishing with %d features\n", len(buildConfig.Enabled))
	} else {
		fmt.Printf("Building with %d features\n", len(buildConfig.Enabled))
	}
	if buildConfig.ImageTag == "" {
		buildConfig.ImageTag = featuresHash(buildConfig.Enabled)
	}

	jobCh := make(chan func(), 10)
	if jobs < 1 {
		jobs = 1
	}
	for w := 0; w < jobs; w++ {
		go worker(jobCh)
	}

	var wg sync.WaitGroup
	for i := range component.Builders {
		b := component.Builders[i]
		if target == component.All || target == b.Name {
			wg.Add(1)
			jobCh <- func() {
				defer wg.Done()
				b.Builder(buildConfig)
			}
		}
	}

	close(jobCh)
	wg.Wait()
	return nil
}

func worker(jobs <-chan func()) {
	for j := range jobs {
		j()
	}
}
