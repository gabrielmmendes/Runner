package cmd

import (
	"fmt"

	"github.com/gabrielmmendes/assinatura/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão do CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assinatura version:", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
