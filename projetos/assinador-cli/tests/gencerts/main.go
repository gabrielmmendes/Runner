package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	now := time.Now()

	// Root CA
	rootKey := must(ecdsa.GenerateKey(elliptic.P256(), rand.Reader))
	rootTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test Root CA", Organization: []string{"Teste"}},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	rootDER := must(x509.CreateCertificate(rand.Reader, rootTmpl, rootTmpl, &rootKey.PublicKey, rootKey))
	rootCert := must(x509.ParseCertificate(rootDER))

	// Leaf cert
	leafKey := must(ecdsa.GenerateKey(elliptic.P256(), rand.Reader))
	leafTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "Test Leaf Cert", Organization: []string{"Teste"}},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection},
	}
	leafDER := must(x509.CreateCertificate(rand.Reader, leafTmpl, rootCert, &leafKey.PublicKey, rootKey))

	out, err := os.Create("tests/fixtures/cert-chain.pem")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// leaf first, then root (standard chain order)
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: leafDER})
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: rootDER})
}
