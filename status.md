# Status de Implementação — Sistema Runner

## Sprint 1 — Fundação & CI/CD

**Progresso: 100%**

| Item | Status |
|------|--------|
| Estrutura CLI Go (Cobra, `cmd/`, `internal/`, `main.go`) | ✅ |
| GitHub Actions build multiplataforma (`build.yml`) | ✅ |
| GitHub Actions release por tag SemVer (`release.yml`) | ✅ |
| Projeto Java base com `SignatureService`, `ValidationService`, `SignatureController` | ✅ |
| Endpoints `/api/sign` e `/api/validate` no Java | ✅ |
| Teste de integração Java básico (`SignatureControllerTest`) | ✅ |
| ADR inicial (`docs/adr/0001-escolhas-tecnologicas.md`) | ✅ |
| README raiz com visão geral, estrutura e convenções | ✅ |
| Convenção de commits e branches documentada | ✅ |
| Nome dos artefatos seguindo `assinatura-<versão>-<os>-<arch>` | ✅ |

---

## Sprint 2 — Assinatura Digital Simulada (Modo Local)

**Progresso: 100%**

| Item | Status |
|------|--------|
| Comando `sign` com todos os flags (`--bundle`, `--provenance`, `--cert-chain`, `--timestamp`, `--strategy`, `--policy-id`, `--config`, `--crypto-type`, flags PKCS#11) | ✅ |
| Validação rigorosa de parâmetros (`internal/sign/validate.go`) | ✅ |
| JDK auto-provisionamento via Adoptium Temurin 21 com SHA-256 (`internal/java/jdk.go`) | ✅ |
| Cache do JDK em `~/.hubsaude/jdk/` | ✅ |
| Auto-start do assinador.jar no modo local | ✅ |
| Saída JSON para stdout ou arquivo via `--output` | ✅ |
| Comando `validate` no CLI (`cmd/validate.go`) | ✅ |
| Contrato JSON CLI↔jar documentado (`docs/contrato-json.md`) | ✅ |
| Testes de integração CLI → assinador.jar (`internal/sign/integration_test.go`) | ✅ |
| Testes cobrindo detecção de JDK presente/ausente (`internal/java/jdk_test.go`) | ✅ |

---

## Sprint 3 — Modo Servidor HTTP & Material Criptográfico

**Progresso: 100%**

| Item | Status |
|------|--------|
| Comandos `server start`, `server stop`, `server status` (`cmd/server.go`) | ✅ |
| Detecção de instância ativa via `java.IsRunning()` antes de auto-start | ✅ |
| HTTP client para `/sign` (`internal/sign/client.go`) | ✅ |
| Fallback para modo local quando servidor offline | ✅ |
| PKCS#11 via `SunPKCS11` no Java (fallback in-memory quando `assinador.pkcs11.library` vazio) | ✅ |
| `--port` no comando `server stop` | ✅ |
| `--timeout` (auto-stop por inatividade) no `server start` | ✅ |
| Graceful shutdown (SIGTERM) no servidor Java (`server.shutdown=graceful`) | ✅ |
| Testes de integração CLI → HTTP → assinador.jar (`-tags=integration`) | ✅ |

---

## Sprint 4 — CLI Simulador & Gestão de Artefatos

**Progresso: 0%**

| Item | Status |
|------|--------|
| Projeto CLI `simulador` | ❌ |
| Pacote Go compartilhado de gestão de processos (PID, health check, registro) | ❌ |
| Download automático do simulador.jar (GitHub Releases + SHA-256) | ❌ |
| Comandos `simulador start`, `simulador stop`, `simulador status` | ❌ |
| Comando `update` (auto-atualização do binário) nos dois CLIs | ❌ |

---

## Sprint 5 — Segurança de Artefatos & Observabilidade

**Progresso: 100%**

| Item | Status |
|------|--------|
| SHA256SUMS + assinatura Cosign keyless (OIDC + transparency log) no `release.yml` | ✅ |
| Verificação de integridade Cosign do jar antes do auto-start (`internal/verify`, sidecars `.sig`/`.pem`, fail-closed) | ✅ |
| Flag `--skip-verify` (com aviso) em `sign` e `server start` | ✅ |
| `--log-format text\|json` / `--log-level` globais nos CLIs (`internal/logging`, slog; campos timestamp/level/message/command/version) | ✅ |
| Endpoint `/metrics` Prometheus no servidor Java (uptime, requests por endpoint/status, erros, histograma de latência) — `assinador.metrics.enabled` | ✅ |
| Testes E2E no pipeline CI (`e2e.yml`: jar Maven + CLI → sign/validate/metrics + lifecycle server) | ✅ |
| Documentação de onboarding e troubleshooting (`docs/onboarding.md`, `docs/troubleshooting.md`) | ✅ |

---

## Resumo

| Sprint | Foco | Progresso |
|--------|------|-----------|
| 1 | Fundação & CI/CD | 100% |
| 2 | Assinatura local | 100% |
| 3 | Modo servidor HTTP | 100% |
| 4 | CLI simulador | 0% |
| 5 | Segurança & observabilidade | 100% |

### Próximos passos prioritários

1. Sprint 4 (CLI simulador + biblioteca compartilhada de processos + auto-update)
