# Plano de Desenvolvimento — Reorganizado por Prioridade

## Premissas

- CLIs desenvolvidos em **Go 1.25** — cross-compiling nativo e biblioteca padrão robusta.
- **assinador.jar** desenvolvido em **Java 21** — restrição de projeto.
- Estratégia **iterativa e incremental**, organizada em **5 sprints de 1 semana**.
- Cada sprint entrega valor ao usuário ou remove risco técnico relevante.
- Histórias nomeadas como `US-XX.Y`, indicando o épico de origem.

> **Nota sobre a reorganização:** O plano original iniciava com modo servidor (funcionalidade avançada) antes de ter o modo local funcionando, posicionava CI/CD apenas no Sprint 3 e duplicava os Sprints 2 e 4. A ordem abaixo prioriza fundação → valor principal → otimização → expansão → segurança.

---

## Rastreabilidade Épicos → Histórias

| Épico | Descrição | Histórias derivadas |
|-------|-----------|---------------------|
| US-01 | Gerenciar Ciclo de Vida do Simulador | US-01.1, US-01.2, US-01.3, US-01.4, US-01.5, US-01.6, US-01.7 |
| US-02 | Simular assinatura digital com validação | US-02.1, US-02.2, US-02.3, US-02.4, US-02.5, US-02.6, US-02.7 |
| US-03 | Invocar assinador.jar via CLI | US-03.1, US-03.2, US-03.3, US-03.4 |
| US-05 | Disponibilizar binários multiplataforma | US-05.1, US-05.2, US-05.3, US-05.4 |

---

## Sprint 1 — Fundação & Entrega Contínua

**Objetivo:** Garantir que qualquer contribuição de código gera artefatos funcionais e verificáveis. Sem essa base, nenhum sprint posterior pode ser validado com segurança.

**Valor entregue:** Pipeline ativo desde o primeiro commit. Estrutura do projeto definida. Projeto Java base compilável.

### US-03.1 — Estrutura base do CLI em Go

**Como** usuário do Sistema Runner,
**quero** que o projeto CLI esteja estruturado com organização de pacotes e build funcional,
**para que** o desenvolvimento possa progredir de forma organizada e incremental.

**Critérios de aceitação:**
- [x] Projeto Go inicializado com `go mod init github.com/gabrielmmendes/assinatura`
- [x] Cobra CLI instalado: `go install github.com/spf13/cobra-cli@latest`
- [x] Estrutura de pacotes definida e documentada
- [x] Aplicação compila e executa nas três plataformas (Windows, Linux, macOS)
- [x] Comando `assinatura version` exibe a versão atual do CLI

### US-03.2 — Pipeline CI/CD multiplataforma

**Como** desenvolvedor do Sistema Runner,
**quero** que alterações no repositório disparem automaticamente a compilação para Windows, Linux e macOS,
**para que** binários atualizados estejam sempre disponíveis após cada mudança.

**Critérios de aceitação:**
- [x] GitHub Actions configurado com workflow de build
- [x] Cross-compilation para `windows/amd64`, `linux/amd64` e `darwin/amd64`
- [x] Build executado a cada push na branch principal
- [x] Artefatos de build disponíveis como artifacts do workflow

### US-03.3 — Publicação de releases com versionamento semântico

**Como** usuário do Sistema Runner,
**quero** baixar binários pré-compilados para minha plataforma a partir do GitHub Releases,
**para que** eu possa utilizar o sistema sem necessidade de compilação, com versionamento claro.

**Critérios de aceitação:**
- [x] Tags de versão seguem SemVer (ex.: `v0.1.0`)
- [x] Workflow de release gera binários nomeados por plataforma
- [x] Binários publicados automaticamente no GitHub Releases ao criar tag
- [x] Nome dos artefatos segue convenção: `assinatura-<versão>-<os>-<arch>`

### US-02.1 — Projeto Java base (assinador.jar)

**Como** usuário do Sistema Runner,
**quero** que o assinador.jar retorne uma assinatura simulada quando receber parâmetros válidos,
**para que** eu possa testar o fluxo de assinatura sem infraestrutura criptográfica real.

**Critérios de aceitação:**
- [x] Projeto Java base inicializado no diretório `projetos/assinador-java`
- [x] Interface `SignatureService` definida com métodos `sign` e `validate`
- [x] Implementação `FakeSignatureService` retorna assinatura pré-construída para parâmetros válidos
- [x] Resposta simulada inclui os campos esperados conforme especificação
- [x] Testes unitários cobrem o cenário de sucesso

### NOVO — ADR e convenções do repositório

**Como** desenvolvedor do Sistema Runner,
**quero** que o repositório tenha convenções claras documentadas desde o início,
**para que** toda a equipe trabalhe com consistência ao longo dos sprints.

**Critérios de aceitação:**
- [ ] README raiz com visão geral do projeto, estrutura de diretórios e links para documentação
- [ ] Estrutura `projetos/` documentada
- [ ] Convenção de commits (ex.: Conventional Commits) definida
- [ ] Convenção de branches (ex.: `feat/`, `fix/`, `chore/`) definida
- [ ] ADR inicial registrando escolhas de tecnologia (Go, Java 21, Cobra)

---

## Sprint 2 — Assinatura Digital Simulada (Modo Local)

**Objetivo:** Entregar o fluxo ponta-a-ponta: o usuário executa `assinatura sign` e obtém uma assinatura simulada, sem precisar instalar o Java manualmente.

**Valor entregue:** Caso de uso principal funcionando. O CLI detecta, provisiona e invoca o assinador.jar de forma transparente.

### NOVO — Contrato JSON CLI ↔ assinador.jar

**Como** desenvolvedor do Sistema Runner,
**quero** que o contrato de entrada e saída entre CLI e assinador.jar esteja formalmente documentado,
**para que** o modo HTTP (Sprint 3) possa reutilizar o mesmo schema sem retrabalho.

**Critérios de aceitação:**
- [ ] Schema JSON de entrada para `sign` e `validate` documentado
- [ ] Schema JSON de saída (sucesso e erro) documentado
- [ ] Documento de contrato versionado junto ao repositório
- [ ] Testes de integração validam conformidade com o contrato

### US-02.4 — Parsing de comandos e parâmetros no CLI

**Como** usuário do Sistema Runner,
**quero** executar comandos `sign` e `validate` com parâmetros via linha de comandos,
**para que** eu possa solicitar operações de assinatura de forma intuitiva.

**Critérios de aceitação:**
- [ ] CLI aceita o comando `sign` com os parâmetros necessários
- [ ] CLI aceita o comando `validate` com os parâmetros necessários
- [ ] Mensagem de ajuda (`--help`) documenta os comandos e parâmetros disponíveis
- [ ] Parâmetros ausentes ou inválidos geram mensagem de erro orientativa
- [ ] Testes cobrem o parsing de comandos e parâmetros

### US-02.2 — Validação de parâmetros de criação de assinatura

**Como** usuário do Sistema Runner,
**quero** que o assinador.jar valide rigorosamente os parâmetros de criação de assinatura,
**para que** eu receba feedback imediato e claro sobre erros antes da operação ser processada.

**Critérios de aceitação:**
- [ ] Todos os parâmetros obrigatórios são verificados (presença e formato)
- [ ] Mensagens de erro indicam qual parâmetro está inválido e o motivo
- [ ] Parâmetros inválidos são rejeitados antes de qualquer processamento
- [ ] Testes unitários cobrem todos os cenários de validação

### US-02.3 — Simulação e validação de parâmetros de validação de assinatura

**Como** usuário do Sistema Runner,
**quero** que o assinador.jar valide os parâmetros de validação de assinatura e retorne resultado pré-determinado,
**para que** eu possa testar o fluxo de validação com feedback claro sobre parâmetros incorretos.

**Critérios de aceitação:**
- [ ] Parâmetros de validação são verificados (presença e formato)
- [ ] Resultado pré-determinado (válido/inválido) retornado baseado em critérios simples
- [ ] Mensagens de erro claras para parâmetros inválidos
- [ ] Testes unitários cobrem cenários de sucesso e falha

### US-02.7 — Detecção e provisionamento automático do JDK

**Como** usuário do Sistema Runner,
**quero** que o sistema detecte se o JDK compatível está presente e, caso não esteja, baixe e configure automaticamente,
**para que** eu possa utilizar o Assinador sem instalar o Java manualmente.

**Critérios de aceitação:**
- [ ] Sistema verifica se JDK 21 está disponível no `PATH` ou em `~/.hubsaude/jdk/`
- [ ] Se ausente, JDK é baixado automaticamente para a plataforma (Windows, Linux, macOS)
- [ ] JDK baixado é armazenado em `~/.hubsaude/jdk/` para reuso
- [ ] Download não é repetido se JDK já estiver provisionado
- [ ] Testes cobrem detecção de JDK presente e ausente nas três plataformas

### US-02.5 — Invocação do assinador.jar no modo local

**Como** usuário do Sistema Runner,
**quero** que o CLI invoque o assinador.jar diretamente via `java -jar` com os parâmetros fornecidos,
**para que** eu possa criar e validar assinaturas sem executar comandos Java manualmente.

**Critérios de aceitação:**
- [ ] CLI localiza o `java` disponível (provisionado ou do sistema)
- [ ] CLI constrói e executa o comando `java -jar assinador.jar` com parâmetros corretamente mapeados
- [ ] Saída do assinador.jar é capturada e repassada ao usuário
- [ ] Erros de execução (ex.: JDK ausente, jar não encontrado) são tratados com mensagens claras
- [ ] Testes de integração validam o fluxo CLI → assinador.jar

### US-02.6 — Exibição legível de resultados

**Como** usuário do Sistema Runner,
**quero** que o CLI apresente o resultado das operações de forma legível e estruturada,
**para que** eu compreenda facilmente o resultado da assinatura ou validação.

**Critérios de aceitação:**
- [ ] Resultado de criação de assinatura é formatado de forma legível
- [ ] Resultado de validação indica claramente se é válida ou inválida
- [ ] Erros são apresentados com mensagem descritiva e orientação para correção
- [ ] Saída é adequada para uso em terminal (não requer pós-processamento)

---

## Sprint 3 — Modo Servidor HTTP & Material Criptográfico

**Objetivo:** Reduzir latência eliminando o cold start da JVM. Suporte a token/smart card via PKCS#11. Gestão completa do ciclo de vida do processo servidor pelo CLI.

**Valor entregue:** Modo de execução com menor latência disponível. Suporte a dispositivo criptográfico real ou simulado.

### US-01.1 — Endpoints HTTP do assinador.jar

**Como** usuário do Sistema Runner,
**quero** que o assinador.jar exponha endpoints HTTP `/sign` e `/validate`,
**para que** o CLI possa invocá-lo via requisições HTTP no modo servidor.

**Critérios de aceitação:**
- [ ] `SignatureController` implementado com endpoints `POST /sign` e `POST /validate`
- [ ] Endpoints reutilizam a mesma lógica de validação e simulação do modo CLI
- [ ] Respostas HTTP seguem estrutura do contrato JSON definido no Sprint 2
- [ ] Testes de integração validam os endpoints

### US-01.3 — Iniciar assinador.jar no modo servidor

**Como** usuário do Sistema Runner,
**quero** que o CLI inicie o assinador.jar no modo servidor usando a porta padrão,
**para que** o assinador.jar fique disponível para requisições HTTP com menor latência.

**Critérios de aceitação:**
- [ ] CLI inicia o assinador.jar como processo em background na porta padrão
- [ ] PID e porta do processo são registrados em `~/.hubsaude/` para gestão posterior
- [ ] Feedback é exibido ao usuário confirmando que o servidor iniciou
- [ ] Porta pode ser personalizada via parâmetro `--port`

### US-01.5 — Detectar instância do assinador.jar em execução

**Como** usuário do Sistema Runner,
**quero** que o CLI detecte se já existe uma instância do assinador.jar em execução e a reutilize,
**para que** não sejam criadas instâncias duplicadas desnecessariamente.

**Critérios de aceitação:**
- [ ] CLI consulta `~/.hubsaude/` para verificar processo registrado
- [ ] Verificação de health check HTTP confirma que o processo está respondendo
- [ ] Se instância ativa é encontrada, CLI a reutiliza em vez de iniciar nova
- [ ] Se processo registrado não responde, é considerado inativo

### US-01.4 — Invocar assinador.jar via HTTP

**Como** usuário do Sistema Runner,
**quero** que o CLI envie requisições HTTP ao assinador.jar no modo servidor por padrão,
**para que** eu tenha menor latência nas operações, eliminando o overhead de cold start.

**Critérios de aceitação:**
- [ ] CLI envia requisições HTTP para os endpoints `/sign` e `/validate`
- [ ] Modo servidor é utilizado por padrão quando há instância em execução
- [ ] Fallback para modo local quando servidor não está disponível (ou via flag `--local`)
- [ ] Testes de integração validam o fluxo CLI → HTTP → assinador.jar

### US-01.6 — Interromper execução do assinador.jar

**Como** usuário do Sistema Runner,
**quero** interromper a execução do assinador.jar em uma porta específica ou na porta padrão,
**para que** eu tenha controle sobre os processos em execução no meu sistema.

**Critérios de aceitação:**
- [ ] Comando `assinatura stop` encerra o assinador.jar na porta padrão
- [ ] Parâmetro `--port` permite especificar a porta do processo a encerrar
- [ ] Feedback é exibido confirmando o encerramento
- [ ] Registro em `~/.hubsaude/` é atualizado após encerramento

### US-01.7 — Agendar interrupção por inatividade

**Como** usuário do Sistema Runner,
**quero** agendar a interrupção automática do assinador.jar após um período sem interação,
**para que** recursos do sistema sejam liberados automaticamente quando não estiverem em uso.

**Critérios de aceitação:**
- [ ] Parâmetro `--timeout <minutos>` define tempo máximo de inatividade
- [ ] Após o período sem requisições, assinador.jar é encerrado automaticamente
- [ ] Mecanismo de timeout é documentado no help do CLI

### US-01.2 — Integração com dispositivo criptográfico via PKCS#11

**Como** usuário do Sistema Runner,
**quero** que o assinador.jar suporte interação com dispositivo criptográfico (token/smart card) via PKCS#11,
**para que** eu possa utilizar material criptográfico real ou simulado nas operações de assinatura.

**Critérios de aceitação:**
- [ ] Integração com PKCS#11 via provider `SunPKCS11`
- [ ] Testes de integração utilizando SoftHSM2 (ou simulador equivalente)
- [ ] Comportamento adequado quando dispositivo não está disponível (mensagem clara)
- [ ] Documentação do setup necessário para uso com dispositivo criptográfico

### NOVO — Graceful shutdown do servidor

**Como** desenvolvedor do Sistema Runner,
**quero** que o assinador.jar trate o sinal de encerramento de forma controlada,
**para que** requisições em andamento sejam concluídas antes do processo ser finalizado.

**Critérios de aceitação:**
- [ ] Servidor captura sinal `SIGTERM` e inicia sequência de encerramento
- [ ] Requisições em curso são concluídas antes do processo ser finalizado
- [ ] Novas requisições são rejeitadas com resposta apropriada durante o shutdown
- [ ] Timeout máximo de graceful shutdown configurável

---

## Sprint 4 — CLI do Simulador & Gestão de Artefatos

**Objetivo:** CLI dedicado para o Simulador do HubSaúde com download automático, verificação de integridade e gestão completa do ciclo de vida. Extração de lógica compartilhada entre os dois CLIs.

**Valor entregue:** Sistema Runner completo. Todos os casos de uso funcionais e integrados.

### NOVO — Biblioteca compartilhada de gestão de processos

**Como** desenvolvedor do Sistema Runner,
**quero** que a lógica de gerenciamento de processos seja encapsulada em um pacote Go reutilizável,
**para que** os CLIs `assinatura` e `simulador` não dupliquem código de PID, health check e registro.

**Critérios de aceitação:**
- [ ] Pacote Go extraído com funções para: registrar processo, verificar health check HTTP, encerrar processo, ler/limpar registro em `~/.hubsaude/`
- [ ] Ambos os CLIs utilizam o pacote compartilhado
- [ ] Testes unitários cobrem o pacote isoladamente

### US-05.3 — Estrutura base do CLI "simulador" em Go

**Como** usuário do Sistema Runner,
**quero** um CLI dedicado para o Simulador com estrutura e organização próprias,
**para que** a gestão do Simulador tenha interface independente e clara.

**Critérios de aceitação:**
- [ ] Projeto CLI `simulador` segue a mesma estrutura do CLI `assinatura`
- [ ] Comandos `start`, `stop` e `status` definidos
- [ ] Pipeline CI/CD gera binários multiplataforma do CLI `simulador`
- [ ] Binários publicados no GitHub Releases junto com o CLI `assinatura`

### US-05.4 — Obter simulador.jar dinamicamente

**Como** usuário do Sistema Runner,
**quero** que o CLI baixe automaticamente a versão mais recente do simulador.jar do GitHub Releases,
**para que** eu sempre utilize a versão atualizada sem necessidade de download manual.

**Critérios de aceitação:**
- [ ] CLI consulta GitHub Releases para identificar a versão mais recente do simulador.jar
- [ ] Download automático quando simulador.jar não está disponível localmente
- [ ] Opção `--source <url>` permite indicar URL alternativa para download
- [ ] Versão já baixada não é baixada novamente (cache local em `~/.hubsaude/`)
- [ ] Verificação de integridade do download (checksum SHA-256)

### US-05.1 — Iniciar o Simulador via CLI

**Como** usuário do Sistema Runner,
**quero** iniciar o Simulador do HubSaúde através do CLI,
**para que** eu possa gerenciá-lo sem conhecer os comandos Java subjacentes.

**Critérios de aceitação:**
- [ ] Comando `simulador start` inicia o simulador.jar
- [ ] CLI verifica se as portas necessárias estão disponíveis antes de iniciar
- [ ] Se o simulador.jar não estiver disponível localmente, é baixado automaticamente
- [ ] Feedback exibido ao usuário sobre o status de inicialização

### US-05.2 — Parar e monitorar o Simulador

**Como** usuário do Sistema Runner,
**quero** parar o Simulador e consultar seu status atual,
**para que** eu tenha visibilidade e controle sobre o ciclo de vida do Simulador.

**Critérios de aceitação:**
- [ ] Comando `simulador stop` encerra o Simulador
- [ ] Comando `simulador status` exibe se o Simulador está em execução ou não
- [ ] Informações de processo (PID, porta) são registradas em `~/.hubsaude/`
- [ ] Encerramento limpo do processo com tratamento adequado de erros

### NOVO — Comando `update` nos dois CLIs

**Como** usuário do Sistema Runner,
**quero** atualizar o próprio binário do CLI sem precisar visitar o GitHub manualmente,
**para que** eu tenha sempre a versão mais recente disponível.

**Critérios de aceitação:**
- [ ] Comando `assinatura update` e `simulador update` verificam e baixam nova versão do binário
- [ ] Exibe changelog resumido da versão mais recente
- [ ] Não substitui o binário em execução (download para arquivo temporário, troca ao concluir)
- [ ] Verificação de integridade (checksum) antes de substituir o binário atual

---

## Sprint 5 — Segurança de Artefatos & Observabilidade

**Objetivo:** Artefatos distribuídos com checksums e assinatura criptográfica para garantia de integridade. Sistema observável com logs estruturados e métricas. Documentação completa para adoção.

**Valor entregue:** Sistema seguro, rastreável e documentado. Pronto para uso em produção.

### US-03.4 — Checksums SHA256 e assinatura de artefatos com Cosign

**Como** usuário do Sistema Runner,
**quero** que os binários distribuídos incluam checksums SHA256 e assinatura via Cosign,
**para que** eu possa verificar a integridade e autenticidade dos artefatos baixados.

**Critérios de aceitação:**
- [x] Cada release inclui arquivo `SHA256SUMS` para todos os binários
- [x] Artefatos assinados com Cosign (identidade OIDC + transparency log)
- [x] Cada artefato acompanhado de `.sig` e `.pem` conforme especificação
- [x] Processo de assinatura automatizado no pipeline CI/CD
- [x] Documentação de como verificar artefatos com `cosign verify-blob`

### US-05.4+ — Verificação de integridade no download

**Como** usuário do Sistema Runner,
**quero** que o CLI verifique a integridade do jar baixado antes de executá-lo,
**para que** eu tenha garantia de que o artefato não foi corrompido ou adulterado.

**Critérios de aceitação:**
- [ ] CLI valida checksum SHA-256 do jar após download
- [ ] CLI valida assinatura Cosign quando certificado disponível
- [ ] Download é rejeitado se verificação falhar, com mensagem clara
- [ ] Verificação pode ser ignorada com flag explícita `--skip-verify` (com aviso ao usuário)

### NOVO — Structured logging (JSON) nos dois CLIs

**Como** operador do Sistema Runner,
**quero** que os CLIs emitam logs em formato estruturado,
**para que** os logs possam ser integrados a pipelines de observabilidade (ex.: Loki, ELK).

**Critérios de aceitação:**
- [ ] Flag `--log-format json` ativa saída em JSON estruturado
- [ ] Campos mínimos: `timestamp`, `level`, `message`, `command`, `version`
- [ ] Saída padrão (sem flag) permanece legível para humanos no terminal
- [ ] Nível de log configurável via `--log-level` (debug, info, warn, error)

### NOVO — Endpoint `/metrics` no servidor assinador

**Como** operador do Sistema Runner,
**quero** que o servidor assinador exponha métricas em formato Prometheus,
**para que** eu possa monitorar saúde e performance em ambientes de produção.

**Critérios de aceitação:**
- [ ] Endpoint `GET /metrics` disponível no servidor HTTP
- [ ] Métricas mínimas: contagem de requisições por endpoint, latência (p50/p95/p99), uptime, erros por tipo
- [ ] Formato compatível com Prometheus (`text/plain; version=0.0.4`)
- [ ] Endpoint pode ser desabilitado via configuração

### NOVO — Testes E2E no pipeline CI

**Como** desenvolvedor do Sistema Runner,
**quero** que o pipeline CI execute testes de ponta-a-ponta após o build,
**para que** regressões no fluxo completo sejam detectadas automaticamente antes do merge.

**Critérios de aceitação:**
- [ ] Smoke tests que sobem o servidor assinador, disparam `sign` e `validate` e verificam saída
- [ ] Smoke tests que iniciam o simulador, verificam status e encerram
- [ ] Testes executados no CI para as três plataformas (ou via matrix strategy)
- [ ] Falha nos E2E bloqueia merge para branch principal

### NOVO — Documentação de onboarding e troubleshooting

**Como** novo usuário do Sistema Runner,
**quero** ter documentação clara sobre como instalar, configurar e usar o sistema,
**para que** eu possa começar a usar sem depender de suporte da equipe.

**Critérios de aceitação:**
- [ ] Guia de quickstart cobrindo instalação e primeiro uso em Windows, Linux e macOS
- [ ] Referência completa de comandos e flags para ambos os CLIs
- [ ] Guia de troubleshooting com erros comuns e soluções
- [ ] Instruções de verificação de integridade dos binários baixados
- [ ] Documentação de setup para uso com dispositivo PKCS#11 real

---

## Resumo dos Sprints

| Sprint | Foco principal | Risco removido |
|--------|---------------|----------------|
| 1 | Fundação & CI/CD | Nenhum artefato sem pipeline |
| 2 | Assinatura local ponta-a-ponta | Caso de uso principal não validado |
| 3 | Modo servidor HTTP & PKCS#11 | Latência e suporte a hardware criptográfico |
| 4 | CLI simulador & gestão de artefatos | Duplicação de código, ausência de auto-update |
| 5 | Segurança, observabilidade & docs | Artefatos sem garantia de integridade, sistema sem visibilidade |
