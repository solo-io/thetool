package cmd

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/spf13/cobra"
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

func BuildCmd() *cobra.Command {
	var verbose bool
	var dryRun bool
	var dockerUser string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build the universe",
		RunE: func(c *cobra.Command, args []string) error {
			return runBuild(verbose, dryRun, dockerUser)
		},
	}
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "show verbose build log")
	cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run; only generate build file")
	cmd.PersistentFlags().StringVarP(&dockerUser, "docker-user", "u", "", "Docker user for publishing images")
	return cmd
}

func runBuild(verbose, dryRun bool, dockerUser string) error {
	if !dryRun && dockerUser == "" {
		return fmt.Errorf("need Docker user ID to publish images")
	}
	features, err := feature.LoadFromFile("features.json")
	if err != nil {
		return err
	}
	var enabled []feature.Feature
	for _, f := range features {
		if f.Enabled {
			enabled = append(enabled, f)
		}
	}
	fmt.Printf("Building with %d features\n", len(enabled))
	featuresHash := featuresHash(enabled)
	err = buildEnvoy(enabled, verbose, dryRun)
	if err != nil {
		return err
	}
	err = publishEnvoy(verbose, dryRun, featuresHash, dockerUser)
	if err != nil {
		return err
	}

	buildGlue(enabled)
	publishGlue()

	generateHelmChart()
	return nil
}

func buildEnvoy(features []feature.Feature, verbose, dryRun bool) error {
	fmt.Println("Building Envoy...")
	// TODO(ashish) for each filter make sure we
	// have Envoy filter to build
	if err := generateFromTemplate(features, buildFile, buildTemplate); err != nil {
		return err
	}
	if err := generateFromTemplate(features, workspaceFile, workspaceTemplate); err != nil {
		return err
	}
	// run build in docker
	generateBuildSh()
	// docker run -t -i -v "$PWD":/source envoyproxy/envoy-build-ubuntu /bin/bash -lc "cd /source && bazel build -c dbg //:envoy"
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	args := []string{
		"run", "-i", "--rm", "-v", pwd + ":/source",
		"envoyproxy/envoy-build-ubuntu", "/source/build.sh",
	}
	err = runCmd(verbose, dryRun, "docker", args...)
	if err != nil {
		return errors.Wrap(err, "unable to build envoy")
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

func generateBuildSh() error {
	err := writeFile("build.sh", buildScript)
	if err != nil {
		return err
	}

	if err = os.Chmod("build.sh", 0755); err != nil {
		return errors.Wrap(err, "unable to make the build script executable")
	}
	return nil
}

func buildGlue(features []feature.Feature) error {
	fmt.Println("Building Glue...")
	return fmt.Errorf("not implemented")
}

func publishEnvoy(verbose, dryRun bool, hash, user string) error {
	fmt.Println("Publishing Envoy...")

	err := writeFile("Dockerfile.envoy", dockerfile)
	if err != nil {
		return err
	}

	buildArgs := []string{
		"build",
		"-f", "Dockerfile.envoy",
		"-t", user + "/envoy:" + hash,
		".",
	}
	err = runCmd(verbose, dryRun, "docker", buildArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to create envoy image")
	}

	pushArgs := []string{
		"push",
		user + "/envoy:" + hash,
	}
	err = runCmd(verbose, dryRun, "docker", pushArgs...)
	if err != nil {
		return errors.Wrap(err, "unable to push envoy image")
	}
	return nil
}

func publishGlue() error {
	fmt.Println("Publishing Glue...")
	return fmt.Errorf("not implemented")
}

func generateHelmChart() {

}

func writeFile(filename, content string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "unable to create file %s", filename)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return errors.Wrapf(err, "unable to write file %s", filename)
	}
	return nil
}

// featuresHash generates a hash for particular envoy and glue build
// based on the features included
func featuresHash(features []feature.Feature) string {
	hash := sha256.New()
	for _, f := range features {
		hash.Write([]byte(f.Name))
		hash.Write([]byte(f.Version))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func runCmd(verbose, dryRun bool, binary string, args ...string) error {
	if verbose {
		fmt.Println(binary, args)
	}
	if dryRun {
		return nil
	}

	cmd := exec.Command(binary, args...)
	if verbose {
		cmdStdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrapf(err, "unable to create StdOut pipe for %s", binary)
		}
		stdoutScanner := bufio.NewScanner(cmdStdout)
		go func() {
			prefix := binary + ": "
			for stdoutScanner.Scan() {
				fmt.Println(prefix, stdoutScanner.Text())
			}
		}()

		cmdStderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrapf(err, "unable to create StdErr pipe for %s", binary)
		}
		stderrScanner := bufio.NewScanner(cmdStderr)
		go func() {
			prefix := binary + ": "
			for stderrScanner.Scan() {
				fmt.Println(prefix, stderrScanner.Text())
			}
		}()
	}
	err := cmd.Start()
	if err != nil {
		return errors.Wrapf(err, "unable to start %s", binary)
	}
	err = cmd.Wait()
	if err != nil {
		return errors.Wrapf(err, "unable to run %s", binary)
	}
	return nil
}
