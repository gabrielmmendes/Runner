package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSHA256File(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "blob")
	if err := os.WriteFile(p, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	// echo -n hello | sha256sum
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	got, err := SHA256File(p)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("SHA256File = %s, want %s", got, want)
	}
	if err := CheckSHA256(p, want); err != nil {
		t.Fatalf("CheckSHA256 ok esperado: %v", err)
	}
	if err := CheckSHA256(p, "deadbeef"); err == nil {
		t.Fatal("CheckSHA256 deveria falhar com digest errado")
	}
}

func TestJarIfSignatures_NoSidecars(t *testing.T) {
	dir := t.TempDir()
	jar := filepath.Join(dir, "assinador.jar")
	if err := os.WriteFile(jar, []byte("jar"), 0o600); err != nil {
		t.Fatal(err)
	}
	ran, err := JarIfSignatures(jar, false, nil)
	if err != nil {
		t.Fatalf("sem sidecars não deve falhar: %v", err)
	}
	if ran {
		t.Fatal("verificação não deveria rodar sem sidecars")
	}
}

func TestJarIfSignatures_Skip(t *testing.T) {
	dir := t.TempDir()
	jar := filepath.Join(dir, "assinador.jar")
	for _, ext := range []string{"", ".sig", ".pem"} {
		if err := os.WriteFile(jar+ext, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var warned bool
	ran, err := JarIfSignatures(jar, true, func(string) { warned = true })
	if err != nil {
		t.Fatalf("skip não deve falhar: %v", err)
	}
	if ran {
		t.Fatal("skip não deveria verificar")
	}
	if !warned {
		t.Fatal("skip deveria emitir aviso")
	}
}
