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

## 4. Como Usar (Em Desenvolvimento)
*(As instruções de instalação e os comandos de uso serão adicionados aqui conforme o desenvolvimento do projeto avançar.)*

## 5. Contexto Acadêmico
|||
|:--|:--|
|**Instituição:**| Universidade Federal de Goiás (UFG) - Instituto de Informática | 
|**Curso:** | Bacharelado em Engenharia de Software|
|**Disciplina:** | Implementação e Integração de Software (INF0466) |
|**Professor:**|Fabio Nogueira de Lucena|
|**Semestre Letivo:**|2026/1|
|**Estudantes:**| Gabriel Matos Mendes,  Eduardo Divino Miranda Pereira|
