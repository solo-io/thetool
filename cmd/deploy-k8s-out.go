package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func DeployK8SOutCmd() *cobra.Command {
	var kubeConfig string

	cmd := &cobra.Command{
		Use:   "k8s-out",
		Short: "deploy out of Kubernetes cluster",
		Run: func(c *cobra.Command, args []string) {
			f := c.InheritedFlags()
			verbose, _ := f.GetBool("verbose")
			dryRun, _ := f.GetBool("dry-run")
			dockerUser, _ := f.GetString("docker-user")
			runDeployK8SOut(verbose, dryRun, dockerUser, kubeConfig)
		},
	}
	if home := homeDir(); home != "" {
		cmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "k",
			filepath.Join(home, ".kube", "config"), "(optional) path to Kubernetes config")
	} else {
		cmd.Flags().StringVarP(&kubeConfig, "kubeconfig", "k", "", "path to Kubernetes config")
		cmd.MarkFlagRequired("kubeconf")
	}
	return cmd
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// Gloo and its components are deployed outside the K8S cluster but
// we will use CRD for storage
func runDeployK8SOut(verbose, dryRun bool, dockerUser, kubeConfig string) {
	fmt.Println("verbose ", verbose, " dryRun ", dryRun, " dockerUser ", dockerUser, " kubeconfig ", kubeConfig)
	fmt.Println("Deploying outside Kubernetes cluster has not been implemented yet.")
	// save the gloo configuration - shared by other tools]
	// storage option is k8s CRD
	// run gloo in docker
	// run envoy in docker -- use --link option to link envoy to gloo
}
