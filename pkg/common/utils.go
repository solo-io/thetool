package common

func GetSshKeyArgs(sshkey string) []string {
	// mount user's key as read only, to guarantee it is unharmed
	args := []string{"-v", sshkey + ":/etc/user-data/ssh-keys/id_rsa:ro"}

	// prepare a tmp volume to hold the key, as we need to chmod \ chown it
	args = append(args, "--mount", "type=tmpfs,destination=/etc/github/")

	return args
}

const PrepareKeyTemplate = `
	if [ -f "/etc/user-data/ssh-keys/id_rsa" ];
	then
	  cp /etc/user-data/ssh-keys/id_rsa /etc/github/id_rsa
	  chmod 400 /etc/github/id_rsa
	  chown thetool /etc/github/id_rsa
	fi
	`
