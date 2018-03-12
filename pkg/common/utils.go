package common

import (
	"os/user"

	"github.com/pkg/errors"
)

func GetUidArgs() ([]string, error) {

	u, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get current user")
	}
	return []string{"--env", "THETOOL_UID=" + u.Uid, "--env", "THETOOL_GID=" + u.Gid}, nil
}

func GetSshKeyArgs(sshkey string) []string {
	// mount user's key as read only, to guarantee it is unharmed
	args := []string{"-v", sshkey + ":/etc/user-data/ssh-keys/id_rsa:ro"}

	// prepare a tmp volume to hold the key, as we need to chmod \ chown it
	args = append(args, "--mount", "type=tmpfs,destination=/etc/github/")

	return args
}

func CreateUserTemplate(homedir string) string {
	return `
if [ -n "$THETOOL_UID" ]; then
groupadd --gid $THETOOL_GID -f thetoolgroup
useradd -o --uid $THETOOL_UID --gid $THETOOL_GID --no-create-home --home-dir ` + homedir + ` thetool
fi

`
}

const PrepareKeyTemplate = `
	if [ -f "/etc/user-data/ssh-keys/id_rsa" ];
	then
	  cp /etc/user-data/ssh-keys/id_rsa /etc/github/id_rsa
	  chmod 400 /etc/github/id_rsa
	  chown thetool /etc/github/id_rsa
	fi
	`
