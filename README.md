

<h1 align="center">
    <img src="theTool-01.png" alt="thetool" width="200" height="233">
  <br>
  </h1>


<h3 align="center">Building, deploying and adding <br> features to Gloo</h3>
<BR>
    
Automates custom builds and deployments of gloo and Envoy, simplifying the development process and the addition of user-contributed plugins. TheTool also supports building a “lean” Gloo, that only contains desired features without bloat.

## Installing
### Prerequisite
`thetool` uses Git to checkout code. Please visit [download Git](https://git-scm.com/downloads)

 > `thetool` uses [bash](https://www.gnu.org/software/bash/manual/html_node/What-is-Bash_003f.html) scripts. Please make sure it is on your PATH. On Windows, git installs a version of bash. For example, `bash.exe` is avialable in your Git install directory `C:\Program Files\Git\bin` 

`thetool` uses Docker to run the build process and build Docker images. To install Docker, please
visit [Install Docker](https://docs.docker.com/install/)

You need [Helm](https://helm.sh/) to deploy gloo to Kubernetes. For information, please visit
[Install Helm](https://docs.helm.sh/using_helm/#installing-helm)

### Downloading and Installing
Download the latest release from https://github.com/solo-io/thetool/releases/latest/

If you prefer to compile your own binary or work on the development of `thetool` please use the following command:

```
go get github.com/solo-io/thetool
``` 

## Getting Started
### Initialize
Create a working directory for `thetool`

```
mkdir gloo
cd gloo
```

Initialize `thetool` with default set of gloo features. Optionally, you specify a default user id for Docker
using the `-u` flag.

```
thetool init -u solo-io
```

You can look at the default set of features by using the `list` command.

```
thetool list

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             aws_lambda
Gloo Directory:   aws
Envoy Directory:  aws/envoy
Enabled:          true

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             google_functions
Gloo Directory:   google
Envoy Directory:  google/envoy
Enabled:          true

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             kubernetes
Gloo Directory:   kubernetes
Envoy Directory:  
Enabled:          true

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             transformation
Gloo Directory:   transformation
Envoy Directory:  transformation/envoy
Enabled:          true


```

### Select the gloo features
You can enable or disable any of the features calling `enable` or `disable` command with the name of the feature.

```
thetool disable -n aws_lambda
```

If you list again, you will see the `echo` feature is disabled.

```
thetool list

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             aws_lambda
Gloo Directory:   aws
Envoy Directory:  aws/envoy
Enabled:          false

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             google_functions
Gloo Directory:   google
Envoy Directory:  google/envoy
Enabled:          true

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             kubernetes
Gloo Directory:   kubernetes
Envoy Directory:  
Enabled:          true

Repository:       https://github.com/solo-io/gloo-plugins.git
Name:             transformation
Gloo Directory:   transformation
Envoy Directory:  transformation/envoy
Enabled:          true

```

### Build
Once you have selected the features you want to include, you can build gloo and its components using the `build` command.

```
thetool build all
```

You can also choose to build individual components of gloo by specifying the name of the component like `envoy` or `gloo`.
To get a complete list of available components please run `thetool build --help`

The build command builds the appropriate binaries and their corresponding Docker images. It then publishes these images to Docker registry. If you do not want to publish, you can pass a flag to `thetool`

```
thetool build all --publish=false
```

Note: In order to deploy gloo to Kubernetes, you need to publish the Docker images.

> When building Envoy, [Bazel](https://bazel.build) build can fail with the error message: `gcc: internal compiler error: Killed (program cc1plus)`, if the virtual machine is out of memory. You can fix it by either reducing the number of cores or increasing the RAM on Docker VM. You can set the VM to 2GB RAM and 2 CPUs for a working configuration.

### Deploy

You can use the `deploy` command to deploy gloo and its components to different environments.

Here, we are looking at deploying gloo to Kubernetes. `thetool` uses Helm to deploy gloo
and its components.

Note: If you used custom Docker tags when building gloo and its components, you must provide
the same tag to `deploy` command to deploy those images.

```
thetool deploy k8s
```

If you want to generate the Helm chart values without deploying please pass the `--dry-run` flag.

The Helm chart used by gloo is available at [gloo-chart](https://github.com/solo-io/gloo-chart)

## Adding Your Own Feature

`thetool` can build gloo with your custom gloo features by adding your own feature repository to the list.

You can add or remove your feature repository using `add` and `delete` commands.

```
thetool add -r https://github.com/axhixh/gloo-magic.git -c 37a53fefe0a267fe3f4704c35e3721a4b6032f2a
```

You can verify by looking at the repository list with `list-repo` command.

```
thetool list-repo

Repository:  https://github.com/solo-io/gloo-plugins.git
Commit:      7bff2ff6c6ee707d8c09100de0bb7f869bd7488d

Repository:  https://github.com/axhixh/gloo-magic.git
Commit:      7bff2ff6c6ee707d8c09100de0bb7f869bd7488d
```

When you add a gloo feature repository, it loads the file `features.json` in the root folder to
find what features are available. It uses the file to identify the gloo plugin folder and envoy
filter folder for the feature.

### Updating a Feature Repository
You can get a list of feature repositories currently being used by `thetool` using the command:

```
thetool list-repo

Repository:  https://github.com/solo-io/gloo-plugins.git
Commit:      1f64f096161269a45aaaa533cab6de786a867287
```

You can update the version of repository being used to a newer version using the command:

```
thetool update -r https://github.com/solo-io/gloo-plugins.git -c 282a844ea3ed2527f5044408c9c98bc7ee027cd2

Updated repository https://github.com/solo-io/gloo-plugins.git to commit hash 282a844ea3ed2527f5044408c9c98bc7ee027cd2
```

For more information, read [Gloo documentation](https://gloo.solo.io/thetool/quickstart/)
To learn more about writing your own gloo feature, please read [Building Custom Gloo](https://gloo.solo.io/thetool/custom/)
