# ADR 0001 — Escolhas Tecnológicas

## Status
Aceito

## Contexto
Necessidade de CLI multiplataforma e integração com Java.

## Decisão

- Go 1.25 → CLI
- Cobra → framework CLI
- Java 21 → assinador.jar

## Consequências

- Build simples e rápido
- Cross-compiling nativo
- Separação clara entre CLI e engine
