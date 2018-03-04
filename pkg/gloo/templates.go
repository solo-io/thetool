package gloo

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/solo-io/thetool/pkg/common"
	"github.com/solo-io/thetool/pkg/feature"
)

const (
	buildScript = `#!/bin/bash

set -ex

chmod 777 $GOPATH/pkg/dep

` + common.CreateUserTemplate + `
` + common.PrepareKeyTemplate + `


if [ -f "/etc/github/id_rsa" ]; 
then
  export GIT_SSH_COMMAND="ssh -i /etc/github/id_rsa -o 'StrictHostKeyChecking no'"
  su thetool -c "PATH=\"$PATH\" && GIT_SSH_COMMAND=\"$GIT_SSH_COMMAND\" && git config --global url.\"git@github.com:\".insteadOf \"https://github.com\" &&cd $GOPATH && mkdir -p -v src/github.com/solo-io && cd src/github.com/solo-io && ln -s /gloo/%s/gloo . && cd gloo && pwd && go get -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only && GOOS=linux CGO_ENABLED=0 go build -o gloo && cp gloo /gloo/gloo-out"
else
  su thetool -c "PATH=\"$PATH\" && cd $GOPATH && mkdir -p -v src/github.com/solo-io && cd src/github.com/solo-io && ln -s /gloo/%s/gloo . && cd gloo && pwd && go get -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only && GOOS=linux CGO_ENABLED=0 go build -o gloo && cp gloo /gloo/gloo-out"
fi


`
	installGo = `package install

import (
{{ range .}}
	_ "{{.Package}}"
{{end}}
)`

	installFile = "gloo/internal/install/install_plugins.go"

	gopkg = `{{range $k, $v := .}}
[[constraint]]
  name = "{{$k}}"
  revision = "{{$v}}"
{{end}}`

	dependencyFile = "gloo/Gopkg.toml"
)

var (
	installTemplate *template.Template
	packageTemplate *template.Template
)

func init() {
	installTemplate = template.Must(template.New("install").Parse(installGo))
	packageTemplate = template.Must(template.New("package").Parse(gopkg))
}

type GlooPlugin struct {
	Package    string
	Revision   string
	Repository string
}

func toGlooPlugins(features []feature.Feature) []GlooPlugin {
	plugins := []GlooPlugin{}
	for _, f := range features {
		if f.GlooDir != "" {
			p := GlooPlugin{
				Package:    filepath.Join(getPackage(f.Repository), f.GlooDir),
				Revision:   f.Revision,
				Repository: f.Repository,
			}
			plugins = append(plugins, p)
		}
	}
	return plugins
}

func getPackage(repo string) string {
	pluginPackage := strings.Replace(repo, "https://", "", 1)
	pluginPackage = strings.Replace(pluginPackage, "http://", "", 1)

	atIndex := strings.Index(pluginPackage, "@")
	if atIndex >= 0 {
		pluginPackage = pluginPackage[atIndex+1:]
	}

	if strings.HasSuffix(pluginPackage, ".git") {
		pluginPackage = pluginPackage[:len(pluginPackage)-len(".git")]
	}
	return pluginPackage
}
