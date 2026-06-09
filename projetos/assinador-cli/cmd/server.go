package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gabrielmmendes/runner/internal/java"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Gerencia o processo assinador-java",
}

var (
	serverJar         string
	serverPort        int
	serverStopPort    int
	serverIdleTimeout int
	serverSkipVerify  bool
)

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia assinador.jar em background (modo servidor HTTP)",
	RunE:  runServerStart,
}

var serverStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra instância registrada do assinador.jar",
	RunE:  runServerStop,
}

var serverStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe status da instância registrada",
	RunE:  runServerStatus,
}

func init() {
	f := serverStartCmd.Flags()
	f.StringVar(&serverJar, "jar", "", "caminho para assinador.jar (ou env ASSINATURA_JAR)")
	f.IntVar(&serverPort, "port", java.DefaultPort, "porta HTTP do servidor")
	f.IntVar(&serverIdleTimeout, "timeout", 0,
		"minutos de inatividade antes de auto-stop (0 = desativado)")
	f.BoolVar(&serverSkipVerify, "skip-verify", false,
		"ignora verificação Cosign do jar (não recomendado)")

	sf := serverStopCmd.Flags()
	sf.IntVar(&serverStopPort, "port", 0,
		"porta da instância a encerrar (0 = usa registro PID)")

	serverCmd.AddCommand(serverStartCmd, serverStopCmd, serverStatusCmd)
	rootCmd.AddCommand(serverCmd)
}

func runServerStart(cmd *cobra.Command, _ []string) error {
	javaPath, err := java.EnsureJava(os.Stderr)
	if err != nil {
		return err
	}
	jarPath, err := java.FindJar(serverJar)
	if err != nil {
		return err
	}
	if err := verifyJar(jarPath, serverSkipVerify); err != nil {
		return err
	}

	baseURL := fmt.Sprintf("http://localhost:%d", serverPort)
	if java.IsRunning(baseURL) {
		fmt.Fprintf(os.Stderr, "assinador-java já está UP em %s\n", baseURL)
		return nil
	}

	pid, err := java.Start(java.StartOptions{
		JavaPath:       javaPath,
		JarPath:        jarPath,
		Port:           serverPort,
		IdleTimeoutMin: serverIdleTimeout,
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "iniciado PID=%d — aguardando %s ...\n", pid, baseURL)

	if !java.WaitReady(baseURL, 90*time.Second) {
		return fmt.Errorf("assinador-java não respondeu em 90s — verifique logs")
	}
	fmt.Printf("assinador-java UP em %s (PID=%d)\n", baseURL, pid)
	return nil
}

func runServerStop(_ *cobra.Command, _ []string) error {
	if err := java.Stop(serverStopPort); err != nil {
		return err
	}
	fmt.Println("assinador-java encerrado")
	return nil
}

func runServerStatus(_ *cobra.Command, _ []string) error {
	running, pid, port := java.Status()
	if !running {
		if pid != 0 {
			fmt.Printf("registrado PID=%d porta=%d mas não responde\n", pid, port)
		} else {
			fmt.Println("nenhuma instância registrada")
		}
		return nil
	}
	fmt.Printf("UP — PID=%d porta=%d url=http://localhost:%d\n", pid, port, port)
	return nil
}
