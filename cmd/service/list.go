package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list supporting services",
		Run: func(c *cobra.Command, args []string) {
			runList()
		},
	}
	return cmd
}

func runList() {
	services, err := load(serviceFilename)
	if err != nil {
		fmt.Printf("Unable to load services %q\n", err)
		return
	}
	for _, s := range services {
		fmt.Println(s)
	}
}
