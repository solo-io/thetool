package downloader

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/solo-io/thetool/pkg/feature"
)

var (
	githubPattern = regexp.MustCompile("https://github.com/([^/]+)/(.+)")
	// based on https://github.com/bazelbuild/bazel/blob/master/tools/build_defs/repo/git.bzl
	gitTemplate = template.Must(template.New("git").Parse(`set -ex
( cd '{{.workDir}}' 
rm -rf '{{.repoDir}}'
git clone '{{.remote}}' '{{.repoDir}}'
cd '{{.repoDir}}'
git reset --hard {{.ref}} || (git fetch origin {{.ref}}:{{.ref}} && git reset --hard {{.ref}})
git clean -xdf 
git submodule update --init --checkout --force )
`))
)

// Checks to see if the URL format is supported
func SupportedURL(repoURL string) bool {
	return strings.HasSuffix(repoURL, ".git") ||
		strings.HasPrefix(repoURL, "http")
}

// Download fetches the feature from its repository and saves it to the folder
func Download(f feature.Feature, folder string, verbose bool) error {
	if strings.HasSuffix(f.Repository, ".git") {
		return withGit(f.Repository, f.Version, folder, verbose)
	}

	if strings.HasPrefix(f.Repository, "http") {
		source := handleGitHub(f.Repository, f.Version)
		srcURL, err := url.Parse(source)
		if err != nil {
			return errors.Wrap(err, "invalid repository URL")
		}
		filename := path.Base(srcURL.Path)
		destination := path.Join(folder, filename)
		err = withHTTP(handleGitHub(f.Repository, f.Version), destination)
		if err != nil {
			return err
		}
		// expand (for now we will assume everything is zip file)
		return expand(folder, filename)
	}

	return fmt.Errorf("unsupported repository scheme %s\nShould either end in '.git' or be HTTP/S URL", f.Repository)
}

func withHTTP(url, destination string) error {
	out, err := os.Create(destination)
	if err != nil {
		return errors.Wrap(err, "unable to create "+destination)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "unable to download "+url)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "unable to save "+url)
	}
	return nil
}

// withGit - uses Git SSH to download the repository
func withGit(url, commit, folder string, verbose bool) error {
	var out bytes.Buffer
	data := map[string]string{
		"workDir": folder,
		"repoDir": RepoDir(url),
		"remote":  url,
		"ref":     commit,
	}
	if err := gitTemplate.Execute(&out, data); err != nil {
		return errors.Wrap(err, "unable to create git script")
	}
	script := out.String()
	if verbose {
		fmt.Println(script)
	}

	cmd := exec.Command("bash", "-c", script)
	if verbose {
		cmdStdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.Wrap(err, "unable to create StdOut pipe for bash")
		}
		stdoutScanner := bufio.NewScanner(cmdStdout)
		go func() {
			for stdoutScanner.Scan() {
				fmt.Println("bash: " + stdoutScanner.Text())
			}
		}()

		cmdStderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.Wrap(err, "unable to create StdErr pipe for bash")
		}
		stderrScanner := bufio.NewScanner(cmdStderr)
		go func() {
			for stderrScanner.Scan() {
				fmt.Println("bash: " + stderrScanner.Text())
			}
		}()
	}
	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "unable to start cloning with git")
	}
	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "unable to clone with git")
	}
	return nil
}

func handleGitHub(repo, version string) string {
	result := githubPattern.FindAllStringSubmatch(repo, -1)
	if len(result) != 1 {
		return repo
	}
	if len(result[0]) != 3 {
		return repo
	}
	return fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip",
		result[0][1], strings.TrimSuffix(result[0][2], ".git"), version)
}

func expand(folder, filename string) error {
	destination := path.Join(folder, filename)
	r, err := zip.OpenReader(destination)
	if err != nil {
		return errors.Wrap(err, "unable to expand the zip file "+filename)
	}
	defer r.Close()

	for _, zf := range r.File {
		rc, err := zf.Open()
		if err != nil {
			return errors.Wrap(err, "unable to expand the zip file "+filename)
		}
		defer rc.Close()

		zpath := filepath.Join(folder, zf.Name)
		if zf.FileInfo().IsDir() {
			os.MkdirAll(zpath, os.ModePerm)
		} else {
			var zdir string
			if lastIndex := strings.LastIndex(zpath, string(os.PathSeparator)); lastIndex > -1 {
				zdir = zpath[:lastIndex]
			}
			err = os.MkdirAll(zdir, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "unable to expand the zip file "+filename)
			}
			f, err := os.OpenFile(zpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
			if err != nil {
				return errors.Wrapf(err, "unable to create %s while expanding %s", zf.Name, filename)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return errors.Wrapf(err, "unable to write %s while expanding %s", zf.Name, filename)
			}
		}
	}
	return nil
}

func RepoDir(remoteURL string) string {
	repo := remoteURL[strings.LastIndex(remoteURL, "/")+1:]
	ext := filepath.Ext(repo)
	if ext != "" {
		return repo[:len(repo)-len(ext)]
	}
	return repo
}
