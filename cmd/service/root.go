package service

import "github.com/spf13/cobra"

func ServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "manage supporting services for Gloo",
	}
	cmd.AddCommand(enableCmd(), disableCmd(), listCmd(),
		markInstallCmd(), markConfigOnlyCmd())
	return cmd
}
