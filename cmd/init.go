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
		Run: func(c *cobra.Command, args []string) {
			runInit(verbose)
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	return cmd
}

func runInit(verbose bool) {
	fmt.Println("Initializing current directory...")
	// check if this directory is already initialized
	if _, err := os.Stat(dataFile); err == nil {
		fmt.Println("thetool already initialized")
		return
	}
	// create directory for external feature repositories
	if _, err := os.Stat(repositoryDir); os.IsNotExist(err) {
		err = os.Mkdir(repositoryDir, 0755)
		if err != nil {
			fmt.Printf("Unable to create repository directory %s: %q\n", repositoryDir, err)
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
