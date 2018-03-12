package addon

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list addons",
		Run: func(c *cobra.Command, args []string) {
			runList()
		},
	}
	return cmd
}

func runList() {
	addons, err := load(addonFilename)
	if err != nil {
		fmt.Printf("Unable to load addons %q\n", err)
		return
	}
	for _, a := range addons {
		fmt.Println(a)
	}
}
