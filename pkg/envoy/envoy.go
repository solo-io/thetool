package envoy

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

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

func Build(enabled []feature.Feature, verbose, dryRun, cache bool, eHash, wDir string) error {
	fmt.Println("Building Envoy...")
	envoyHash = eHash
	workDir = wDir
	features := envoyFilters(enabled)
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
	name := "thetool-envoy"
	var args []string
	if cache {
		args = []string{
			"run", "-i", "--rm", "--name", name,
			"-v", pwd + ":/source",
			"-v", pwd + "/cache/envoy:/root/.cache/bazel",
			"envoyproxy/envoy-build-ubuntu", "/source/build.sh",
		}
	} else {
		args = []string{
			"run", "-i", "--rm", "--name", name,
			"-v", pwd + ":/source",
			"envoyproxy/envoy-build-ubuntu", "/source/build.sh",
		}
	}

	err = util.DockerRun(verbose, dryRun, name, args...)
	if err != nil {
		var msg string
		if cache {
			msg = fmt.Sprintf("unable to build envoy; please look at %s for details", envoyBuildLog())
		} else {
			msg = "unable to build enovy; consider running in verbose mode"
		}
		return errors.Wrap(err, msg)
	}
	return nil
}

func Publish(verbose, dryRun, publish bool, imageTag, user string) error {
	fmt.Println("Publishing Envoy...")

	err := ioutil.WriteFile("Dockerfile.envoy", []byte(dockerfile), 0644)
	if err != nil {
		return err
	}

	image := user + "/envoy:" + imageTag
	buildArgs := []string{
		"build",
		"-f", "Dockerfile.envoy",
		"-t", image,
		".",
	}
	err = util.RunCmd(verbose, dryRun, "docker", buildArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to create envoy image")
	}

	if publish {
		pushArgs := []string{"push", image}
		err = util.RunCmd(verbose, dryRun, "docker", pushArgs...)
		if err != nil {
			return errors.Wrap(err, "unable to push envoy image")
		}
		fmt.Printf("Pushed Envoy image %s\n", image)
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

func envoyBuildLog() string {
	logFile := "command.log"
	bazelDir := "cache/envoy/_bazel_root"
	files, err := ioutil.ReadDir(bazelDir)
	if err == nil {
		for _, f := range files {
			if f.Name() != "install" {
				return filepath.Join(bazelDir, f.Name(), logFile)
			}
		}
	}
	return filepath.Join(bazelDir, "<hash>", logFile)
}

func envoyFilters(enabled []feature.Feature) []feature.Feature {
	out := []feature.Feature{}
	for _, f := range enabled {
		if f.EnvoyDir != "" {
			out = append(out, f)
		}
	}
	return out
}
