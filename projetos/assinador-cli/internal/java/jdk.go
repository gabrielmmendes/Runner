package java

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const adoptiumAPI = "https://api.adoptium.net/v3/assets/latest/21/hotspot"

type adoptiumAsset struct {
	Binary struct {
		Package struct {
			Link     string `json:"link"`
			Checksum string `json:"checksum"`
			Size     int64  `json:"size"`
		} `json:"package"`
	} `json:"binary"`
}

func adoptiumOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

func adoptiumArch() string {
	if runtime.GOARCH == "arm64" {
		return "aarch64"
	}
	return "x64"
}

// hubJavaBin returns the java binary inside ~/.hubsaude/jdk/ if provisioned.
func hubJavaBin() (string, error) {
	hub, err := HubDir()
	if err != nil {
		return "", err
	}
	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}
	jdkBase := filepath.Join(hub, "jdk")
	entries, err := os.ReadDir(jdkBase)
	if err != nil {
		return "", fmt.Errorf("~/.hubsaude/jdk não existe")
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := filepath.Join(jdkBase, e.Name(), "bin", bin)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("java não encontrado em ~/.hubsaude/jdk/")
}

// EnsureJava returns a java binary path, provisioning JRE 21 if necessary.
// Progress messages are written to out (typically os.Stderr).
func EnsureJava(out io.Writer) (string, error) {
	// 1. system java (JAVA_HOME, PATH)
	if p, err := FindJava(); err == nil {
		return p, nil
	}
	// 2. previously provisioned JRE
	if p, err := hubJavaBin(); err == nil {
		return p, nil
	}
	// 3. provision
	fmt.Fprintln(out, "Java não encontrado — baixando JRE 21 (Adoptium Temurin)...")
	return provisionJRE(out)
}

func provisionJRE(out io.Writer) (string, error) {
	apiURL := fmt.Sprintf("%s?os=%s&arch=%s&image_type=jre",
		adoptiumAPI, adoptiumOS(), adoptiumArch())

	asset, err := fetchAdoptiumMeta(apiURL)
	if err != nil {
		return "", fmt.Errorf("metadados JRE 21: %w", err)
	}

	hub, err := HubDir()
	if err != nil {
		return "", err
	}
	jdkBase := filepath.Join(hub, "jdk")
	if err := os.MkdirAll(jdkBase, 0o700); err != nil {
		return "", fmt.Errorf("criar ~/.hubsaude/jdk: %w", err)
	}

	tmp, err := os.CreateTemp(hub, "jre-dl-*")
	if err != nil {
		return "", fmt.Errorf("arquivo temporário: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	sizeMB := float64(asset.Binary.Package.Size) / 1e6
	fmt.Fprintf(out, "baixando JRE 21 (%.0f MB) ...\n", sizeMB)
	if err := downloadWithProgress(out, asset.Binary.Package.Link, tmp, asset.Binary.Package.Size); err != nil {
		tmp.Close()
		return "", fmt.Errorf("download: %w", err)
	}
	tmp.Close()

	fmt.Fprintln(out, "verificando integridade (SHA-256)...")
	if err := verifySHA256(tmpPath, asset.Binary.Package.Checksum); err != nil {
		return "", err
	}

	fmt.Fprintln(out, "extraindo JRE...")
	if err := extractArchive(tmpPath, jdkBase); err != nil {
		return "", fmt.Errorf("extração: %w", err)
	}

	p, err := hubJavaBin()
	if err != nil {
		return "", fmt.Errorf("JRE extraído mas binário não encontrado: %w", err)
	}
	fmt.Fprintf(out, "JRE 21 provisionado: %s\n", p)
	return p, nil
}

func fetchAdoptiumMeta(apiURL string) (*adoptiumAsset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Adoptium API HTTP %d", resp.StatusCode)
	}
	var assets []adoptiumAsset
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		return nil, fmt.Errorf("nenhum JRE 21 disponível para %s/%s", adoptiumOS(), adoptiumArch())
	}
	return &assets[0], nil
}

func downloadWithProgress(out io.Writer, url string, dst *os.File, totalBytes int64) error {
	resp, err := http.Get(url) //nolint:noctx // download sem timeout intencional
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download HTTP %d", resp.StatusCode)
	}
	buf := make([]byte, 64*1024)
	var downloaded int64
	lastPct := -1
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				return werr
			}
			downloaded += int64(n)
			if totalBytes > 0 {
				pct := int(downloaded * 100 / totalBytes)
				if pct != lastPct && pct%10 == 0 {
					fmt.Fprintf(out, "  %d%%\n", pct)
					lastPct = pct
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func verifySHA256(path, expected string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got, expected) {
		return fmt.Errorf("SHA-256 inválido: esperado %s obtido %s", expected, got)
	}
	return nil
}

func extractArchive(src, dst string) error {
	if strings.HasSuffix(src, ".zip") {
		return extractZip(src, dst)
	}
	return extractTarGz(src, dst)
}

func extractZip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if err := writeZipEntry(f, dst); err != nil {
			return err
		}
	}
	return nil
}

func writeZipEntry(f *zip.File, dst string) error {
	target := safeJoin(dst, f.Name)
	if target == "" {
		return nil
	}
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

func extractTarGz(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := safeJoin(dst, hdr.Name)
		if target == "" {
			continue
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0o755)
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		case tar.TypeSymlink:
			// best-effort — ignore errors (symlinks may not work on all platforms)
			os.Symlink(hdr.Linkname, target)
		}
	}
	return nil
}

// safeJoin prevents path traversal — returns "" for paths escaping dst.
func safeJoin(dst, name string) string {
	clean := filepath.Clean(filepath.Join(dst, filepath.FromSlash(name)))
	dstClean := filepath.Clean(dst)
	if clean != dstClean && !strings.HasPrefix(clean, dstClean+string(os.PathSeparator)) {
		return ""
	}
	return clean
}
