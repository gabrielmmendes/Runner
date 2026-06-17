# Assinatura CLI

CLI multiplataforma (Go + Cobra) que orquestra o serviço `assinador-java` para
criar e validar assinaturas digitais FHIR. Provisiona o JDK quando ausente e
gerencia o `assinador.jar` em modo local (auto-start) ou servidor HTTP.

> Especificação (fonte da verdade): [`kyriosdata/runner@0cd5461`](https://github.com/kyriosdata/runner/tree/0cd5461481861b320b6e6d4f9af85648206cea56).

## O que é

| Comando | Função |
|---------|--------|
| `sign` | Cria assinatura JWS a partir de Bundle + Provenance + cadeia X.509 |
| `validate` | Valida parâmetros/assinatura |
| `verify` | Verifica assinatura Cosign do `assinador.jar` |
| `server start\|stop\|status` | Ciclo de vida do `assinador.jar` como servidor HTTP |
| `version` | Identificador rastreável (tag + SHA curto + data) |

## Gerar o executável

```bash
go build -o assinatura .

# Build rastreável (tag + SHA injetados via ldflags)
PKG=github.com/gabrielmmendes/runner/internal/version
go build -ldflags="-X $PKG.Version=$(git describe --tags --always) -X $PKG.Commit=$(git rev-parse --short HEAD)" -o assinatura .
```

## Executar o artefato

```bash
./assinatura --version      # tag + SHA + data
./assinatura --help         # ajuda com exemplos
./assinatura sign --help    # exemplos por comando
```

Flags globais de log: `--log-format text|json`, `--log-level`,
`--verbose`/`-v` (debug), `--quiet`/`-q` (apenas erros). Logs vão para `stderr`;
o `stdout` fica reservado à resposta JSON.

## Rodar os testes

```bash
go test ./...                                # unitários
go test -tags=integration ./internal/sign/... # integração (sobe o jar real)
go vet ./...                                 # análise estática
gofmt -l .                                   # formatação (saída vazia = ok)
```

## Estrutura

- `cmd/` → comandos Cobra (`sign`, `server`, `validate`, `verify`, `version`)
- `internal/java/` → provisionamento JDK e gestão de processos
- `internal/sign/` → carga de payloads, validação e client HTTP
- `internal/logging/` → logging estruturado (slog)
- `internal/version/` → identificadores de build injetados via ldflags
- `docs/adr/` → Architecture Decision Records

## Contribuir

Branches `feat/*`, `fix/*`, `chore/*`, `docs/*`. Commits seguem
[Conventional Commits](https://www.conventionalcommits.org/pt-br/). PRs pequenos,
vinculados a issues; CI (lint + testes Linux/Windows + build) precisa passar.

## Licença

[MIT](../../LICENSE).
