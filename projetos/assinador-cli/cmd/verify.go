package cmd

import (
	"fmt"
	"os"

	"github.com/gabrielmmendes/runner/internal/logging"
	"github.com/gabrielmmendes/runner/internal/verify"
)

// verifyJar runs Cosign integrity verification on jarPath when sidecar
// signature files exist. Fails closed if verification cannot run.
func verifyJar(jarPath string, skip bool) error {
	ran, err := verify.JarIfSignatures(jarPath, skip, func(msg string) {
		logging.L().Warn(msg)
		fmt.Fprintf(os.Stderr, "aviso: %s\n", msg)
	})
	if err != nil {
		return fmt.Errorf("verificação de integridade do jar falhou: %w", err)
	}
	if ran {
		logging.L().Info("jar verificado via Cosign", "jar", jarPath)
		fmt.Fprintf(os.Stderr, "integridade do jar verificada (Cosign): %s\n", jarPath)
	}
	return nil
}
