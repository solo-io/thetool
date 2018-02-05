package glue

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"
)

var (
	RepositoryDirectory = "external"
)

func Build(verbose, dryRun bool, features []feature.Feature) error {
	fmt.Println("Building Glue...")
	if err := ioutil.WriteFile("build-glue.sh", []byte(buildScript), 0755); err != nil {
		return errors.Wrap(err, "unable to write build script")
	}

	if !dryRun {
		f := feature.Feature{
			Name:       "glue",
			Repository: "https://github.com/solo-io/glue.git",
			Version:    "5309cb36385555b7c2d5278fc230b2b27d8a0787",
		}
		if err := downloader.Download(f, RepositoryDirectory, verbose); err != nil {
			return errors.Wrap(err, "unable to download glue repository")
		}
	}

	// what about plugins from features?

	// let's build it all in Docker
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	args := []string{
		"run", "-i", "--rm", "-v", pwd + ":/glue",
		"golang:1.9", "/glue/build-glue.sh",
	}
	err = util.RunCmd(verbose, dryRun, "docker", args...)
	if err != nil {
		return errors.Wrap(err, "unable to build glue")
	}
	return nil
}

func Publish(verbose, dryRun bool, hash, user string) error {
	fmt.Println("Publishing Glue...")
	if err := ioutil.WriteFile("Dockerfile.glue", []byte(dockerFile), 0644); err != nil {
		return errors.Wrap(err, "unable to create Dockerfile")
	}

	tag := user + "/glue:" + hash

	buildArgs := []string{
		"build",
		"-f", "Dockerfile.glue",
		"-t", tag, ".",
	}
	if err := util.RunCmd(verbose, dryRun, "docker", buildArgs...); err != nil {
		return errors.Wrap(err, "unable to create glue image")
	}
	pushArgs := []string{
		"push", tag,
	}
	if err := util.RunCmd(verbose, dryRun, "docker", pushArgs...); err != nil {
		return errors.Wrap(err, "unable to push glue image ")
	}
	return nil
}
