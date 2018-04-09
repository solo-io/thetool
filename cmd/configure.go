package cmd

import (
	"fmt"

	"github.com/solo-io/thetool/pkg/config"
	"github.com/spf13/cobra"
)

func ConfigureCmd() *cobra.Command {
	conf := config.Config{}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "configure the tool",
		Run: func(c *cobra.Command, args []string) {
			runConfigure(&conf)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&conf.EnvoyHash, "envoy-hash", "e", "", "Envoy commit hash to use")
	flags.StringVar(&conf.EnvoyRepoUser, "envoy-repo-user", "", "Envoy repository user")
	flags.StringVar(&conf.EnvoyCommonHash, "envoy-common-hash", "", "Hash for Soloio Envoy Common")
	flags.StringVarP(&conf.GlooHash, "gloo-hash", "g", "", "Gloo commit hash to use")
	flags.StringVar(&conf.GlooRepo, "gloo-repo", "", "Gloo git repository")
	flags.StringVarP(&conf.DockerUser, "user", "u", "", "default Docker user")
	flags.StringVar(&conf.EnvoyBuilderHash, "envoy-builder-hash", "", "hash for envoy build container")

	return cmd
}

func runConfigure(c *config.Config) {
	existing, err := config.Load(config.ConfigFile)
	if err != nil {
		fmt.Println("Unable to read current configuration:", err)
		return
	}

	if c.DockerUser != "" {
		existing.DockerUser = c.DockerUser
	}
	if c.EnvoyBuilderHash != "" {
		existing.EnvoyBuilderHash = c.EnvoyBuilderHash
	}
	if c.EnvoyHash != "" {
		existing.EnvoyHash = c.EnvoyHash
	}
	if c.EnvoyRepoUser != "" {
		existing.EnvoyRepoUser = c.EnvoyRepoUser
	}
	if c.EnvoyCommonHash != "" {
		existing.EnvoyCommonHash = c.EnvoyCommonHash
	}
	if c.GlooHash != "" {
		existing.GlooHash = c.GlooHash
	}
	if c.GlooRepo != "" {
		existing.GlooRepo = c.GlooRepo
	}

	if err := existing.Save(config.ConfigFile); err != nil {
		fmt.Printf("unable to save the configuration to %s: %q\n", config.ConfigFile, err)
		return
	}

	show(existing)
}

func show(c *config.Config) {
	fmt.Printf("%-20s: %s\n", "Docker User", c.DockerUser)
	fmt.Printf("%-20s: %s\n", "Envoy Builder Hash", c.EnvoyBuilderHash)
	fmt.Printf("%-20s: %s\n", "Envoy Repo User", c.EnvoyRepoUser)
	fmt.Printf("%-20s: %s\n", "Envoy Hash", c.EnvoyHash)
	fmt.Printf("%-20s: %s\n", "Envoy Common Hash", c.EnvoyCommonHash)
	fmt.Printf("%-20s: %s\n", "Gloo Hash", c.GlooHash)
	fmt.Printf("%-20s: %s\n", "Gloo Repo", c.GlooRepo)
}
