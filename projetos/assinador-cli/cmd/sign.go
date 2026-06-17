package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gabrielmmendes/runner/internal/java"
	"github.com/gabrielmmendes/runner/internal/sign"
	"github.com/spf13/cobra"
)

var signOpts sign.Options
var signJar string
var signSkipVerify bool

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria assinatura digital FHIR (caso de uso Goiás)",
	Long:  "Envia Bundle + Provenance + cadeia de certificados ao serviço assinador-java para geração da assinatura JWS.",
	Example: `  # Assinatura completa (PIN via prompt seguro)
  assinatura sign \
    --bundle bundle.json --provenance prov.json --cert-chain chain.pem \
    --timestamp 1751328000 --strategy iat \
    --policy-id "https://exemplo/policy|1.0.0" --config config.json \
    --crypto-type token --pkcs11-id chave1

  # PIN via variável de ambiente (CI) e saída para arquivo
  ASSINATURA_PKCS11_PIN=1234 assinatura sign ... --output assinatura.json

  # Apontar para servidor já em execução
  assinatura sign ... --service-url http://localhost:8085`,
	RunE: runSign,
}

func init() {
	f := signCmd.Flags()
	f.StringVar(&signOpts.BundlePath, "bundle", "", "caminho para FHIR Bundle JSON (obrigatório)")
	f.StringVar(&signOpts.ProvenancePath, "provenance", "", "caminho para FHIR Provenance JSON (obrigatório)")
	f.StringVar(&signOpts.CertChainPath, "cert-chain", "", "caminho para cadeia X.509 (PEM ou JSON array base64) (obrigatório)")
	f.Int64Var(&signOpts.Timestamp, "timestamp", 0, "Unix UTC seconds, range [1751328000, 4102444800] (obrigatório)")
	f.StringVar(&signOpts.Strategy, "strategy", "", "iat|tsa (obrigatório)")
	f.StringVar(&signOpts.PolicyId, "policy-id", "", "URI no formato <baseURI>|<major.minor.patch> (obrigatório)")
	f.StringVar(&signOpts.ConfigPath, "config", "", "caminho para JSON de operationalConfig (obrigatório)")
	f.StringVar(&signOpts.CryptoType, "crypto-type", "", "smartcard|token (obrigatório)")
	f.StringVar(&signOpts.Pin, "pkcs11-pin", "", "PIN do dispositivo (ou env ASSINATURA_PKCS11_PIN, ou prompt stdin)")
	f.StringVar(&signOpts.Identifier, "pkcs11-id", "", "identificador da chave privada (obrigatório)")
	f.IntVar(&signOpts.SlotId, "pkcs11-slot", 0, "slot ID PKCS#11 (≥0)")
	f.StringVar(&signOpts.TokenLabel, "pkcs11-token-label", "", "label do token (≤32 chars UTF-8)")
	f.StringVar(&signOpts.ServiceURL, "service-url", "", "URL do assinador-java (default http://localhost:8080 ou env ASSINATURA_SERVICE_URL)")
	f.StringVar(&signOpts.OutputPath, "output", "", "arquivo p/ resposta JSON (default stdout)")
	f.IntVar(&signOpts.TimeoutSeconds, "timeout", 30, "timeout HTTP em segundos")
	f.StringVar(&signJar, "jar", "", "caminho para assinador.jar para auto-start (ou env ASSINATURA_JAR)")
	f.BoolVar(&signSkipVerify, "skip-verify", false, "ignora verificação Cosign do jar (não recomendado)")

	rootCmd.AddCommand(signCmd)
}

func runSign(cmd *cobra.Command, args []string) error {
	signOpts.SlotIdSet = cmd.Flags().Changed("pkcs11-slot")

	if signOpts.Pin == "" {
		if env := os.Getenv("ASSINATURA_PKCS11_PIN"); env != "" {
			signOpts.Pin = env
		} else {
			pin, err := promptPin()
			if err != nil {
				return err
			}
			signOpts.Pin = pin
		}
	}

	if signOpts.ServiceURL == "" {
		if env := os.Getenv("ASSINATURA_SERVICE_URL"); env != "" {
			signOpts.ServiceURL = env
		} else {
			signOpts.ServiceURL = "http://localhost:8080"
		}
	}

	if !java.IsRunning(signOpts.ServiceURL) {
		javaPath, jErr := java.EnsureJava(os.Stderr)
		if jErr != nil {
			fmt.Fprintf(os.Stderr, "aviso: auto-start indisponível — %v\n", jErr)
		} else if jarPath, jarErr := java.FindJar(signJar); jarErr != nil {
			fmt.Fprintf(os.Stderr, "aviso: auto-start indisponível — %v\n", jarErr)
		} else if vErr := verifyJar(jarPath, signSkipVerify); vErr != nil {
			return vErr
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
			if !java.WaitReady(signOpts.ServiceURL, 90*time.Second) {
				return fmt.Errorf("assinador-java não respondeu em 90s — verifique logs")
			}
			fmt.Fprintln(os.Stderr, "assinador-java UP")
		}
	}

	if err := sign.Validate(&signOpts); err != nil {
		return err
	}

	bundle, err := sign.LoadBundle(signOpts.BundlePath)
	if err != nil {
		return err
	}
	provenance, err := sign.LoadProvenance(signOpts.ProvenancePath)
	if err != nil {
		return err
	}
	certChain, err := sign.LoadCertChain(signOpts.CertChainPath)
	if err != nil {
		return err
	}
	config, err := sign.LoadConfig(signOpts.ConfigPath)
	if err != nil {
		return err
	}

	req := &sign.SignRequest{
		Bundle:            bundle,
		Provenance:        provenance,
		CertChain:         certChain,
		Timestamp:         signOpts.Timestamp,
		Strategy:          signOpts.Strategy,
		PolicyId:          signOpts.PolicyId,
		OperationalConfig: config,
		CryptoMaterial: sign.CryptoMaterial{
			Type:       signOpts.CryptoType,
			Pin:        signOpts.Pin,
			Identifier: signOpts.Identifier,
			TokenLabel: signOpts.TokenLabel,
		},
	}
	if signOpts.SlotIdSet {
		s := signOpts.SlotId
		req.CryptoMaterial.SlotId = &s
	}

	ctx := context.Background()
	resp, err := sign.Post(ctx, signOpts.ServiceURL, req, time.Duration(signOpts.TimeoutSeconds)*time.Second)
	if err != nil {
		return err
	}

	if signOpts.OutputPath != "" {
		if err := os.WriteFile(signOpts.OutputPath, resp, 0o600); err != nil {
			return fmt.Errorf("escrita em %s falhou: %w", signOpts.OutputPath, err)
		}
		return nil
	}
	fmt.Println(string(resp))
	return nil
}

func promptPin() (string, error) {
	fmt.Fprint(os.Stderr, "PIN: ")
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("leitura de PIN falhou: %w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}
