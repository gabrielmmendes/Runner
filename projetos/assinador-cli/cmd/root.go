package cmd

import (
	"fmt"
	"os"

	"github.com/gabrielmmendes/runner/internal/logging"
	"github.com/spf13/cobra"
)

var (
	logFormat string
	logLevel  string
)

var rootCmd = &cobra.Command{
	Use:               "assinatura",
	Short:             "CLI para simulação de assinatura digital",
	Long:              "Ferramenta CLI para integração com assinador.jar",
	PersistentPreRunE: initLogging,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&logFormat, "log-format", "text", "formato dos logs: text|json")
	pf.StringVar(&logLevel, "log-level", "info", "nível de log: debug|info|warn|error")
}

func initLogging(cmd *cobra.Command, _ []string) error {
	format, err := logging.ParseFormat(logFormat)
	if err != nil {
		return err
	}
	level, err := logging.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logging.Init(format, level)
	logging.WithCommand(cmd.Name()).Debug("comando iniciado")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
