# Troubleshooting — Sistema Runner

Erros comuns do CLI `assinatura` e do servidor `assinador.jar`, com soluções.

## Instalação e integridade

### `SHA-256 inválido` ao baixar o JRE

O download do JRE 21 (Adoptium) foi corrompido. O CLI rejeita e aborta.
**Solução:** rode o comando de novo (rede instável) ou limpe o cache:

```bash
rm -rf ~/.hubsaude/jdk          # Linux/macOS
Remove-Item -Recurse ~/.hubsaude/jdk   # Windows
```

### `verificação de integridade do jar falhou`

Há arquivos `assinador.jar.sig` e `assinador.jar.pem` ao lado do jar, mas a
verificação Cosign falhou (jar adulterado, sig/cert errados, ou identidade
divergente). **O start é bloqueado (fail-closed).**

- Verifique manualmente com `cosign verify-blob` (ver
  [onboarding.md](onboarding.md#verificar-integridade-dos-binários-recomendado)).
- Se confia na origem e quer pular: `--skip-verify` (emite aviso).

### `cosign não encontrado no PATH`

Sidecars `.sig`/`.pem` presentes mas o Cosign não está instalado.
**Solução:** instale o Cosign (<https://docs.sigstore.dev/cosign/installation/>)
ou remova os sidecars / use `--skip-verify`.

## Java / jar

### `java não encontrado (instale JDK 21 ou defina JAVA_HOME)`

O CLI tenta `JAVA_HOME`, `PATH` e `~/.hubsaude/jdk/`. Se todos falharem,
provisiona o JRE 21 automaticamente. Se o provisionamento falhar (sem rede):

```bash
export JAVA_HOME=/caminho/para/jdk-21   # Linux/macOS
$env:JAVA_HOME = "C:\caminho\jdk-21"    # Windows
```

### `assinador.jar não encontrado`

Forneça o jar por um destes meios:

```bash
assinatura server start --jar /caminho/assinador.jar
# ou
export ASSINATURA_JAR=/caminho/assinador.jar
# ou coloque em ~/.hubsaude/assinador.jar
```

### `assinador-java não respondeu em 90s`

O processo subiu mas a porta HTTP não respondeu. Causas comuns:

- Porta ocupada → use `--port` diferente.
- JRE corrompido → limpe `~/.hubsaude/jdk` e tente de novo.
- Veja os logs do jar (stderr do processo) para stack traces.

## Servidor / rede

### `instância registrada usa porta X, não Y`

O `server stop --port Y` não bate com a porta registrada no PID file
(`~/.hubsaude/assinador-java.pid`). Use a porta correta ou `server stop` sem
`--port` (usa o registro).

### `nenhuma instância registrada`

Não há PID file. O servidor não foi iniciado por este CLI, ou já foi encerrado.
Verifique processos Java manualmente se necessário.

### `/metrics` retorna 404

O endpoint está desabilitado. Habilite no `application.properties`:

```properties
assinador.metrics.enabled=true
```

## Assinatura

### `PIN` solicitado interativamente em ambiente sem TTY

Em CI/scripts, forneça o PIN sem prompt:

```bash
export ASSINATURA_PKCS11_PIN=1234
# ou
assinatura sign --pkcs11-pin 1234 ...
```

### Erro de validação de flags (`timestamp`, `policy-id`, etc.)

A validação é rigorosa (`internal/sign/validate.go`):

- `--timestamp` deve estar em `[1751328000, 4102444800]`.
- `--policy-id` no formato `<baseURI>|<major.minor.patch>`.
- `--strategy` ∈ {`iat`, `tsa`}.
- `--crypto-type` ∈ {`smartcard`, `token`}.
- `--pkcs11-token-label` ≤ 32 chars UTF-8.

### `bad request` do servidor ao assinar

O payload chegou mas o serviço rejeitou. Rode com `--log-level debug` e veja a
mensagem do `ErrorResponse` retornada pelo servidor.

## PKCS#11 (dispositivo real)

### `assinador.pkcs11.library` vazio → modo simulado

Sem a lib configurada, o servidor assina com chave RSA **em memória** (apenas
para testes). Para hardware real, configure o caminho da lib PKCS#11 do
fabricante (`opensc-pkcs11.so`, `eToken.dll`, etc.).

### Slot/token não encontrado

Confirme o dispositivo conectado e o slot correto:

```bash
pkcs11-tool --list-slots
```

Passe `--pkcs11-slot` e/ou `--pkcs11-token-label` conforme o slot listado.

## Logs estruturados

Para integrar a Loki/ELK, ative JSON:

```bash
assinatura --log-format json --log-level info server start
```

Campos emitidos: `time`, `level`, `msg`, `command`, `version` (mais atributos
contextuais por evento).
