// Package verify checks artifact integrity (SHA-256) and authenticity
// (Cosign keyless signatures) before the CLI executes a downloaded jar.
package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// CertIdentityRegexp / CertOIDCIssuer pin the expected keyless signer.
// They match the GitHub Actions OIDC identity used in release.yml.
const (
	CertIdentityRegexp = `^https://github.com/gabrielmmendes/runner/`
	CertOIDCIssuer     = "https://token.actions.githubusercontent.com"
)

// SHA256File returns the lowercase hex SHA-256 of the file at path.
func SHA256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// CheckSHA256 fails if the file digest does not match expected (case-insensitive).
func CheckSHA256(path, expected string) error {
	got, err := SHA256File(path)
	if err != nil {
		return err
	}
	if !strings.EqualFold(got, strings.TrimSpace(expected)) {
		return fmt.Errorf("SHA-256 inválido para %s: esperado %s, obtido %s", path, expected, got)
	}
	return nil
}

// CosignAvailable reports whether the cosign binary is on PATH.
func CosignAvailable() bool {
	_, err := exec.LookPath("cosign")
	return err == nil
}

// Blob verifies a blob against its Cosign keyless signature (.sig) and
// certificate (.pem) using transparency-log verification.
func Blob(blobPath, sigPath, certPath string) error {
	if !CosignAvailable() {
		return fmt.Errorf("cosign não encontrado no PATH")
	}
	cmd := exec.Command("cosign", "verify-blob",
		"--certificate", certPath,
		"--signature", sigPath,
		"--certificate-identity-regexp", CertIdentityRegexp,
		"--certificate-oidc-issuer", CertOIDCIssuer,
		blobPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cosign verify-blob falhou: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// JarIfSignatures verifies jarPath when sidecar .sig/.pem files exist next to it.
//   - skip=true  → skipped with a warning (returned as nil, msg via warn).
//   - no sidecars → skipped silently (nil).
//   - sidecars present but cosign missing → error (fail closed).
//
// The bool return reports whether verification actually ran.
func JarIfSignatures(jarPath string, skip bool, warn func(string)) (bool, error) {
	sig := jarPath + ".sig"
	cert := jarPath + ".pem"
	if !fileExists(sig) || !fileExists(cert) {
		return false, nil
	}
	if skip {
		if warn != nil {
			warn(fmt.Sprintf("verificação Cosign IGNORADA (--skip-verify) para %s", jarPath))
		}
		return false, nil
	}
	if err := Blob(jarPath, sig, cert); err != nil {
		return false, err
	}
	return true, nil
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
