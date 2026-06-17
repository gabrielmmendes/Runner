package cmd

import (
	"fmt"
	"os"

	"github.com/gabrielmmendes/runner/internal/logging"
	"github.com/gabrielmmendes/runner/internal/version"
	"github.com/spf13/cobra"
)

var (
	logFormat string
	logLevel  string
	verbose   bool
	quiet     bool
)

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI multiplataforma para assinatura digital FHIR via assinador.jar",
	Long: `assinatura orquestra o serviço assinador-java para criar e validar
assinaturas digitais FHIR (caso de uso SES-GO/UFG).

Provisiona o JDK automaticamente quando ausente, sobe o assinador.jar sob
demanda (modo local) ou o gerencia como servidor HTTP de longa duração.`,
	Example: `  # Assinar um Bundle FHIR (auto-start do assinador.jar)
  assinatura sign --bundle bundle.json --provenance prov.json \
    --cert-chain chain.pem --timestamp 1751328000 --strategy iat \
    --policy-id "https://ex/policy|1.0.0" --config config.json \
    --crypto-type token --pkcs11-id chave1

  # Modo servidor de longa duração
  assinatura server start --port 8085
  assinatura server status
  assinatura server stop

  # Diagnóstico
  assinatura --version
  assinatura sign --verbose ...    # logs debug em stderr`,
	Version:           version.String(),
	PersistentPreRunE: initLogging,
	SilenceUsage:      true,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&logFormat, "log-format", "text", "formato dos logs: text|json")
	pf.StringVar(&logLevel, "log-level", "info", "nível de log: debug|info|warn|error")
	pf.BoolVarP(&verbose, "verbose", "v", false, "logs detalhados (equivale a --log-level debug)")
	pf.BoolVarP(&quiet, "quiet", "q", false, "apenas erros (equivale a --log-level error)")

	rootCmd.SetVersionTemplate("assinatura {{.Version}}\n")
}

func initLogging(cmd *cobra.Command, _ []string) error {
	if verbose && quiet {
		return fmt.Errorf("--verbose e --quiet são mutuamente exclusivos")
	}
	// Precedência: --quiet > --verbose > --log-level.
	switch {
	case quiet:
		logLevel = "error"
	case verbose:
		logLevel = "debug"
	}

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
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}
