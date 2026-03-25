# 📌 Planejamento de Execução do Projeto Runner

## 🧭 1. Visão Geral

Este planejamento define a execução do projeto **Sistema Runner**, com abordagem incremental e iterativa, organizada por funcionalidades (User Stories) e focada na entrega progressiva de valor.

- Estratégia: Incremental / Ágil
- Organização: Por funcionalidades (User Stories)
- Prioridade: Núcleo funcional → Integrações → Automação → Distribuição

---

## 🚀 2. Fases do Projeto

### 🔹 Fase 1 — Setup e Arquitetura Base (Semana 1)

**Objetivo:** Preparar o ambiente e estrutura do projeto

#### Atividades

- Definir linguagem do CLI
- Criar repositórios:
  - `assinatura` (CLI)
  - `assinador` (Java)
- Estruturar diretórios e dependências
- Definir padrão de comandos CLI
- Criar pipeline básico de build

#### Entregáveis

- Estrutura inicial dos projetos
- README inicial

---

### 🔹 Fase 2 — Assinador (Core Java) (Semanas 2–3)

**Objetivo:** Implementar o núcleo do sistema

#### Atividades

- Implementar parsing de parâmetros
- Validar rigorosamente entradas
- Simular:
  - criação de assinatura
  - validação de assinatura
- Implementar tratamento de erros
- Criar modos:
  - CLI (execução direta)
  - HTTP (modo servidor)

#### Entregáveis

- `assinador.jar` funcional
- Testes unitários

---

### 🔹 Fase 3 — CLI (assinatura) + Integração (Semanas 3–4)

**Objetivo:** Permitir interação completa com o usuário

#### Atividades

- Criar comandos CLI:
  - `assinar`
  - `validar`
- Implementar integração:
  - Execução local
  - Execução via HTTP
- Detectar instância ativa do servidor
- Gerenciar ciclo de vida do assinador

#### Entregáveis

- CLI funcional integrado
- Testes de integração

---

### 🔹 Fase 4 — Simulador HubSaúde (Semana 5)

**Objetivo:** Gerenciar o simulador externo

#### Atividades

- Implementar download automático (GitHub Releases)
- Verificar versão local existente
- Criar comandos:
  - `start`
  - `stop`
  - `status`
- Validar portas disponíveis

#### Entregáveis

- Gerenciamento completo do simulador

---

### 🔹 Fase 5 — Provisionamento de JDK (Semana 6)

**Objetivo:** Automatizar dependência do Java

#### Atividades

- Detectar presença do JDK
- Baixar automaticamente quando necessário
- Configurar uso interno

#### Entregáveis

- Execução independente de instalação manual

---

### 🔹 Fase 6 — Empacotamento e Distribuição (Semana 7)

**Objetivo:** Preparar entregáveis finais

#### Atividades

- Gerar binários multiplataforma:
  - Windows
  - Linux
  - macOS
- Criar releases no GitHub
- Gerar checksums SHA256
- Aplicar versionamento SemVer

#### Entregáveis

- Binários distribuíveis

---

### 🔹 Fase 7 — Segurança e CI/CD (Semana 8)

**Objetivo:** Garantir integridade e automação

#### Atividades

- Configurar pipeline CI/CD
- Assinar artefatos com Cosign
- Publicar arquivos `.sig` e `.pem`
- Automatizar processo de release

#### Entregáveis

- Pipeline automatizado completo

---

### 🔹 Fase 8 — Testes e Documentação (Contínuo + Final)

#### Testes

- Testes unitários
- Testes de integração
- Casos de erro
- Testes de aceitação

#### Documentação

- Manual do usuário
- Guia de instalação
- Documentação técnica
- Exemplos de uso

---

## 📊 3. Priorização das Histórias de Usuário

1. HU-02 — Assinador (núcleo)
2. HU-01 — CLI + Integração
3. HU-03 — Simulador
4. HU-04 — Provisionamento JDK
5. HU-05 — Distribuição

---

## ⚠️ 4. Riscos e Mitigações

| Risco | Mitigação |
|------|----------|
| Integração CLI ↔ Java | Prototipação antecipada |
| Problemas com portas HTTP | Implementar fallback automático |
| Download de JDK | Utilizar fontes confiáveis |
| Multiplataforma | Testar via CI constantemente |
| Complexidade CI/CD | Iniciar pipeline cedo |

---

## 📅 5. Cronograma Resumido

| Semana | Entrega |
|------|--------|
| 1 | Setup |
| 2–3 | Assinador |
| 3–4 | CLI + Integração |
| 5 | Simulador |
| 6 | JDK automático |
| 7 | Binários |
| 8 | CI/CD + Finalização |

---

## ✅ 6. Resultado Esperado

Ao final do projeto, o sistema deverá:

- Executar aplicações Java sem configuração manual
- Disponibilizar uma CLI simples e intuitiva
- Integrar com o assinador via CLI e HTTP
- Gerenciar o simulador automaticamente
- Ser distribuído de forma segura e multiplataforma
