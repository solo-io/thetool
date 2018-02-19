# thetool
Build the Gloo universe

## Installing
### Prerequisite
`thetool` uses Docker to run the build process and build Docker images. Please make sure you have Docker installed. For more information please visit [Install Docker](https://docs.docker.com/install/)

### Downloading and Installing
Download the latest release from https://github.com/solo-io/thetool/releases/latest/

If you prefer to compile your own binary or work on the development `thetool` please use the following command:

```
go get github.com/solo-io/thetool
``` 

## Getting Started
Create a working directory for `thetool`

```
mkdir gloo
cd gloo
```

Initialize `thetool` with default set of Gloo features. Optionally, you specify a default user id for Docker.

```
thetool init -u solo-io
```

You can look at the default set of features by using the `list` command.

```
thetool list

Name:        squash
Repository:  https://github.com/axhixh/envoy-squash.git
Commit:      9397617b238cc4f17a0a3f0dc24194baf506ac97
Enabled:     true

Name:        echo
Repository:  https://github.com/axhixh/echo.git
Commit:      37a53fefe0a267fe3f4704c35e3721a4b6032f2a
Enabled:     true

Name:        lambda
Repository:  git@github.com:solo-io/glue-lambda.git
Commit:      9aeb7747286c9116d9f531f7cc2c3331a8a23c7f
Enabled:     true
```

You can enable or disable any of the features calling `enable` or `disable` command with the name of the feature.

```
thetool disable -n echo
```

If you list again, you will see the `echo` feature is disabled.

```
thetool list

Name:        squash
Repository:  https://github.com/axhixh/envoy-squash.git
Commit:      9397617b238cc4f17a0a3f0dc24194baf506ac97
Enabled:     true

Name:        echo
Repository:  https://github.com/axhixh/echo.git
Commit:      37a53fefe0a267fe3f4704c35e3721a4b6032f2a
Enabled:     false

Name:        lambda
Repository:  git@github.com:solo-io/glue-lambda.git
Commit:      9aeb7747286c9116d9f531f7cc2c3331a8a23c7f
Enabled:     true
```

Once you have selected the features you want to include, you can build Gloo and its components using the `build` command.

```
thetool build all
```

TODO: show how docker images are built

You can also choose to build individual components of Gloo by specifying the name of the component like `envoy` or `gloo`.

The build command builds the appropriate binaries and their corresponding Docker images. It then publishes these images to Docker registry. If you do not want to publish, you can pass a flag to `thetool`

```
thetool build all --publish=false
```

Note: in order to deploy Gloo to Kubernetes, you need to publish the Docker images.

TODO:

```
thetool deploy
deploy the universe

Usage:
  thetool deploy [command]

Available Commands:
  k8s         deploy the universe in Kubernetes
  k8s-out     deploy out of Kubernetes cluster
  local       deploy the universe locally

Flags:
  -u, --docker-user string   Docker user for publishing images
  -d, --dry-run              dry run; only generate build file
  -h, --help                 help for deploy
  -v, --verbose              show verbose build log

Use "thetool deploy [command] --help" for more information about a command.
```

## Adding Your Own Feature

`thetool` can build Gloo with your custom Gloo features by adding your feature's repository to list of features.

You can add or remove your custom feature from the features list using `add` and `delete` commands.

```
thetool add -n magic -r https://github.com/axhixh/gloo-magic.git -c 37a53fefe0a267fe3f4704c35e3721a4b6032f2a
```

You can verify by looking at the feature list.

```
thetool list

Name:        squash
Repository:  https://github.com/axhixh/envoy-squash.git
Commit:      9397617b238cc4f17a0a3f0dc24194baf506ac97
Enabled:     true

Name:        echo
Repository:  https://github.com/axhixh/echo.git
Commit:      37a53fefe0a267fe3f4704c35e3721a4b6032f2a
Enabled:     false

Name:        lambda
Repository:  git@github.com:solo-io/glue-lambda.git
Commit:      9aeb7747286c9116d9f531f7cc2c3331a8a23c7f
Enabled:     true

Name:        magic
Repository:  https://github.com/axhixh/glue-magic.git
Commit:      37a53fefe0a267fe3f4704c35e3721a4b6032f2a
Enabled:     true

```

Adding your feature repository to the list of features will check out the repository to `repositories` 
directory in your working directory (e.g. `gloo`).

To learn more about writing your own Gloo feature, please read the Gloo documentation.