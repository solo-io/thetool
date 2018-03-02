package envoy

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/util"
)

const (
	buildFile     = "BUILD"
	workspaceFile = "WORKSPACE"
)

func Build(enabled []feature.Feature, verbose, dryRun, cache bool, sshKeyFile, eHash, wDir, buildContainerHash string) error {
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
	// create output directory
	os.Mkdir("envoy-out", 0777)

	// run build in docker
	ioutil.WriteFile("build-envoy.sh", []byte(fmt.Sprintf(buildScript, envoyHash)), 0755)
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
	args := []string{"run", "-i", "--rm", "--name", name,
		"-v", pwd + ":/source"}
	if cache {
		args = append(args, "-v", pwd+"/cache/envoy:/source/.cache/bazel")
	}
	if sshKeyFile != "" {
		args = append(args, "-v", sshKeyFile+":/etc/github/id_rsa")
	}

	u, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "unable to get current user")
	}
	args = append(args, "--env", "THETOOL_UID="+u.Uid, "--env", "THETOOL_GID="+u.Gid)
	args = append(args, "envoyproxy/envoy-build-ubuntu@sha256:"+buildContainerHash, "/source/build-envoy.sh")

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

	err := ioutil.WriteFile("envoy-out/Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		return err
	}

	image := user + "/envoy:" + imageTag
	buildArgs := []string{
		"build",
		"-t", image,
		".",
	}
	oldDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "not able to get working directory")
	}
	if err := os.Chdir("envoy-out"); err != nil {
		return errors.Wrap(err, "unable to change working directory to envoy-out")
	}
	defer os.Chdir(oldDir)
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
	bazelDir := "cache/envoy/_bazel_thetool"
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
