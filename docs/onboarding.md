# Onboarding — Sistema Runner

Guia de início rápido para instalar, configurar e usar o CLI `assinatura`
(assinador-cli) em Windows, Linux e macOS.

## 1. Instalação

### Baixar o binário

Baixe o binário da sua plataforma na página de
[Releases](https://github.com/gabrielmmendes/runner/releases):

| Plataforma | Artefato |
|------------|----------|
| Windows    | `assinatura-<versão>-windows-amd64.exe` |
| Linux      | `assinatura-<versão>-linux-amd64` |
| macOS      | `assinatura-<versão>-darwin-amd64` |

Cada release inclui também `SHA256SUMS`, e para cada binário os arquivos
`.sig` (assinatura Cosign) e `.pem` (certificado).

### Verificar integridade dos binários (recomendado)

**Checksum SHA-256:**

```bash
# Linux/macOS
sha256sum -c SHA256SUMS --ignore-missing

# Windows (PowerShell)
Get-FileHash .\assinatura-<versão>-windows-amd64.exe -Algorithm SHA256
# compare com a linha correspondente em SHA256SUMS
```

**Assinatura Cosign (keyless, transparency log):**

```bash
cosign verify-blob \
  --certificate assinatura-<versão>-linux-amd64.pem \
  --signature   assinatura-<versão>-linux-amd64.sig \
  --certificate-identity-regexp '^https://github.com/gabrielmmendes/runner/' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  assinatura-<versão>-linux-amd64
```

Instale o Cosign: <https://docs.sigstore.dev/cosign/installation/>.

### Tornar executável (Linux/macOS)

```bash
chmod +x assinatura-<versão>-linux-amd64
sudo mv assinatura-<versão>-linux-amd64 /usr/local/bin/assinatura
```

No Windows, adicione o `.exe` ao `PATH` ou execute pelo caminho completo.

## 2. Pré-requisitos

- **Java 21**: necessário para rodar o `assinador.jar`. Se ausente, o CLI
  provisiona automaticamente o JRE 21 (Adoptium Temurin) em `~/.hubsaude/jdk/`,
  com verificação SHA-256.
- **assinador.jar**: localizado nesta ordem:
  1. flag `--jar`
  2. env `ASSINATURA_JAR`
  3. mesmo diretório do executável
  4. `~/.hubsaude/assinador.jar`
  5. `target/` do projeto Maven (modo dev)

## 3. Primeiro uso

### Modo servidor (recomendado)

```bash
# inicia o assinador.jar em background
assinatura server start

# verifica status
assinatura server status

# assina (auto-start se o servidor estiver offline)
assinatura sign \
  --bundle bundle.json \
  --provenance provenance.json \
  --cert-chain cert-chain.pem \
  --timestamp 1751500000 \
  --strategy iat \
  --policy-id "https://policy.example/x|1.0.0" \
  --config config.json \
  --crypto-type smartcard \
  --pkcs11-id default

# encerra o servidor
assinatura server stop
```

### Observabilidade

- `--log-format json` — logs estruturados (campos `timestamp`, `level`,
  `message`, `command`, `version`) para Loki/ELK. Default: `text` (legível).
- `--log-level debug|info|warn|error` — verbosidade.
- `GET /metrics` no servidor — métricas Prometheus (`text/plain; version=0.0.4`):
  uptime, requisições por endpoint/status, erros, histograma de latência
  (p50/p95/p99 via `histogram_quantile`). Desabilite com
  `assinador.metrics.enabled=false`.

```bash
assinatura --log-format json --log-level debug server start
curl http://localhost:8080/metrics
```

## 4. Referência de comandos

| Comando | Descrição |
|---------|-----------|
| `assinatura sign [flags]` | Cria assinatura JWS (auto-start do jar) |
| `assinatura validate [flags]` | Valida assinatura |
| `assinatura server start` | Sobe o assinador.jar em background |
| `assinatura server status` | Status da instância registrada |
| `assinatura server stop` | Encerra a instância (graceful) |
| `assinatura version` | Versão do CLI |

### Flags globais

| Flag | Default | Descrição |
|------|---------|-----------|
| `--log-format` | `text` | `text` ou `json` |
| `--log-level` | `info` | `debug`/`info`/`warn`/`error` |

### Flags `server start`

| Flag | Default | Descrição |
|------|---------|-----------|
| `--jar` | (auto) | caminho do assinador.jar |
| `--port` | `8080` | porta HTTP |
| `--timeout` | `0` | minutos de inatividade até auto-stop (0 = off) |
| `--skip-verify` | `false` | ignora verificação Cosign do jar |

### Flags `sign` (principais)

| Flag | Obrigatório | Descrição |
|------|-------------|-----------|
| `--bundle` | sim | FHIR Bundle JSON |
| `--provenance` | sim | FHIR Provenance JSON |
| `--cert-chain` | sim | cadeia X.509 (PEM ou JSON base64) |
| `--timestamp` | sim | Unix UTC `[1751328000, 4102444800]` |
| `--strategy` | sim | `iat` ou `tsa` |
| `--policy-id` | sim | `<baseURI>|<major.minor.patch>` |
| `--config` | sim | JSON de operationalConfig |
| `--crypto-type` | sim | `smartcard` ou `token` |
| `--pkcs11-id` | sim | identificador da chave privada |
| `--pkcs11-pin` | não | PIN (ou env `ASSINATURA_PKCS11_PIN`, ou prompt) |
| `--pkcs11-slot` | não | slot ID PKCS#11 (≥0) |
| `--pkcs11-token-label` | não | label do token (≤32 chars) |
| `--service-url` | não | URL do assinador (default `http://localhost:8080`) |
| `--output` | não | arquivo p/ resposta JSON (default stdout) |
| `--timeout` | `30` | timeout HTTP em segundos |
| `--skip-verify` | `false` | ignora verificação Cosign do jar |

## 5. PKCS#11 com dispositivo real

Por padrão o servidor opera em modo simulado (chave RSA em memória). Para
usar um dispositivo real (smartcard/token), configure no `assinador.jar`:

```properties
# application.properties
assinador.pkcs11.library=/usr/lib/opensc-pkcs11.so   # caminho da lib do fabricante
```

Forneça o PIN via `--pkcs11-pin`, env `ASSINATURA_PKCS11_PIN`, ou prompt
interativo. O `--pkcs11-slot` e `--pkcs11-token-label` selecionam o dispositivo
quando há mais de um conectado.

Veja [troubleshooting.md](troubleshooting.md) para erros comuns.
