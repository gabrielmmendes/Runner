package sign

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

func LoadData(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("leitura de %s falhou: %w", path, err)
	}
	return string(raw), nil
}

func LoadResource(path, expectedType string) (json.RawMessage, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leitura de %s falhou: %w", path, err)
	}
	var head struct {
		ResourceType string `json:"resourceType"`
	}
	if err := json.Unmarshal(raw, &head); err != nil {
		return nil, fmt.Errorf("%s: JSON inválido: %w", path, err)
	}
	if head.ResourceType != expectedType {
		return nil, fmt.Errorf("%s: resourceType=%q (esperado %q)", path, head.ResourceType, expectedType)
	}
	return json.RawMessage(raw), nil
}

func LoadBundle(path string) (json.RawMessage, error) {
	return LoadResource(path, "Bundle")
}

func LoadProvenance(path string) (json.RawMessage, error) {
	raw, err := LoadResource(path, "Provenance")
	if err != nil {
		return nil, err
	}
	var p struct {
		Target []json.RawMessage `json:"target"`
	}
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("%s: parse de target[] falhou: %w", path, err)
	}
	if len(p.Target) == 0 {
		return nil, fmt.Errorf("%s: Provenance.target[] vazio", path)
	}
	return raw, nil
}

func LoadCertChain(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leitura de %s falhou: %w", path, err)
	}
	if certs, err := tryJsonArray(raw); err == nil {
		if len(certs) < 2 {
			return nil, errors.New("cert chain: mínimo 2 certificados (folha + raiz)")
		}
		return certs, nil
	}
	certs, err := tryPemBundle(raw)
	if err != nil {
		return nil, fmt.Errorf("cert chain: formato inválido (esperado PEM ou JSON array): %w", err)
	}
	if len(certs) < 2 {
		return nil, errors.New("cert chain: mínimo 2 certificados (folha + raiz)")
	}
	return certs, nil
}

func tryJsonArray(raw []byte) ([]string, error) {
	trimmed := strings.TrimSpace(string(raw))
	if !strings.HasPrefix(trimmed, "[") {
		return nil, errors.New("não é JSON array")
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err != nil {
		return nil, err
	}
	for i, c := range arr {
		if _, err := base64.StdEncoding.DecodeString(c); err != nil {
			return nil, fmt.Errorf("cert[%d]: base64 inválido", i)
		}
	}
	return arr, nil
}

func tryPemBundle(raw []byte) ([]string, error) {
	var out []string
	rest := raw
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("bloco PEM inesperado: %s", block.Type)
		}
		out = append(out, base64.StdEncoding.EncodeToString(block.Bytes))
	}
	if len(out) == 0 {
		return nil, errors.New("nenhum bloco CERTIFICATE encontrado")
	}
	return out, nil
}

func LoadConfig(path string) (json.RawMessage, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("leitura de %s falhou: %w", path, err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		return nil, fmt.Errorf("%s: JSON inválido: %w", path, err)
	}
	known := map[string]bool{
		"verification": true, "trustStore": true, "temporalPolicy": true,
		"security": true, "middlewareCrypto": true,
	}
	for k := range top {
		if !known[k] {
			return nil, fmt.Errorf("%s: chave desconhecida em operationalConfig: %q", path, k)
		}
	}
	return json.RawMessage(raw), nil
}
