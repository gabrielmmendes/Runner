package cmd

import (
	"fmt"

	"github.com/gabrielmmendes/runner/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Exibe a versão do CLI (tag + SHA curto + data de build)",
	Example: "  assinatura version\n  assinatura --version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assinatura", version.String())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
