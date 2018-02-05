package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/envoy"
	"github.com/solo-io/thetool/pkg/feature"
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
		Use:   "build",
		Short: "build the universe",
		RunE: func(c *cobra.Command, args []string) error {
			target := componentAll
			if len(args) > 0 {
				switch strings.ToLower(args[0]) {
				case "envoy":
					target = componentEnvoy
				case "glue":
					target = componentGlue
				default:
					target = componentAll
				}
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
	features, err := feature.LoadFromFile(dataFile)
	if err != nil {
		return err
	}
	var enabled []feature.Feature
	for _, f := range features {
		if f.Enabled {
			enabled = append(enabled, f)
		}
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
		if err := glue.Build(enabled, verbose, dryRun, conf.GlueHash, conf.WorkDir); err != nil {
			fmt.Println(err)
			return
		}

		if err := glue.Publish(verbose, dryRun, featuresHash, dockerUser); err != nil {
			fmt.Println(err)
		}
	}()

	generateHelmValues(verbose, featuresHash, dockerUser)
	wg.Wait()
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

// featuresHash generates a hash for particular envoy and glue build
// based on the features included
func featuresHash(features []feature.Feature) string {
	hash := sha256.New()
	for _, f := range features {
		hash.Write([]byte(f.Name))
		hash.Write([]byte(f.Version))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))[:8]
}
