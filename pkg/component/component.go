package component

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/cmd/addon"
	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/util"

	"github.com/solo-io/thetool/pkg/common"
	"github.com/solo-io/thetool/pkg/config"
	"github.com/solo-io/thetool/pkg/envoy"
	"github.com/solo-io/thetool/pkg/feature"
	"github.com/solo-io/thetool/pkg/gloo"
)

type BuilderConfig struct {
	Enabled      []feature.Feature
	Verbose      bool
	DryRun       bool
	UseCache     bool
	PublishImage bool
	ImageTag     string
	DockerUser   string
	SSHKeyFile   string
	Config       *config.Config
}

type Builder struct {
	Name    string
	Builder func(BuilderConfig)
}

const (
	All = "all"
)

var (
	Builders []Builder
)

func Components() []string {
	c := make([]string, len(Builders)+1)
	for i, b := range Builders {
		c[i] = b.Name
	}
	c[len(Builders)] = All
	return c
}

func init() {
	Builders = append(Builders, Builder{
		Name: "envoy",
		Builder: func(b BuilderConfig) {
			if err := envoy.Build(b.Enabled, b.Verbose, b.DryRun, b.UseCache, b.SSHKeyFile,
				b.Config.EnvoyHash, b.Config.EnvoyCommonHash, b.Config.EnvoyRepoUser, config.WorkDir,
				b.Config.EnvoyBuilderHash); err != nil {
				fmt.Println(err)
				return
			}
			if err := envoy.Publish(b.Verbose, b.DryRun, b.PublishImage, b.ImageTag, b.DockerUser); err != nil {
				fmt.Println(err)
				return
			}
		},
	})

	Builders = append(Builders, Builder{
		Name: "gloo",
		Builder: func(b BuilderConfig) {
			if err := gloo.Build(b.Enabled, b.Verbose, b.DryRun, b.UseCache, b.SSHKeyFile,
				b.Config.GlooRepo, b.Config.GlooHash, config.WorkDir); err != nil {
				fmt.Println(err)
				return
			}

			if err := gloo.Publish(b.Verbose, b.DryRun, b.PublishImage,
				config.WorkDir, b.ImageTag, b.DockerUser); err != nil {
				fmt.Println(err)
			}
		},
	})

	addons, err := addon.List()
	if err != nil {
		// this is possible during init
		return
	}
	for _, a := range addons {
		srv := a // necessary for the way pointers work
		if isGlooAddon(srv) {
			builder := Builder{
				Name: srv.Name,
				Builder: func(b BuilderConfig) {
					if err := buildRepo(b.Verbose, b.DryRun, b.UseCache, b.SSHKeyFile,
						b.Config.GlooRepo, b.Config.GlooHash, config.WorkDir); err != nil {
						fmt.Println(err)
						return
					}

					if err := publishRepo(b.Verbose, b.DryRun, b.PublishImage,
						b.Config.GlooRepo, config.WorkDir, b.Config.GlooHash, b.DockerUser); err != nil {
						fmt.Println(err)
						return
					}
				},
			}
			Builders = append(Builders, builder)
		}
	}
}

func isGlooAddon(a *addon.Addon) bool {
	_, exists := a.Configuration["gloo"]
	return exists // only gloo addons have this key
}

func buildRepo(verbose, dryRun, useCache bool, sshKeyFile, repo, hash, workDir string) error {
	name := downloader.RepoDir(repo)
	fmt.Printf("Building %s...\n", name)
	scriptFilename := fmt.Sprintf("build-%s.sh", name)
	generateBuildScript(scriptFilename, workDir, repo)
	if err := downloader.Download(repo, hash, workDir, verbose); err != nil {
		return errors.Wrapf(err, "unable to download %s repository", name)
	}
	if !dryRun {
		if useCache {
			if err := os.MkdirAll(filepath.Join("cache", name), 0755); err != nil {
				return errors.Wrap(err, "unable to create cache directory for "+name)
			}
		}
	}
	// let's build it all in Docker
	// create output directory
	os.Mkdir(name+"-out", 0777)
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get working directory")
	}
	containerName := "thetool-" + name
	args := []string{"run", "-i", "--rm", "--name", containerName, "-v", pwd + ":/code"}

	if sshKeyFile != "" {
		args = append(args, common.GetSshKeyArgs(sshKeyFile)...)
	}

	uargs, err := common.GetUidArgs()
	if err != nil {
		return err
	}
	args = append(args, uargs...)
	args = append(args, "golang:1.10", filepath.Join("/code", scriptFilename))
	err = util.DockerRun(verbose, dryRun, name, args...)
	if err != nil {
		return errors.Wrapf(err, "unable to build %s; consider running with verbose flag", name)
	}
	return nil
}

func publishRepo(verbose, dryRun, publish bool, repo, workDir, imageTag, dockerUser string) error {
	name := downloader.RepoDir(repo)
	fmt.Printf("Publishing %s...\n", name)

	if err := util.Copy(filepath.Join(workDir, name, "Dockerfile"),
		filepath.Join(name+"-out", "Dockerfile")); err != nil {
		return errors.Wrap(err, "unable to copy the Dockerfile")
	}
	oldDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "not able to get the working directory")
	}
	if err := os.Chdir(name + "-out"); err != nil {
		return errors.Wrap(err, "unable to change the working directory to "+name+"-out")
	}
	defer os.Chdir(oldDir)

	tag := dockerUser + "/" + name + ":" + imageTag
	buildArgs := []string{
		"build",
		"-t", tag, ".",
	}
	if err := util.RunCmd(verbose, dryRun, "docker", buildArgs...); err != nil {
		return errors.Wrapf(err, "unable to create %s image", name)
	}
	if publish {
		pushArgs := []string{"push", tag}
		if err := util.RunCmd(verbose, dryRun, "docker", pushArgs...); err != nil {
			return errors.Wrapf(err, "unable to push %s image ", name)
		}
		fmt.Printf("Pushed %s image %s\n", name, tag)
	}
	return nil
}

func generateBuildScript(filename, workDir, repoURL string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.Wrap(err, "unable to create file: "+filename)
	}
	defer f.Close()
	err = buildSriptTemplate.Execute(f, map[string]string{
		"repoParent": repoParent(repoURL),
		"repoDir":    downloader.RepoDir(repoURL),
		"workDir":    workDir,
	})
	if err != nil {
		return errors.Wrap(err, "unable to write file: "+filename)
	}
	return nil
}

func repoParent(repoURL string) string {
	// TODO (ashish)
	return "github.com/solo-io"
}
