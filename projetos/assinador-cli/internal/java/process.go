package java

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultPort = 8080
	JarName     = "assinador.jar"
	hubDirName  = ".hubsaude"
	pidFileName = "assinador-java.pid"
)

func HubDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir inacessível: %w", err)
	}
	return filepath.Join(home, hubDirName), nil
}

func FindJava() (string, error) {
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		bin := "java"
		if runtime.GOOS == "windows" {
			bin = "java.exe"
		}
		p := filepath.Join(jh, "bin", bin)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	p, err := exec.LookPath("java")
	if err != nil {
		return "", fmt.Errorf("java não encontrado (instale JDK 21 ou defina JAVA_HOME)")
	}
	return p, nil
}

func FindJar(jarFlag string) (string, error) {
	if jarFlag != "" {
		if _, err := os.Stat(jarFlag); err != nil {
			return "", fmt.Errorf("jar não encontrado: %s", jarFlag)
		}
		return jarFlag, nil
	}
	if env := os.Getenv("ASSINATURA_JAR"); env != "" {
		if _, err := os.Stat(env); err == nil {
			return env, nil
		}
	}
	// same dir as exe
	if exe, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exe), JarName)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	// ~/.hubsaude/assinador.jar
	if hub, err := HubDir(); err == nil {
		p := filepath.Join(hub, JarName)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	// dev: look for assinador-*.jar in Maven target/ relative to cwd and exe
	if p, err := findMavenJar(); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("assinador.jar não encontrado (use --jar, env ASSINATURA_JAR, ou coloque em ~/.hubsaude/)")
}

// findMavenJar searches for the built jar in Maven target/ directories,
// relative to the current working directory and the executable path.
func findMavenJar() (string, error) {
	roots := []string{}
	if cwd, err := os.Getwd(); err == nil {
		roots = append(roots, cwd)
	}
	if exe, err := os.Executable(); err == nil {
		roots = append(roots, filepath.Dir(exe))
	}
	candidates := []string{
		filepath.Join("assinador-java", "target"),
		filepath.Join("..", "assinador-java", "target"),
		filepath.Join("..", "..", "assinador-java", "target"),
	}
	for _, root := range roots {
		for _, rel := range candidates {
			dir := filepath.Join(root, rel)
			entries, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range entries {
				name := e.Name()
				if !e.IsDir() && strings.HasSuffix(name, ".jar") &&
					!strings.Contains(name, "-sources") && !strings.Contains(name, "-javadoc") {
					p := filepath.Join(dir, name)
					return p, nil
				}
			}
		}
	}
	return "", fmt.Errorf("jar Maven não encontrado")
}

func pidFilePath() (string, error) {
	hub, err := HubDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(hub, pidFileName), nil
}

func writePID(pid, port int) error {
	hub, err := HubDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(hub, 0o700); err != nil {
		return fmt.Errorf("criar ~/.hubsaude: %w", err)
	}
	pf, err := pidFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(pf, []byte(fmt.Sprintf("%d\n%d\n", pid, port)), 0o600)
}

func readPID() (pid, port int, err error) {
	pf, err := pidFilePath()
	if err != nil {
		return 0, 0, err
	}
	data, err := os.ReadFile(pf)
	if err != nil {
		return 0, 0, err
	}
	lines := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
	if len(lines) < 2 {
		return 0, 0, fmt.Errorf("pid file corrompido")
	}
	pid, err = strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("pid inválido: %w", err)
	}
	port, err = strconv.Atoi(strings.TrimSpace(lines[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("port inválido: %w", err)
	}
	return pid, port, nil
}

func clearPID() {
	pf, _ := pidFilePath()
	os.Remove(pf)
}

// IsRunning returns true if HTTP server answers at baseURL.
func IsRunning(baseURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		strings.TrimRight(baseURL, "/"), nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

// WaitReady polls IsRunning until server answers or timeout elapses.
func WaitReady(baseURL string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if IsRunning(baseURL) {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

// StartOptions configures a Start call.
type StartOptions struct {
	JavaPath       string
	JarPath        string
	Port           int
	IdleTimeoutMin int       // 0 = sem auto-stop
	LogOut         io.Writer // nil = discard
}

// Start launches assinador.jar as a detached background process.
func Start(opts StartOptions) (pid int, err error) {
	out := io.Discard
	if opts.LogOut != nil {
		out = opts.LogOut
	}
	args := []string{"-jar", opts.JarPath,
		fmt.Sprintf("--server.port=%d", opts.Port)}
	if opts.IdleTimeoutMin > 0 {
		args = append(args,
			fmt.Sprintf("--assinador.idle-timeout-min=%d", opts.IdleTimeoutMin))
	}
	cmd := exec.Command(opts.JavaPath, args...)
	cmd.Stdout = out
	cmd.Stderr = out
	setSysProcAttr(cmd)
	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("falha ao iniciar assinador.jar: %w", err)
	}
	_ = writePID(cmd.Process.Pid, opts.Port)
	return cmd.Process.Pid, nil
}

// Stop sends interrupt to the registered process and clears the PID file.
// If port > 0, only stops when registered port matches.
func Stop(port int) error {
	pid, regPort, err := readPID()
	if err != nil {
		return fmt.Errorf("nenhuma instância registrada: %w", err)
	}
	if port > 0 && port != regPort {
		return fmt.Errorf("instância registrada usa porta %d, não %d", regPort, port)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		clearPID()
		return fmt.Errorf("processo %d não encontrado", pid)
	}
	if err := proc.Signal(os.Interrupt); err != nil {
		_ = proc.Kill()
	}
	clearPID()
	return nil
}

// Status reads the PID file and checks if the server is responding.
func Status() (running bool, pid, port int) {
	pid, port, err := readPID()
	if err != nil {
		return false, 0, 0
	}
	baseURL := fmt.Sprintf("http://localhost:%d", port)
	return IsRunning(baseURL), pid, port
}
