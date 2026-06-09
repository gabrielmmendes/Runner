//go:build integration

package sign

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestIntegration_Metrics verifica que o endpoint Prometheus /metrics
// responde no formato esperado e contabiliza a requisição de sign disparada.
func TestIntegration_Metrics(t *testing.T) {
	baseURL, stop := startServer(t)
	defer stop()

	// dispara uma assinatura para gerar métricas
	req := &SignRequest{
		Bundle:    json.RawMessage(`{"resourceType":"Bundle","id":"metrics"}`),
		Strategy:  "iat",
		PolicyId:  "https://policy.example/metrics|0.1.0",
		Timestamp: 1751500000,
		CryptoMaterial: CryptoMaterial{
			Type:       "smartcard",
			Pin:        "1234",
			Identifier: "default",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := Post(ctx, baseURL, req, 30*time.Second); err != nil {
		t.Fatalf("Post: %v", err)
	}

	resp, err := http.Get(baseURL + "/metrics")
	if err != nil {
		t.Fatalf("GET /metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("/metrics HTTP %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "version=0.0.4") {
		t.Fatalf("Content-Type inesperado: %q", ct)
	}

	body, _ := io.ReadAll(resp.Body)
	text := string(body)
	for _, want := range []string{
		"assinador_uptime_seconds",
		"assinador_requests_total",
		"assinador_request_duration_seconds_bucket",
		`path="/api/sign"`,
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("/metrics não contém %q\n%s", want, text)
		}
	}
}
