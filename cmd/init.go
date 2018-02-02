package cmd

import (
	"fmt"
	"os"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

const (
	dataFile      = "features.json"
	repositoryDir = "external"
)

// InitCmd initialize current directory for thetool
func InitCmd() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "init",
		Short: "intialize the tool",
		RunE: func(c *cobra.Command, args []string) error {
			return runInit(verbose)
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	return cmd
}

func runInit(verbose bool) error {
	fmt.Println("Initializing current directory...")
	// check if this directory is already initialized
	if _, err := os.Stat(dataFile); err == nil {
		// TODO(ashish) check it is thetool file
		return fmt.Errorf("thetool already initialized")
	}
	// create directory for external feature repositories
	if _, err := os.Stat(repositoryDir); os.IsNotExist(err) {
		err = os.Mkdir(repositoryDir, 0755)
		if err != nil {
			return err
		}
	}

	fmt.Println("Adding default features...")
	if err := feature.SaveToFile([]feature.Feature{}, dataFile); err != nil {
		return err
	}
	// get list of available filters
	features := feature.ListDefaultFeatures()
	for _, f := range features {
		if err := runAdd(f.Name, f.Repository, f.Version, verbose); err != nil {
			return err
		}
	}

	fmt.Println("Initialized.")
	return nil
}
