package component

import "html/template"

const (
	buildScript = `#!/bin/bash

	set -ex
	
	if [ -f "/etc/github/id_rsa" ]; then
	  chmod 400 /etc/github/id_rsa
	  export GIT_SSH_COMMAND="ssh -i /etc/github/id_rsa -o 'StrictHostKeyChecking no'"
	  git config --global url."git@github.com:".insteadOf "https://github.com"
	fi
	
	cd $GOPATH
	mkdir -p -v src/{{ .repoParent }}
	cd src/{{ .repoParent }}
	ln -s /code/{{ .workDir }}/{{ .repoDir }} .
	cd {{ .repoDir }}
	pwd
	
	go get -u github.com/golang/dep/cmd/dep
	
	dep ensure -vendor-only
	GOOS=linux CGO_ENABLED=0 go build -o {{ .repoDir }}
	cp {{ .repoDir }} /code/{{ .repoDir }}-out
	`
)

var (
	buildSriptTemplate = template.Must(template.New("build").Parse(buildScript))
)
