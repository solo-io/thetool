package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/envoy"
	"github.com/solo-io/thetool/pkg/glue"
	"github.com/spf13/cobra"
)

type component int

const (
	componentAll component = iota
	componentEnvoy
	componentGlue
)

func BuildCmd() *cobra.Command {
	var verbose bool
	var dryRun bool
	var dockerUser string
	cmd := &cobra.Command{
		Use:       "build [target to build]",
		Short:     "build the universe",
		Long:      "build glue, envoy or all",
		ValidArgs: []string{"envoy", "glue", "all"},
		Args:      cobra.OnlyValidArgs,
		RunE: func(c *cobra.Command, args []string) error {
			target := componentAll
			if len(args) != 1 {
				return fmt.Errorf("please specify a build target")
			}
			switch strings.ToLower(args[0]) {
			case "envoy":
				target = componentEnvoy
			case "glue":
				target = componentGlue
			default:
				target = componentAll
			}
			return runBuild(verbose, dryRun, dockerUser, target)
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show verbose build log")
	cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run; only generate build file")
	cmd.PersistentFlags().StringVarP(&dockerUser, "docker-user", "u", "", "Docker user for publishing images")
	return cmd
}

func runBuild(verbose, dryRun bool, dockerUser string, target component) error {
	conf, err := config.Load(config.ConfigFile)
	if err != nil {
		fmt.Printf("Unable to load configuration from %s: %q\n", config.ConfigFile, err)
		return nil
	}
	if dockerUser == "" {
		dockerUser = conf.DockerUser
	}
	if dockerUser == "" {
		return fmt.Errorf("need Docker user ID to publish images")
	}
	enabled, err := loadEnabledFeatures()
	if err != nil {
		return err
	}
	fmt.Printf("Building with %d features\n", len(enabled))
	featuresHash := featuresHash(enabled)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if target != componentAll && target != componentEnvoy {
			return
		}
		if err := envoy.Build(enabled, verbose, dryRun, conf.EnvoyHash, conf.WorkDir); err != nil {
			fmt.Println(err)
			return
		}
		if err := envoy.Publish(verbose, dryRun, featuresHash, dockerUser); err != nil {
			fmt.Println(err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		if target != componentAll && target != componentGlue {
			return
		}
		if err := glue.Build(enabled, verbose, dryRun, conf.GlueRepo, conf.GlueHash, conf.WorkDir); err != nil {
			fmt.Println(err)
			return
		}

		if err := glue.Publish(verbose, dryRun, featuresHash, dockerUser); err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()
	return nil
}
