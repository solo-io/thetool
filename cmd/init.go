package cmd

import (
	"fmt"
	"os"

	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

const (
	dataFile = "features.json"
)

// InitCmd initialize current directory for thetool
func InitCmd() *cobra.Command {
	var verbose bool
	conf := config.Config{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "intialize the tool",
		Run: func(c *cobra.Command, args []string) {
			runInit(verbose, conf)
		},
	}
	pflags := cmd.PersistentFlags()
	pflags.BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	pflags.StringVarP(&conf.WorkDir, "work-dir", "w", config.WorkDir, "working directory")
	pflags.StringVarP(&conf.EnvoyHash, "envoy-hash", "e", config.EnvoyHash, "Envoy commit hash to use")
	pflags.StringVarP(&conf.GlueHash, "glue-hash", "g", config.GlueHash, "Glue commit hash to use")
	pflags.StringVarP(&conf.DockerUser, "user", "u", config.DockerUser, "default Docker user")
	return cmd
}

func runInit(verbose bool, conf config.Config) {
	fmt.Println("Initializing current directory...")
	// check if this directory is already initialized
	if _, err := os.Stat(dataFile); err == nil {
		fmt.Println("thetool already initialized")
		return
	}

	// Let's save the configuration file
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

	fmt.Println("Adding default features...")
	if err := feature.SaveToFile([]feature.Feature{}, dataFile); err != nil {
		fmt.Printf("Unable to save features file %s: %q\n", dataFile, err)
		return
	}
	// get list of available filters
	features := feature.ListDefaultFeatures()
	for _, f := range features {
		if err := runAdd(f.Name, f.Repository, f.Version, verbose); err != nil {
			fmt.Printf("Error setting up default feature %s: %q\n", f.Name, err)
			return
		}
	}

	fmt.Println("Initialized.")
}
