# Status de Implementação — Sistema Runner

## Sprint 1 — Fundação & CI/CD

**Progresso: ~85%**

| Item | Status |
|------|--------|
| Estrutura CLI Go (Cobra, `cmd/`, `internal/`, `main.go`) | ✅ |
| GitHub Actions build multiplataforma (`build.yml`) | ✅ |
| GitHub Actions release por tag SemVer (`release.yml`) | ✅ |
| Projeto Java base com `SignatureService`, `ValidationService`, `SignatureController` | ✅ |
| Endpoints `/api/sign` e `/api/validate` no Java | ✅ |
| Teste de integração Java básico (`SignatureControllerTest`) | ✅ |
| ADR inicial (`docs/adr/0001-escolhas-tecnologicas.md`) | ✅ |
| README raiz com visão geral, estrutura e convenções | ❌ |
| Convenção de commits e branches documentada | ❌ |
| Nome dos artefatos seguindo `assinatura-<versão>-<os>-<arch>` | ❌ (usa `runner-*`) |

---

## Sprint 2 — Assinatura Digital Simulada (Modo Local)

**Progresso: ~70%**

| Item | Status |
|------|--------|
| Comando `sign` com todos os flags (`--bundle`, `--provenance`, `--cert-chain`, `--timestamp`, `--strategy`, `--policy-id`, `--config`, `--crypto-type`, flags PKCS#11) | ✅ |
| Validação rigorosa de parâmetros (`internal/sign/validate.go`) | ✅ |
| JDK auto-provisionamento via Adoptium Temurin 21 com SHA-256 (`internal/java/jdk.go`) | ✅ |
| Cache do JDK em `~/.hubsaude/jdk/` | ✅ |
| Auto-start do assinador.jar no modo local | ✅ |
| Saída JSON para stdout ou arquivo via `--output` | ✅ |
| Comando `validate` no CLI | ❌ |
| Contrato JSON CLI↔jar documentado | ❌ |
| Testes de integração CLI → assinador.jar | ❌ |
| Testes cobrindo detecção de JDK presente/ausente | ❌ |

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

**Progresso: 0%**

| Item | Status |
|------|--------|
| SHA256SUMS + assinatura Cosign no `release.yml` ⚠️ marcado `[x]` no plano mas ausente no workflow real | ❌ |
| Verificação de integridade Cosign no download do jar | ❌ |
| `--log-format json` / `--log-level` nos CLIs | ❌ |
| Endpoint `/metrics` Prometheus no servidor Java | ❌ |
| Testes E2E no pipeline CI | ❌ |
| Documentação de onboarding e troubleshooting | ❌ |

---

## Resumo

| Sprint | Foco | Progresso |
|--------|------|-----------|
| 1 | Fundação & CI/CD | ~85% |
| 2 | Assinatura local | ~70% |
| 3 | Modo servidor HTTP | 100% |
| 4 | CLI simulador | 0% |
| 5 | Segurança & observabilidade | 0% |

### Próximos passos prioritários

1. Comando `validate` no CLI (Sprint 2)
2. Testes de integração CLI→jar local (Sprint 2)
3. SHA256SUMS + Cosign no `release.yml` (Sprint 5 — corrigir marcação incorreta no plano)
4. Início da Sprint 4 (CLI simulador + biblioteca compartilhada de processos)
