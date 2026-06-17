# ADR 0002 — Convenções Operacionais (porta, descoberta, segredos, logs)

## Status
Aceito — 2026-06-16

## Contexto
Várias decisões não óbvias precisam ficar registradas para garantir
rastreabilidade e comportamento previsível ("falhar bem"), em vez de ficarem
implícitas no código.

## Decisões

1. **Porta HTTP padrão `8080`** (`internal/java.DefaultPort`). Sobrescrita por
   `--port` no `server start` e por `--service-url`/`ASSINATURA_SERVICE_URL` no
   `sign`. Escolha alinhada ao default do Spring Boot.

2. **Descoberta do `assinador.jar`** (`FindJar`), em ordem de precedência:
   `--jar` → env `ASSINATURA_JAR` → `~/.hubsaude/assinador.jar`. O diretório
   `~/.hubsaude` também guarda o PID file (`assinador-java.pid`) que registra a
   instância em execução.

3. **PIN PKCS#11 nunca em flag obrigatória persistida**: precedência
   `--pkcs11-pin` → env `ASSINATURA_PKCS11_PIN` → prompt interativo no `stdin`.
   Evita vazamento de segredo no histórico de shell quando possível.

4. **Logs em `stderr`, resultado em `stdout`**: a resposta JSON da assinatura
   sai limpa no `stdout`, permitindo `| jq` e redirecionamento sem ruído de log.

5. **Versão rastreável**: `version.String()` combina tag + SHA curto + data,
   injetados via `-ldflags` no CI. Builds locais sem ldflags reportam
   `dev (none, built unknown)`.

## Consequências
- Comportamento previsível e documentado para operadores.
- Mensagens de erro explícitas com causa e remédio nos pontos de descoberta.
