//go:build integration

package sign

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/gabrielmmendes/runner/internal/java"
)

// Para rodar:
//   go test -tags=integration ./internal/sign/...
//
// Requisitos:
//   - JDK 21 acessivel (JAVA_HOME, PATH ou ~/.hubsaude/jdk)
//   - assinador.jar disponivel (flag ASSINATURA_JAR ou Maven target/)
//
// Sem ambos, teste e pulado.

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("free port: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func startServer(t *testing.T) (baseURL string, stop func()) {
	t.Helper()

	javaPath, err := java.FindJava()
	if err != nil {
		t.Skipf("Java indisponivel: %v", err)
	}
	jarPath, err := java.FindJar(os.Getenv("ASSINATURA_JAR"))
	if err != nil {
		t.Skipf("assinador.jar indisponivel: %v", err)
	}

	port := freePort(t)
	pid, err := java.Start(java.StartOptions{
		JavaPath: javaPath,
		JarPath:  jarPath,
		Port:     port,
		LogOut:   os.Stderr,
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	baseURL = fmt.Sprintf("http://localhost:%d", port)

	if !java.WaitReady(baseURL, 90*time.Second) {
		_ = java.Stop(port)
		t.Fatalf("servidor PID=%d nao respondeu em 90s", pid)
	}

	return baseURL, func() {
		_ = java.Stop(port)
	}
}

func TestIntegration_Sign(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	req := &SignRequest{
		Bundle:    json.RawMessage(`{"resourceType":"Bundle","id":"itg"}`),
		Strategy:  "iat",
		PolicyId:  "https://policy.example/itg|0.1.0",
		Timestamp: 1751500000,
		CryptoMaterial: CryptoMaterial{
			Type:       "smartcard",
			Pin:        "1234",
			Identifier: "default",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := Post(ctx, baseURL, req, 30*time.Second)
	if err != nil {
		t.Fatalf("Post: %v", err)
	}

	var body struct {
		Success   bool   `json:"success"`
		Signature string `json:"signature"`
		Algorithm string `json:"algorithm"`
	}
	if err := json.Unmarshal(resp, &body); err != nil {
		t.Fatalf("unmarshal: %v — body=%s", err, string(resp))
	}
	if !body.Success {
		t.Fatalf("success=false: %s", string(resp))
	}
	if body.Signature == "" {
		t.Fatalf("signature vazia: %s", string(resp))
	}
	if body.Algorithm != "SHA256withRSA" {
		t.Fatalf("algorithm=%q (esperado SHA256withRSA)", body.Algorithm)
	}
}

func TestIntegration_Sign_WithFullPayload(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	req := &SignRequest{
		Bundle:     json.RawMessage(`{"resourceType":"Bundle","id":"full","type":"document","entry":[]}`),
		Provenance: json.RawMessage(`{"resourceType":"Provenance","target":[{"reference":"Bundle/full"}]}`),
		CertChain:  []string{"dGVzdA==", "dGVzdDI="},
		Strategy:   "iat",
		PolicyId:   "https://policy.example/full|1.0.0",
		Timestamp:  1751500000,
		OperationalConfig: json.RawMessage(`{
			"verification":{"tsaUrl":"","checkRevocation":false},
			"trustStore":{"type":"JKS","path":"","password":""},
			"temporalPolicy":{"allowedClockSkewSeconds":300},
			"security":{"requireSecureChannel":false},
			"middlewareCrypto":{"library":"","slotDescription":""}
		}`),
		CryptoMaterial: CryptoMaterial{
			Type:       "token",
			Pin:        "5678",
			Identifier: "full-key",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := Post(ctx, baseURL, req, 30*time.Second)
	if err != nil {
		t.Fatalf("Post: %v", err)
	}

	var body struct {
		Success   bool   `json:"success"`
		Signature string `json:"signature"`
		Algorithm string `json:"algorithm"`
	}
	if err := json.Unmarshal(resp, &body); err != nil {
		t.Fatalf("unmarshal: %v — body=%s", err, string(resp))
	}
	if !body.Success {
		t.Fatalf("success=false: %s", string(resp))
	}
	if body.Signature == "" {
		t.Fatal("signature vazia")
	}
}

func TestIntegration_Sign_InvalidPayload(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	req := &SignRequest{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := Post(ctx, baseURL, req, 30*time.Second)
	if err == nil {
		t.Fatal("expected error for empty request")
	}
}

func TestIntegration_Validate(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	req := &ValidateRequest{
		Data:      `{"resourceType":"Bundle","id":"val-test"}`,
		Signature: "dGVzdFNpZ25hdHVyZQ==",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := PostValidate(ctx, baseURL, req, 30*time.Second)
	if err != nil {
		t.Fatalf("PostValidate: %v", err)
	}

	var body struct {
		Valid   bool   `json:"valid"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(resp, &body); err != nil {
		t.Fatalf("unmarshal: %v — body=%s", err, string(resp))
	}
	if !body.Valid {
		t.Fatalf("valid=false: %s", string(resp))
	}
	if body.Message == "" {
		t.Fatal("message vazia")
	}
}

func TestIntegration_Validate_EmptyFields(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	req := &ValidateRequest{
		Data:      "",
		Signature: "",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := PostValidate(ctx, baseURL, req, 30*time.Second)
	if err != nil {
		t.Fatalf("PostValidate: %v", err)
	}

	var body struct {
		Valid bool `json:"valid"`
	}
	if err := json.Unmarshal(resp, &body); err != nil {
		t.Fatalf("unmarshal: %v — body=%s", err, string(resp))
	}
}

func TestIntegration_StopByPort(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	if !java.IsRunning(baseURL) {
		t.Fatal("servidor deveria estar UP antes do stop")
	}

	// stop() ja invoca java.Stop(port) — apenas verifica que de fato cai
	stop()

	// pequena espera pro processo encerrar
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if !java.IsRunning(baseURL) {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatal("servidor continuou UP apos Stop")
}
