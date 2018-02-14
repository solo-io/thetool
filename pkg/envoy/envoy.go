package envoy

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"
)

const (
	buildFile     = "BUILD"
	workspaceFile = "WORKSPACE"
	buildScript   = `#!/bin/bash

set -e
cd /source
bazel build -c dbg //:envoy
cp -f bazel-bin/envoy .

`
)

func Build(features []feature.Feature, verbose, dryRun, cache bool, eHash, wDir string) error {
	fmt.Println("Building Envoy...")
	envoyHash = eHash
	workDir = wDir
	// TODO(ashish) for each filter make sure we
	// have Envoy filter to build
	if err := generateFromTemplate(features, buildFile, buildTemplate); err != nil {
		return err
	}
	if err := generateFromTemplate(features, workspaceFile, workspaceTemplate); err != nil {
		return err
	}
	// run build in docker
	ioutil.WriteFile("build.sh", []byte(buildScript), 0755)
	if cache {
		if err := os.MkdirAll("cache/envoy", 0755); err != nil {
			return errors.Wrap(err, "unable to create cache for envoy")
		}
	}
	// docker run -t -i -v "$PWD":/source envoyproxy/envoy-build-ubuntu /bin/bash -lc "cd /source && bazel build -c dbg //:envoy"
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	var args []string
	if cache {
		args = []string{
			"run", "-i", "--rm", "-v", pwd + ":/source",
			"-v", pwd + "/cache/envoy:/root/.cache/bazel",
			"envoyproxy/envoy-build-ubuntu", "/source/build.sh",
		}
	} else {
		args = []string{
			"run", "-i", "--rm", "-v", pwd + ":/source",
			"envoyproxy/envoy-build-ubuntu", "/source/build.sh",
		}
	}
	err = util.RunCmd(verbose, dryRun, "docker", args...)
	if err != nil {
		return errors.Wrap(err, "unable to build envoy")
	}
	return nil
}

func Publish(verbose, dryRun bool, hash, user string) error {
	fmt.Println("Publishing Envoy...")

	err := ioutil.WriteFile("Dockerfile.envoy", []byte(dockerfile), 0644)
	if err != nil {
		return err
	}

	buildArgs := []string{
		"build",
		"-f", "Dockerfile.envoy",
		"-t", user + "/envoy:" + hash,
		".",
	}
	err = util.RunCmd(verbose, dryRun, "docker", buildArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to create envoy image")
	}

	pushArgs := []string{
		"push",
		user + "/envoy:" + hash,
	}
	err = util.RunCmd(verbose, dryRun, "docker", pushArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to push envoy image")
	}
	return nil
}

func generateFromTemplate(features []feature.Feature, filename string, t *template.Template) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()
	err = t.Execute(f, features)
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}
