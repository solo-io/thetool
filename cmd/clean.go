package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func CleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "clean up files and directory",
		Run: func(c *cobra.Command, args []string) {
			runClean()
		},
	}
	return cmd
}

func runClean() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("Unable to list directory:", err)
		return
	}

	for _, f := range files {
		if shouldDelete(f) {
			if err := os.RemoveAll(f.Name()); err != nil {
				if !os.IsNotExist(err) {
					fmt.Printf("Unable to delete %v: %q\n", f.Name(), err)
				}
			}
		}
	}
}

func shouldDelete(f os.FileInfo) bool {
	name := f.Name()
	if strings.HasPrefix(name, "build-") && strings.HasSuffix(name, ".sh") {
		return true
	}

	if strings.HasSuffix(name, "-out") && f.IsDir() {
		return true
	}

	switch name {
	case "gloo-chart.yaml":
		return true
	case "cache":
		return true
	case "envoy":
		return true
	default:
		return false
	}
}
