package glue

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
)

var (
	RepositoryDirectory = "external"
)

func Build(verbose, dryRun bool, features []feature.Feature) error {
	fmt.Println("Building Glue...")
	f := feature.Feature{
		Name:       "glue",
		Repository: "https://github.com/solo-io/glue.git",
		Version:    "5309cb36385555b7c2d5278fc230b2b27d8a0787",
	}
	if err := downloader.Download(f, RepositoryDirectory, verbose); err != nil {
		return errors.Wrap(err, "unable to download glue repository")
	}

	// what about plugins from features?

	// let's build it all in Docker
	return fmt.Errorf("not implemented")
}

func Publish(verbose, dryRun bool, hash, user string) error {
	fmt.Println("Publishing Glue...")
	return fmt.Errorf("not implemented")
}
