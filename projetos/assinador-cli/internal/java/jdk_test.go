package java

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFindJava_ViaJAVA_HOME(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}
	javaPath := filepath.Join(binDir, bin)
	if err := os.WriteFile(javaPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	orig := os.Getenv("JAVA_HOME")
	t.Setenv("JAVA_HOME", dir)
	defer os.Setenv("JAVA_HOME", orig)

	p, err := FindJava()
	if err != nil {
		t.Fatalf("expected java found via JAVA_HOME, got %v", err)
	}
	if p != javaPath {
		t.Fatalf("expected %s, got %s", javaPath, p)
	}
}

func TestFindJava_JAVA_HOME_Invalid(t *testing.T) {
	t.Setenv("JAVA_HOME", "/nonexistent/jdk")
	t.Setenv("PATH", "")

	_, err := FindJava()
	if err == nil {
		t.Fatal("expected error when JAVA_HOME is invalid and PATH is empty")
	}
}

func TestFindJava_NoJAVA_HOME_FallbackPATH(t *testing.T) {
	t.Setenv("JAVA_HOME", "")

	_, err := FindJava()
	// Can't guarantee java is on PATH in test env — just verify no panic
	// If java is on PATH, err == nil. If not, err != nil. Both valid.
	_ = err
}

func TestHubJavaBin_NoneProvisioned(t *testing.T) {
	dir := t.TempDir()
	jdkDir := filepath.Join(dir, "jdk")
	if err := os.MkdirAll(jdkDir, 0o755); err != nil {
		t.Fatal(err)
	}

	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
	defer func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	}()

	// Create .hubsaude/jdk with no java binary inside
	hubJdk := filepath.Join(dir, hubDirName, "jdk")
	if err := os.MkdirAll(hubJdk, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := hubJavaBin()
	if err == nil {
		t.Fatal("expected error when no JDK provisioned under .hubsaude/jdk")
	}
}

func TestHubJavaBin_Provisioned(t *testing.T) {
	dir := t.TempDir()

	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}

	hubJdk := filepath.Join(dir, hubDirName, "jdk", "jre-21", "bin")
	if err := os.MkdirAll(hubJdk, 0o755); err != nil {
		t.Fatal(err)
	}
	javaPath := filepath.Join(hubJdk, bin)
	if err := os.WriteFile(javaPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	origHome := os.Getenv("HOME")
	origUserProfile := os.Getenv("USERPROFILE")
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
	defer func() {
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserProfile)
	}()

	p, err := hubJavaBin()
	if err != nil {
		t.Fatalf("expected java found in .hubsaude/jdk, got %v", err)
	}
	if p != javaPath {
		t.Fatalf("expected %s, got %s", javaPath, p)
	}
}

func TestHubJavaBin_NoDirExists(t *testing.T) {
	dir := t.TempDir()

	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	// .hubsaude/jdk doesn't exist at all
	_, err := hubJavaBin()
	if err == nil {
		t.Fatal("expected error when .hubsaude/jdk dir missing")
	}
}

func TestEnsureJava_FindsSystemJava(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}
	javaPath := filepath.Join(binDir, bin)
	if err := os.WriteFile(javaPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("JAVA_HOME", dir)

	p, err := EnsureJava(os.Stderr)
	if err != nil {
		t.Fatalf("expected java via JAVA_HOME, got %v", err)
	}
	if p != javaPath {
		t.Fatalf("expected %s, got %s", javaPath, p)
	}
}

func TestAdoptiumOS(t *testing.T) {
	got := adoptiumOS()
	valid := map[string]bool{"mac": true, "windows": true, "linux": true}
	if !valid[got] {
		t.Fatalf("unexpected adoptiumOS: %q", got)
	}
}

func TestAdoptiumArch(t *testing.T) {
	got := adoptiumArch()
	valid := map[string]bool{"x64": true, "aarch64": true}
	if !valid[got] {
		t.Fatalf("unexpected adoptiumArch: %q", got)
	}
}

func TestSafeJoin_PreventTraversal(t *testing.T) {
	cases := []struct {
		name string
		dst  string
		path string
		ok   bool
	}{
		{"normal", "/tmp/dst", "subdir/file.txt", true},
		{"traversal", "/tmp/dst", "../../../etc/passwd", false},
		{"dot-dot-abs", "/tmp/dst", "../../other/file", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := safeJoin(c.dst, c.path)
			if c.ok && result == "" {
				t.Fatal("expected valid path, got empty")
			}
			if !c.ok && result != "" {
				t.Fatalf("expected empty (blocked), got %s", result)
			}
		})
	}
}

func TestVerifySHA256_Mismatch(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.bin")
	if err := os.WriteFile(f, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := verifySHA256(f, "0000000000000000000000000000000000000000000000000000000000000000")
	if err == nil {
		t.Fatal("expected SHA-256 mismatch error")
	}
}

func TestVerifySHA256_Correct(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.bin")
	if err := os.WriteFile(f, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

	// SHA-256 of "hello"
	err := verifySHA256(f, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824")
	if err != nil {
		t.Fatalf("expected valid SHA-256, got %v", err)
	}
}
