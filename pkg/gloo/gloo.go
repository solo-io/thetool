package gloo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/common"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"
)

func Build(enabled []feature.Feature, verbose, dryRun, cache bool, sshKeyFile, glooRepo, glooHash, workDir string) error {
	fmt.Println("Building Gloo...")

	script := fmt.Sprintf(buildScript, workDir)
	if err := ioutil.WriteFile("build-gloo.sh", []byte(script), 0755); err != nil {
		return errors.Wrap(err, "unable to write build script")
	}

	if err := downloader.Download(glooRepo, glooHash, workDir, verbose); err != nil {
		return errors.Wrap(err, "unable to download gloo repository")
	}
	if cache {
		if err := os.MkdirAll("cache/gloo", 0777); err != nil {
			return errors.Wrap(err, "unable to create cache directory for gloo")
		}
	}

	plugins := toGlooPlugins(enabled)

	fmt.Println("Adding plugins to Gloo...")
	pf := filepath.Join(workDir, installFile)
	if err := installPlugins(plugins,
		pf, installTemplate); err != nil {
		return errors.Wrapf(err, "unable to update %s", pf)
	}

	fmt.Println("Constraining plugins to given revisions...")
	df := filepath.Join(workDir, dependencyFile)
	if err := updateDep(plugins, df, glooRepo); err != nil {
		return errors.Wrapf(err, "unable to update to dependencies file %s", df)
	}
	// let's build it all in Docker
	// create output directory
	os.Mkdir("gloo-out", 0777)
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	name := "thetool-gloo"
	args := []string{"run", "-i", "--rm", "--name", name, "-v", pwd + ":/gloo"}
	if cache {
		gloocache := filepath.Join(pwd, "cache", "gloo")
		// create it first to make sure it's with the current user.
		os.MkdirAll(gloocache, 0755)

		args = append(args, "-v", gloocache+":/go/pkg/dep/sources")
	}

	if sshKeyFile != "" {
		args = append(args, common.GetSshKeyArgs(sshKeyFile)...)
	}
	uargs, err := common.GetUidArgs()
	if err != nil {
		// doesn't return current user in Jenkins
		fmt.Println("warning: unable to get current user id:", err)
	} else {
		args = append(args, uargs...)
	}

	args = append(args, "golang:1.10", "/gloo/build-gloo.sh")
	err = util.DockerRun(verbose, dryRun, name, args...)
	if err != nil {
		return errors.Wrap(err, "unable to build gloo; consider running with verbose flag")
	}
	return nil
}

func Publish(verbose, dryRun, publish bool, workDir, imageTag, user string) error {
	fmt.Println("Publishing Gloo...")

	if !dryRun {
		if err := util.Copy(filepath.Join(workDir, "gloo", "cmd", "control-plane", "Dockerfile"), filepath.Join("gloo-out", "Dockerfile")); err != nil {
			return errors.Wrap(err, "not able to copy the Dockerfile")
		}

		oldDir, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "not able to get working directory")
		}
		if err := os.Chdir("gloo-out"); err != nil {
			return errors.Wrap(err, "unable to change working directory to gloo-out")
		}
		defer os.Chdir(oldDir)
	}
	tag := user + "/control-plane:" + imageTag
	buildArgs := []string{
		"build",
		"-t", tag, ".",
	}
	if err := util.RunCmd(verbose, dryRun, "docker", buildArgs...); err != nil {
		return errors.Wrap(err, "unable to create gloo image")
	}
	if publish {
		pushArgs := []string{"push", tag}
		if err := util.RunCmd(verbose, dryRun, "docker", pushArgs...); err != nil {
			return errors.Wrap(err, "unable to push gloo image ")
		}
		fmt.Printf("Pushed Gloo image %s\n", tag)
	}
	return nil
}

func installPlugins(packages []GlooPlugin, filename string, t *template.Template) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()
	err = t.Execute(f, packages)
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}

func updateDep(plugins []GlooPlugin, filename, glooRepo string) error {
	// get unique Repositories
	repos := make(map[string]string)
	for _, p := range plugins {
		if p.Repository != glooRepo {
			repos[getPackage(p.Repository)] = p.Revision
		}
	}
	var w bytes.Buffer
	if err := packageTemplate.Execute(&w, repos); err != nil {
		return errors.Wrap(err, "unable to generate dependency constraints")
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to open depencency file for update")
	}
	defer f.Close()

	if _, err = f.Write(w.Bytes()); err != nil {
		return errors.Wrap(err, "unable to update the dependency files")
	}

	return nil
}
