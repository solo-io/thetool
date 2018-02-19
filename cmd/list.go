package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
)

func ListReposCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-repo",
		Short: "list all registered Gloo repositories",
		Run: func(c *cobra.Command, args []string) {
			runListRepos()
		},
	}
	return cmd
}

func runListRepos() {
	store := &feature.FileRepoStore{Filename: feature.ReposFileName}
	repos, err := store.List()
	if err != nil {
		fmt.Printf("Unable to load repository list: %q\n", err)
		return
	}
	if len(repos) == 0 {
		fmt.Println("No repositories added yet!")
	}
	for _, r := range repos {
		fmt.Println("Repository: ", r.URL)
		fmt.Println("Commit:     ", r.Commit)
		fmt.Println("")
	}
}
