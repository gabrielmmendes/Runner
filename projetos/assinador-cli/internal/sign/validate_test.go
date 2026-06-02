package sign

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func baseOpts(t *testing.T) *Options {
	dir := t.TempDir()
	return &Options{
		BundlePath:     writeFile(t, dir, "bundle.json", `{"resourceType":"Bundle"}`),
		ProvenancePath: writeFile(t, dir, "prov.json", `{"resourceType":"Provenance"}`),
		CertChainPath:  writeFile(t, dir, "chain.pem", "x"),
		ConfigPath:     writeFile(t, dir, "cfg.json", `{}`),
		Timestamp:      1751500000,
		Strategy:       "iat",
		PolicyId:       "https://x.example/policy|0.1.2",
		CryptoType:     "smartcard",
		Pin:            "1234",
		Identifier:     "key1",
	}
}

func TestValidate_OK(t *testing.T) {
	if err := Validate(baseOpts(t)); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	cases := []struct {
		name  string
		mut   func(*Options)
		match string
	}{
		{"bundle", func(o *Options) { o.BundlePath = "" }, "--bundle"},
		{"provenance", func(o *Options) { o.ProvenancePath = "" }, "--provenance"},
		{"cert-chain", func(o *Options) { o.CertChainPath = "" }, "--cert-chain"},
		{"config", func(o *Options) { o.ConfigPath = "" }, "--config"},
		{"strategy", func(o *Options) { o.Strategy = "" }, "--strategy"},
		{"policy-id", func(o *Options) { o.PolicyId = "" }, "--policy-id"},
		{"crypto-type", func(o *Options) { o.CryptoType = "" }, "--crypto-type"},
		{"pin", func(o *Options) { o.Pin = "" }, "PIN"},
		{"identifier", func(o *Options) { o.Identifier = "" }, "--pkcs11-id"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			o := baseOpts(t)
			c.mut(o)
			err := Validate(o)
			if err == nil || !strings.Contains(err.Error(), c.match) {
				t.Fatalf("expected error matching %q, got %v", c.match, err)
			}
		})
	}
}

func TestValidate_TimestampRange(t *testing.T) {
	for _, ts := range []int64{0, TimestampMin - 1, TimestampMax + 1} {
		o := baseOpts(t)
		o.Timestamp = ts
		if err := Validate(o); err == nil || !strings.Contains(err.Error(), "--timestamp") {
			t.Fatalf("ts=%d: expected timestamp error, got %v", ts, err)
		}
	}
}

func TestValidate_StrategyInvalid(t *testing.T) {
	o := baseOpts(t)
	o.Strategy = "lol"
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "iat|tsa") {
		t.Fatalf("expected strategy error, got %v", err)
	}
}

func TestValidate_PolicyIdMalformed(t *testing.T) {
	for _, pid := range []string{"foo", "https://x|", "https://x|1.2", "no-version"} {
		o := baseOpts(t)
		o.PolicyId = pid
		if err := Validate(o); err == nil || !strings.Contains(err.Error(), "--policy-id") {
			t.Fatalf("pid=%q: expected policy-id error, got %v", pid, err)
		}
	}
}

func TestValidate_CryptoTypeInvalid(t *testing.T) {
	o := baseOpts(t)
	o.CryptoType = "pem"
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "smartcard|token") {
		t.Fatalf("expected crypto-type error, got %v", err)
	}
}

func TestValidate_SlotNegative(t *testing.T) {
	o := baseOpts(t)
	o.SlotIdSet = true
	o.SlotId = -1
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "--pkcs11-slot") {
		t.Fatalf("expected slot error, got %v", err)
	}
}

func TestValidate_TokenLabelTooLong(t *testing.T) {
	o := baseOpts(t)
	o.TokenLabel = strings.Repeat("a", TokenLabelMaxRunes+1)
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "token-label") {
		t.Fatalf("expected token-label error, got %v", err)
	}
}

func TestValidate_TsaRequiresTsaUrl(t *testing.T) {
	o := baseOpts(t)
	o.Strategy = "tsa"
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "tsaUrl") {
		t.Fatalf("expected tsa url error, got %v", err)
	}

	dir := t.TempDir()
	o = baseOpts(t)
	o.Strategy = "tsa"
	o.ConfigPath = writeFile(t, dir, "cfg.json", `{"verification":{"tsaUrl":"https://tsa.example"}}`)
	if err := Validate(o); err != nil {
		t.Fatalf("expected ok with tsaUrl, got %v", err)
	}
}

func TestValidate_FileMissing(t *testing.T) {
	o := baseOpts(t)
	o.BundlePath = "/nope/missing.json"
	if err := Validate(o); err == nil || !strings.Contains(err.Error(), "inacessível") {
		t.Fatalf("expected file error, got %v", err)
	}
}

// --- ValidateOpts tests ---

func baseValidateOpts(t *testing.T) *ValidateOptions {
	t.Helper()
	dir := t.TempDir()
	return &ValidateOptions{
		DataPath:       writeFile(t, dir, "data.json", `{"resourceType":"Bundle"}`),
		SignaturePath:  writeFile(t, dir, "sig.b64", "dGVzdA=="),
		TimeoutSeconds: 30,
	}
}

func TestValidateOpts_OK(t *testing.T) {
	if err := ValidateOpts(baseValidateOpts(t)); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
}

func TestValidateOpts_MissingData(t *testing.T) {
	o := baseValidateOpts(t)
	o.DataPath = ""
	if err := ValidateOpts(o); err == nil || !strings.Contains(err.Error(), "--data") {
		t.Fatalf("expected --data error, got %v", err)
	}
}

func TestValidateOpts_MissingSignature(t *testing.T) {
	o := baseValidateOpts(t)
	o.SignaturePath = ""
	if err := ValidateOpts(o); err == nil || !strings.Contains(err.Error(), "--signature") {
		t.Fatalf("expected --signature error, got %v", err)
	}
}

func TestValidateOpts_DataFileNotFound(t *testing.T) {
	o := baseValidateOpts(t)
	o.DataPath = "/nope/missing.json"
	if err := ValidateOpts(o); err == nil || !strings.Contains(err.Error(), "inacessível") {
		t.Fatalf("expected file error, got %v", err)
	}
}

func TestValidateOpts_SignatureFileNotFound(t *testing.T) {
	o := baseValidateOpts(t)
	o.SignaturePath = "/nope/missing.sig"
	if err := ValidateOpts(o); err == nil || !strings.Contains(err.Error(), "inacessível") {
		t.Fatalf("expected file error, got %v", err)
	}
}
