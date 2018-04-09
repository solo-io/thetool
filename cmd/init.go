package cmd

import (
	"fmt"
	"os"

	"github.com/solo-io/thetool/cmd/addon"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

// InitCmd initialize current directory for thetool
func InitCmd() *cobra.Command {
	var verbose bool
	var noDefaults bool
	conf := config.Config{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "intialize the tool",
		Run: func(c *cobra.Command, args []string) {
			runInit(verbose, noDefaults, conf)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	flags.StringVar(&conf.EnvoyRepoUser, "envoy-repo-user", config.EnvoyRepoUser, "Envoy repository user")
	flags.StringVarP(&conf.EnvoyHash, "envoy-hash", "e", config.EnvoyHash, "Envoy commit hash to use")
	flags.StringVar(&conf.EnvoyCommonHash, "envoy-common-hash", config.EnvoyCommonHash, "Hash for Soloio Envoy Common")
	flags.StringVarP(&conf.GlooHash, "gloo-hash", "g", config.GlooHash, "Gloo commit hash to use")
	flags.StringVar(&conf.GlooRepo, "gloo-repo", config.GlooRepo, "Gloo git repository")
	flags.StringVarP(&conf.DockerUser, "user", "u", config.DockerUser, "default Docker user")

	flags.BoolVar(&noDefaults, "no-defaults", false, "do not add default features")
	return cmd
}

func runInit(verbose, noDefaults bool, conf config.Config) {
	fmt.Println("Initializing current directory...")
	// check if this directory is already initialized
	if _, err := os.Stat(feature.ReposFileName); err == nil {
		fmt.Println("thetool already initialized")
		return
	}

	// Let's save the configuration file that aren't changed via CLI args
	conf.EnvoyBuilderHash = config.EnvoyBuilderHash

	if err := conf.Save(config.ConfigFile); err != nil {
		fmt.Printf("unable to save the configuration to %s: %q\n", config.ConfigFile, err)
		return
	}
	// create directory for external feature repositories
	if _, err := os.Stat(config.WorkDir); os.IsNotExist(err) {
		err = os.Mkdir(config.WorkDir, 0755)
		if err != nil {
			fmt.Printf("Unable to create repository directory %s: %q\n", config.WorkDir, err)
			return
		}
	}

	repoStore := feature.FileRepoStore{Filename: feature.ReposFileName}
	if err := repoStore.Init(); err != nil {
		fmt.Printf("Unable to initialize repositories file %s: %q\n", feature.ReposFileName, err)
		return
	}

	featureStore := feature.FileFeatureStore{Filename: feature.FeaturesFileName}
	if err := featureStore.Init(); err != nil {
		fmt.Printf("Unable to initialize features file %s: %q\n", feature.FeaturesFileName, err)
		return
	}

	if err := addon.Init(); err != nil {
		fmt.Printf("Unable to initialize supporting addons file: %q\n", err)
		return
	}
	if !noDefaults {
		fmt.Println("Adding default repositories...")
		// add the plugins in Gloo as default features
		if err := runAdd(verbose, conf.GlooRepo, conf.GlooHash, "pkg/plugins/features.json"); err != nil {
			fmt.Printf("Error setting up default features: %q\n", err)
			return
		}
	}
	fmt.Println("Initialized.")
}
