package component

import (
	"text/template"

	"github.com/solo-io/thetool/pkg/common"
)

var (
	buildScript = `#!/bin/bash

	set -ex
	` + common.CreateUserTemplate("/code") + ` 	
` + common.PrepareKeyTemplate + `

if [ -f "/etc/github/id_rsa" ]; 
then
	export GIT_SSH_COMMAND="ssh -i /etc/github/id_rsa -o 'StrictHostKeyChecking no'"
fi

# create a script to run in su 
cat << EOF > build_user.sh
#!/bin/bash
set -ex
PATH="$PATH"
cd $GOPATH
if [ -n "$GIT_SSH_COMMAND" ]; then
	GIT_SSH_COMMAND="$GIT_SSH_COMMAND"
	git config --global url.\"git@github.com:\".insteadOf \"https://github.com\"
fi
mkdir -p -v src/{{ .repoParent }}
cd src/{{ .repoParent }}
ln -s /code/{{ .workDir }}/{{ .repoDir }} .
cd {{ .repoDir }}
pwd
go get -u github.com/golang/dep/cmd/dep
dep ensure -vendor-only
GOOS=linux CGO_ENABLED=0 go build -o {{ .repoDir }}
cp {{ .repoDir }} /code/{{ .repoDir }}-out
EOF

chmod a+rx ./build_user.sh
if [ -n "$THETOOL_UID" ]; then
su thetool -c ./build_user.sh
else
bash -c ./build_user.sh
fi
`
)

var (
	buildSriptTemplate = template.Must(template.New("build").Parse(buildScript))
)
