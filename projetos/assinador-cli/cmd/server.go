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
	serverJar  string
	serverPort int
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

	baseURL := fmt.Sprintf("http://localhost:%d", serverPort)
	if java.IsRunning(baseURL) {
		fmt.Fprintf(os.Stderr, "assinador-java já está UP em %s\n", baseURL)
		return nil
	}

	pid, err := java.Start(java.StartOptions{
		JavaPath: javaPath,
		JarPath:  jarPath,
		Port:     serverPort,
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
	if err := java.Stop(); err != nil {
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
