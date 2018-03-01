package cmd

import (
	"fmt"
	"os"

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
	flags.StringVarP(&conf.WorkDir, "work-dir", "w", config.WorkDir, "working directory")
	flags.StringVarP(&conf.EnvoyHash, "envoy-hash", "e", config.EnvoyHash, "Envoy commit hash to use")
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
	conf.GlooChartRepo = config.GlooChartRepo
	conf.GlooChartHash = config.GlooChartHash
	// Other Gloo components
	conf.GlooFuncDRepo = config.GlooFuncDiscoveryRepo
	conf.GlooFuncDHash = config.GlooFuncDiscoveryHash
	conf.GlooIngressRepo = config.GlooIngressRepo
	conf.GlooIngressHash = config.GlooIngressHash
	conf.GlooK8SDRepo = config.GlooK8SDiscvoeryRepo
	conf.GlooK8SDHash = config.GlooK8SDiscoveryHash

	if err := conf.Save(config.ConfigFile); err != nil {
		fmt.Printf("unable to save the configuration to %s: %q\n", config.ConfigFile, err)
		return
	}
	// create directory for external feature repositories
	if _, err := os.Stat(conf.WorkDir); os.IsNotExist(err) {
		err = os.Mkdir(conf.WorkDir, 0755)
		if err != nil {
			fmt.Printf("Unable to create repository directory %s: %q\n", conf.WorkDir, err)
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
	if !noDefaults {
		fmt.Println("Adding default repositories...")
		// get list of available repositories
		repos := feature.ListDefaultRepos()
		for _, r := range repos {
			if err := runAdd(verbose, r.URL, r.Commit); err != nil {
				fmt.Printf("Error setting up default repository %s: %q\n", r.URL, err)
				return
			}
		}
	}
	fmt.Println("Initialized.")
}
