package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/envoy"
	"github.com/solo-io/thetool/pkg/gloo"
	"github.com/spf13/cobra"
)

type component int

const (
	componentAll component = iota
	componentEnvoy
	componentGloo
)

// BuildConfig stores the configuration used for building and
// publishing components of Gloo
type BuildConfig struct {
	UseCache      bool
	PublishImages bool
	DockerUser    string
	ImageTag      string
	SSHKeyFile    string
}

func BuildCmd() *cobra.Command {
	var verbose bool
	var dryRun bool
	var jobs int
	config := BuildConfig{}
	cmd := &cobra.Command{
		Use:       "build [target to build]",
		Short:     "build the universe",
		Long:      "build gloo, envoy or all",
		ValidArgs: []string{"envoy", "gloo", "all"},
		Args:      cobra.OnlyValidArgs,
		RunE: func(c *cobra.Command, args []string) error {
			target := componentAll
			if len(args) != 1 {
				return fmt.Errorf("please specify a build target")
			}
			switch strings.ToLower(args[0]) {
			case "envoy":
				target = componentEnvoy
			case "gloo":
				target = componentGloo
			default:
				target = componentAll
			}
			return runBuild(verbose, dryRun, jobs, config, target)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "show verbose build log")
	flags.BoolVarP(&dryRun, "dry-run", "d", false, "dry run; only generate build file")
	flags.BoolVar(&config.UseCache, "cache", true, "use cache for builds")
	flags.BoolVarP(&config.PublishImages, "publish", "p", true, "publish Docker images to registry")
	flags.StringVarP(&config.ImageTag, "image-tag", "t", "", "tag for Docker images; uses auto-generated hash if empty")
	flags.StringVarP(&config.DockerUser, "docker-user", "u", "", "Docker user for publishing images")
	flags.StringVar(&config.SSHKeyFile, "ssh-key", "", "file containg SSH key for git to use with private repositories")
	flags.IntVarP(&jobs, "jobs", "j", 1, "number of jobs to run simultaneously")
	return cmd
}

func runBuild(verbose, dryRun bool, jobs int, buildConfig BuildConfig, target component) error {
	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s: %q\n", config.ConfigFile, err)
		return nil
	}
	if buildConfig.DockerUser == "" {
		buildConfig.DockerUser = conf.DockerUser
	}
	if buildConfig.DockerUser == "" && buildConfig.PublishImages {
		return fmt.Errorf("need Docker user ID to publish images")
	}
	enabled, err := loadEnabledFeatures()
	if err != nil {
		return err
	}
	if buildConfig.PublishImages {
		fmt.Printf("Building and publishing with %d features\n", len(enabled))
	} else {
		fmt.Printf("Building with %d features\n", len(enabled))
	}
	if buildConfig.ImageTag == "" {
		buildConfig.ImageTag = featuresHash(enabled)
	}

	jobCh := make(chan func(), 10)
	if jobs < 1 {
		jobs = 1
	}
	for w := 0; w < jobs; w++ {
		go worker(jobCh)
	}

	var wg sync.WaitGroup
	if target == componentAll || target == componentEnvoy {
		wg.Add(1)
		jobCh <- func() {
			defer wg.Done()
			if err := envoy.Build(enabled, verbose, dryRun,
				buildConfig.UseCache, buildConfig.SSHKeyFile,
				conf.EnvoyHash, conf.WorkDir, conf.EnvoyBuilderHash); err != nil {
				fmt.Println(err)
				return
			}
			if err := envoy.Publish(verbose, dryRun,
				buildConfig.PublishImages, buildConfig.ImageTag, buildConfig.DockerUser); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	if target == componentAll || target == componentGloo {
		wg.Add(1)
		jobCh <- func() {
			defer wg.Done()
			if err := gloo.Build(enabled, verbose, dryRun,
				buildConfig.UseCache, buildConfig.SSHKeyFile,
				conf.GlooRepo, conf.GlooHash, conf.WorkDir); err != nil {
				fmt.Println(err)
				return
			}

			if err := gloo.Publish(verbose, dryRun,
				buildConfig.PublishImages, conf.WorkDir, buildConfig.ImageTag, buildConfig.DockerUser); err != nil {
				fmt.Println(err)
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
