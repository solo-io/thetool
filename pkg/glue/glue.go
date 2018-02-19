package glue

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"
)

func Build(enabled []feature.Feature, verbose, dryRun, cache bool, sshKeyFile, glueRepo, glueHash, workDir string) error {
	fmt.Println("Building Gloo...")

	script := fmt.Sprintf(buildScript, workDir)
	if err := ioutil.WriteFile("build-gloo.sh", []byte(script), 0755); err != nil {
		return errors.Wrap(err, "unable to write build script")
	}

	if !dryRun {
		if err := downloader.Download(glueRepo, glueHash, workDir, verbose); err != nil {
			return errors.Wrap(err, "unable to download glue repository")
		}
		if cache {
			if err := os.MkdirAll("cache/gloo", 0755); err != nil {
				return errors.Wrap(err, "unable to create cache directory for glue")
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
		if err := updateDep(plugins, df); err != nil {
			return errors.Wrapf(err, "unable to update to dependencies file %s", df)
		}
	}
	// let's build it all in Docker
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	name := "thetool-gloo"
	args := []string{"run", "-i", "--rm", "--name", name, "-v", pwd + ":/gloo"}
	if cache {
		args = append(args, "-v", pwd+"/cache/gloo:/go/pkg/dep/sources")
	}

	if sshKeyFile != "" {
		args = append(args, "-v", sshKeyFile+":/etc/github/id_rsa")
	}
	args = append(args, "golang:1.9", "/gloo/build-gloo.sh")
	err = util.DockerRun(verbose, dryRun, name, args...)
	if err != nil {
		return errors.Wrap(err, "unable to build gloo; consider running with verbose flag")
	}
	return nil
}

func Publish(verbose, dryRun, publish bool, workDir, imageTag, user string) error {
	fmt.Println("Publishing Gloo...")

	tag := user + "/gloo:" + imageTag
	buildArgs := []string{
		"build",
		"-f", filepath.Join(workDir, "gloo", "Dockerfile"),
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

func updateDep(plugins []GlooPlugin, filename string) error {
	// get unique Repositories
	repos := make(map[string]string)
	for _, p := range plugins {
		repos[getPackage(p.Repository)] = p.Revision
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
