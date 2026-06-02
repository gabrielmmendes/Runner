package sign

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"unicode/utf8"
)

const (
	TimestampMin       int64 = 1751328000
	TimestampMax       int64 = 4102444800
	TokenLabelMaxRunes       = 32
)

var policyIdRegex = regexp.MustCompile(`^https?://[^\s|]+\|\d+\.\d+\.\d+$`)

func Validate(opts *Options) error {
	if opts.BundlePath == "" {
		return errors.New("--bundle é obrigatório")
	}
	if opts.ProvenancePath == "" {
		return errors.New("--provenance é obrigatório")
	}
	if opts.CertChainPath == "" {
		return errors.New("--cert-chain é obrigatório")
	}
	if opts.ConfigPath == "" {
		return errors.New("--config é obrigatório")
	}
	for flag, path := range map[string]string{
		"--bundle":     opts.BundlePath,
		"--provenance": opts.ProvenancePath,
		"--cert-chain": opts.CertChainPath,
		"--config":     opts.ConfigPath,
	} {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("%s: arquivo inacessível: %s", flag, path)
		}
	}

	if opts.Timestamp < TimestampMin || opts.Timestamp > TimestampMax {
		return fmt.Errorf("--timestamp fora do range [%d, %d]", TimestampMin, TimestampMax)
	}

	switch opts.Strategy {
	case "iat", "tsa":
	case "":
		return errors.New("--strategy é obrigatório (iat|tsa)")
	default:
		return fmt.Errorf("--strategy inválida: %q (esperado: iat|tsa)", opts.Strategy)
	}

	if opts.PolicyId == "" {
		return errors.New("--policy-id é obrigatório")
	}
	if !policyIdRegex.MatchString(opts.PolicyId) {
		return errors.New("--policy-id mal-formado (esperado: <baseURI>|<major.minor.patch>)")
	}

	switch opts.CryptoType {
	case "smartcard", "token":
	case "":
		return errors.New("--crypto-type é obrigatório (smartcard|token)")
	default:
		return fmt.Errorf("--crypto-type inválido: %q (esperado: smartcard|token)", opts.CryptoType)
	}

	if opts.Pin == "" {
		return errors.New("PIN é obrigatório (--pkcs11-pin, env ASSINATURA_PKCS11_PIN ou stdin)")
	}
	if opts.Identifier == "" {
		return errors.New("--pkcs11-id é obrigatório")
	}
	if opts.SlotIdSet && opts.SlotId < 0 {
		return errors.New("--pkcs11-slot deve ser ≥ 0")
	}
	if opts.TokenLabel != "" && utf8.RuneCountInString(opts.TokenLabel) > TokenLabelMaxRunes {
		return fmt.Errorf("--pkcs11-token-label excede %d chars UTF-8", TokenLabelMaxRunes)
	}

	if opts.Strategy == "tsa" {
		if err := validateTsaUrl(opts.ConfigPath); err != nil {
			return err
		}
	}

	return nil
}

func ValidateOpts(opts *ValidateOptions) error {
	if opts.DataPath == "" {
		return errors.New("--data é obrigatório")
	}
	if _, err := os.Stat(opts.DataPath); err != nil {
		return fmt.Errorf("--data: arquivo inacessível: %s", opts.DataPath)
	}
	if opts.SignaturePath == "" {
		return errors.New("--signature é obrigatório")
	}
	if _, err := os.Stat(opts.SignaturePath); err != nil {
		return fmt.Errorf("--signature: arquivo inacessível: %s", opts.SignaturePath)
	}
	return nil
}

func validateTsaUrl(configPath string) error {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("--config: leitura falhou: %w", err)
	}
	var cfg struct {
		Verification struct {
			TsaUrl string `json:"tsaUrl"`
		} `json:"verification"`
	}
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return fmt.Errorf("--config: JSON inválido: %w", err)
	}
	if cfg.Verification.TsaUrl == "" {
		return errors.New("--strategy=tsa exige config.verification.tsaUrl")
	}
	return nil
}
