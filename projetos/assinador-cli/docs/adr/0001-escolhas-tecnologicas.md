# ADR 0001 — Escolhas Tecnológicas

## Status
Aceito — 2026-04-28

## Contexto
O sistema precisa de um executável multiplataforma (Windows, Linux, macOS) que
o usuário rode sem conhecer detalhes do ambiente Java, integrando-se a um
serviço de assinatura escrito em Java.

## Decisão

- **Go 1.25** para o CLI: binário único estático, cross-compiling nativo
  (`GOOS`/`GOARCH`), sem runtime a instalar no cliente.
- **Cobra** como framework CLI: subcomandos, geração de `--help`/`--version` e
  autocompletar, padrão de fato no ecossistema Go.
- **Java 21 (Spring Boot)** para o `assinador.jar`: reaproveita bibliotecas
  criptográficas/FHIR maduras do ecossistema Java.

## Alternativas consideradas
- CLI em Java (uberjar): exigiria JRE no cliente — descartado.
- CLI em Rust: cross-compiling e curva de adoção menos convenientes para a equipe.

## Consequências
- Build simples e rápido; distribuição de um único binário por plataforma.
- Separação clara entre CLI (orquestração) e engine (assinatura).
- O CLI precisa provisionar/descobrir o JDK em runtime (ver [ADR 0002](0002-convencoes-operacionais.md)).
