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

func Build(features []feature.Feature, verbose, dryRun, cache bool, glueRepo, glueHash, workDir string) error {
	fmt.Println("Building Glue...")
	script := fmt.Sprintf(buildScript, workDir)
	if err := ioutil.WriteFile("build-glue.sh", []byte(script), 0755); err != nil {
		return errors.Wrap(err, "unable to write build script")
	}

	if !dryRun {
		f := feature.Feature{
			Name:       "glue",
			Repository: glueRepo,
			Version:    glueHash,
		}
		if err := downloader.Download(f, workDir, verbose); err != nil {
			return errors.Wrap(err, "unable to download glue repository")
		}
		if cache {
			if err := os.MkdirAll("cache/glue", 0755); err != nil {
				return errors.Wrap(err, "unable to create cache directory for glue")
			}
		}
	}

	// what about plugins from features?

	// let's build it all in Docker
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	var args []string
	if cache {
		args = []string{
			"run", "-i", "--rm", "-v", pwd + ":/glue",
			"-v", pwd + "/cache/glue:/go/pkg/dep/sources",
			"golang:1.9", "/glue/build-glue.sh",
		}
	} else {
		args = []string{
			"run", "-i", "--rm", "-v", pwd + ":/glue",
			"golang:1.9", "/glue/build-glue.sh",
		}
	}
	err = util.RunCmd(verbose, dryRun, "docker", args...)
	if err != nil {
		return errors.Wrap(err, "unable to build glue; consider running with verbose flag")
	}
	return nil
}

func Publish(verbose, dryRun, publish bool, imageTag, user string) error {
	fmt.Println("Publishing Glue...")
	if err := ioutil.WriteFile("Dockerfile.glue", []byte(dockerFile), 0644); err != nil {
		return errors.Wrap(err, "unable to create Dockerfile")
	}

	tag := user + "/glue:" + imageTag
	buildArgs := []string{
		"build",
		"-f", "Dockerfile.glue",
		"-t", tag, ".",
	}
	if err := util.RunCmd(verbose, dryRun, "docker", buildArgs...); err != nil {
		return errors.Wrap(err, "unable to create glue image")
	}
	if publish {
		pushArgs := []string{"push", tag}
		if err := util.RunCmd(verbose, dryRun, "docker", pushArgs...); err != nil {
			return errors.Wrap(err, "unable to push glue image ")
		}
		fmt.Printf("Pushed Glue image %s\n", tag)
	}
	return nil
}
