# Sistema Runner

## 1. Visão Geral

O **Sistema Runner** é um trabalho prático desenvolvido para a disciplina de Implementação e Integração de Software do Bacharelado em Engenharia de Software (2026) da UFG. Este projeto é de interesse real da Secretaria de Estado de Saúde de Goiás (SES) e da Universidade Federal de Goiás (UFG), que realizam um esforço conjunto na definição e implementação de uma plataforma de interoperabilidade de dados em saúde.

O objetivo principal do sistema é facilitar o acesso à funcionalidade de execução de aplicações Java via linha de comandos, permitindo que os usuários executem essas aplicações sem a necessidade de conhecer detalhes de configuração ou instalação do ambiente Java.

## 2. Componentes do Sistema

* **Assinatura (CLI):** Uma interface de linha de comandos simples e intuitiva, multiplataforma (Windows, Linux e macOS).
* **Assinador (Java):** Aplicação `assinador.jar` que realiza a simulação de assinatura digital e validação rigorosa de parâmetros.
* **Simulador do HubSaúde:** Aplicação `simulador.jar` cujo ciclo de vida é gerenciado por este sistema.

## 3. Principais Funcionalidades

* **Execução Flexível:** O CLI pode invocar o Assinador de forma direta (linha de comando) ou via HTTP (modo servidor).
* **Provisionamento Automático de JDK:** O sistema detecta se o Java está instalado na máquina do usuário; se não estiver, baixa e configura a versão necessária automaticamente.
* **Simulação Criptográfica:** O Assinador simula a criação e a validação de assinaturas digitais, focando na validação rigorosa de parâmetros conforme especificações FHIR.
* **Gerenciamento do Simulador:** Permite iniciar, parar e monitorar o status do Simulador do HubSaúde, baixando-o dinamicamente caso não esteja disponível localmente.
* **Segurança e Integridade:** Os binários gerados são assinados criptograficamente utilizando Cosign/Sigstore, garantindo proteção contra adulterações e verificabilidade para o usuário final.

## 4. Estrutura do Repositório

```
Runner/
├── .github/
│   └── workflows/
│       ├── build.yml          # CI — build multiplataforma a cada push
│       └── release.yml        # CD — release por tag SemVer
├── projetos/
│   ├── assinador-cli/         # CLI em Go (Cobra)
│   │   ├── cmd/               # Comandos Cobra (sign, server, version)
│   │   ├── internal/          # Pacotes internos
│   │   │   ├── java/          # Provisionamento JDK e gestão de processos
│   │   │   ├── sign/          # Lógica de assinatura, validação, client HTTP
│   │   │   └── version/       # Versão injetada via ldflags
│   │   ├── docs/adr/          # Architecture Decision Records
│   │   ├── tests/             # Fixtures e helpers de teste
│   │   ├── main.go
│   │   └── go.mod
│   └── assinador-java/        # Aplicação Spring Boot (Java 21)
│       ├── src/
│       └── pom.xml
├── planejamento.md            # Plano de execução por fases
├── plano.md                   # User stories organizadas por sprint
├── status.md                  # Status de implementação
└── README.md
```

## 5. Tecnologias

| Componente | Tecnologia | Versão |
|------------|-----------|--------|
| CLI | Go | 1.25 |
| Framework CLI | Cobra | 1.8 |
| Assinador | Java (Spring Boot) | 21 |
| CI/CD | GitHub Actions | — |
| Distribuição | GitHub Releases | — |

Decisões de arquitetura registradas em [`projetos/assinador-cli/docs/adr/`](projetos/assinador-cli/docs/adr/).

## 6. Como Usar

### Pré-requisitos

- Go 1.25 (para compilar o CLI)
- JDK 21 (ou deixar o CLI provisionar automaticamente)
- Maven (para compilar o assinador.jar)

### Build

```bash
# CLI
cd projetos/assinador-cli
go build -o assinatura.exe

# Assinador JAR
cd projetos/assinador-java
mvn -B -DskipTests package
```

### Comandos principais

```bash
# Assinar documento (modo local — auto-start do jar)
assinatura sign --data "conteúdo" --cert cert.pem --key key.pem

# Modo servidor
assinatura server start --port 8085
assinatura server status
assinatura sign --data "conteúdo"       # usa servidor automaticamente
assinatura server stop

# Versão
assinatura version
```

## 7. Convenções

### Commits

O projeto segue [Conventional Commits](https://www.conventionalcommits.org/pt-br/):

```
<tipo>(<escopo>): <descrição curta>

[corpo opcional]

[rodapé opcional]
```

**Tipos:**

| Tipo | Uso |
|------|-----|
| `feat` | Nova funcionalidade |
| `fix` | Correção de bug |
| `docs` | Documentação |
| `refactor` | Refatoração sem mudança de comportamento |
| `test` | Adição ou correção de testes |
| `chore` | Tarefas auxiliares (CI, build, deps) |
| `style` | Formatação, sem mudança de lógica |

**Escopos comuns:** `cli`, `java`, `ci`, `docs`

**Exemplos:**
```
feat(cli): add server start command with --port flag
fix(java): handle null pointer in SignatureService
chore(ci): update Go version to 1.25
```

### Branches

| Padrão | Uso |
|--------|-----|
| `main` | Branch principal — sempre estável |
| `feat/<nome>` | Nova funcionalidade |
| `fix/<nome>` | Correção de bug |
| `chore/<nome>` | Manutenção, CI, deps |
| `docs/<nome>` | Documentação |

### Versionamento

Segue [SemVer](https://semver.org/lang/pt-BR/). Tags no formato `v<MAJOR>.<MINOR>.<PATCH>` (ex.: `v0.1.0`).

### Artefatos

Binários nomeados como `assinatura-<versão>-<os>-<arch>` (ex.: `assinatura-v0.2.0-linux-amd64`).

## 8. Contexto Acadêmico

|||
|:--|:--|
|**Instituição:**| Universidade Federal de Goiás (UFG) - Instituto de Informática |
|**Curso:** | Bacharelado em Engenharia de Software|
|**Disciplina:** | Implementação e Integração de Software (INF0466) |
|**Professor:**|Fabio Nogueira de Lucena|
|**Semestre Letivo:**|2026/1|
|**Estudantes:**| Gabriel Matos Mendes, Eduardo Divino Miranda Pereira|
