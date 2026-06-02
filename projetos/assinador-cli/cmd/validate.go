package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gabrielmmendes/runner/internal/java"
	"github.com/gabrielmmendes/runner/internal/sign"
	"github.com/spf13/cobra"
)

var validateOpts sign.ValidateOptions
var validateJar string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida assinatura digital FHIR",
	Long:  "Envia dados + assinatura ao serviço assinador-java para verificação da assinatura.",
	RunE:  runValidate,
}

func init() {
	f := validateCmd.Flags()
	f.StringVar(&validateOpts.DataPath, "data", "", "caminho para arquivo de dados assinados (obrigatório)")
	f.StringVar(&validateOpts.SignaturePath, "signature", "", "caminho para arquivo com assinatura base64 (obrigatório)")
	f.StringVar(&validateOpts.ServiceURL, "service-url", "", "URL do assinador-java (default http://localhost:8080 ou env ASSINATURA_SERVICE_URL)")
	f.StringVar(&validateOpts.OutputPath, "output", "", "arquivo p/ resposta JSON (default stdout)")
	f.IntVar(&validateOpts.TimeoutSeconds, "timeout", 30, "timeout HTTP em segundos")
	f.StringVar(&validateJar, "jar", "", "caminho para assinador.jar para auto-start (ou env ASSINATURA_JAR)")

	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	if validateOpts.ServiceURL == "" {
		if env := os.Getenv("ASSINATURA_SERVICE_URL"); env != "" {
			validateOpts.ServiceURL = env
		} else {
			validateOpts.ServiceURL = "http://localhost:8080"
		}
	}

	if !java.IsRunning(validateOpts.ServiceURL) {
		javaPath, jErr := java.EnsureJava(os.Stderr)
		if jErr != nil {
			fmt.Fprintf(os.Stderr, "aviso: auto-start indisponível — %v\n", jErr)
		} else if jarPath, jarErr := java.FindJar(validateJar); jarErr != nil {
			fmt.Fprintf(os.Stderr, "aviso: auto-start indisponível — %v\n", jarErr)
		} else {
			fmt.Fprintf(os.Stderr, "assinador-java offline — iniciando %s ...\n", jarPath)
			pid, err := java.Start(java.StartOptions{
				JavaPath: javaPath,
				JarPath:  jarPath,
				Port:     java.DefaultPort,
			})
			if err != nil {
				return fmt.Errorf("auto-start falhou: %w", err)
			}
			fmt.Fprintf(os.Stderr, "aguardando assinador-java (PID=%d)...\n", pid)
			if !java.WaitReady(validateOpts.ServiceURL, 90*time.Second) {
				return fmt.Errorf("assinador-java não respondeu em 90s — verifique logs")
			}
			fmt.Fprintln(os.Stderr, "assinador-java UP")
		}
	}

	if err := sign.ValidateOpts(&validateOpts); err != nil {
		return err
	}

	data, err := sign.LoadData(validateOpts.DataPath)
	if err != nil {
		return err
	}
	signature, err := sign.LoadData(validateOpts.SignaturePath)
	if err != nil {
		return err
	}

	req := &sign.ValidateRequest{
		Data:      data,
		Signature: signature,
	}

	ctx := context.Background()
	resp, err := sign.PostValidate(ctx, validateOpts.ServiceURL, req, time.Duration(validateOpts.TimeoutSeconds)*time.Second)
	if err != nil {
		return err
	}

	if validateOpts.OutputPath != "" {
		if err := os.WriteFile(validateOpts.OutputPath, resp, 0o600); err != nil {
			return fmt.Errorf("escrita em %s falhou: %w", validateOpts.OutputPath, err)
		}
		return nil
	}
	fmt.Println(string(resp))
	return nil
}
