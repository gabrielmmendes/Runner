# Contrato JSON — CLI ↔ assinador-java

Comunicação entre `assinatura` CLI (Go) e `assinador-java` (Spring Boot) via HTTP POST com `Content-Type: application/json`.

---

## POST `/api/sign`

### Request

```jsonc
{
  "bundle": { /* FHIR Bundle JSON completo */ },
  "provenance": { /* FHIR Provenance JSON completo */ },
  "certChain": [
    "MIIBiDCC...",   // base64 DER — certificado folha
    "MIIBfzCC..."    // base64 DER — certificado raiz (mínimo 2)
  ],
  "timestamp": 1751500000,           // Unix UTC seconds [1751328000, 4102444800]
  "strategy": "iat",                 // "iat" | "tsa"
  "policyId": "https://example/policy|1.0.0",  // <baseURI>|<semver>
  "cryptoMaterial": {
    "type": "smartcard",             // "smartcard" | "token"
    "pin": "1234",                   // PIN PKCS#11
    "identifier": "key-alias",      // alias da chave privada
    "slotId": 0,                     // (opcional) slot PKCS#11 ≥ 0
    "tokenLabel": "MyToken"          // (opcional) ≤ 32 chars UTF-8
  },
  "operationalConfig": {
    "verification": {
      "tsaUrl": "https://tsa.example",  // obrigatório se strategy=tsa
      "checkRevocation": false
    },
    "trustStore": {
      "type": "JKS",
      "path": "",
      "password": ""
    },
    "temporalPolicy": {
      "allowedClockSkewSeconds": 300
    },
    "security": {
      "requireSecureChannel": false
    },
    "middlewareCrypto": {
      "library": "",
      "slotDescription": ""
    }
  }
}
```

**Campos legados** (aceitos pelo Java, não usados pelo CLI):
- `data` (String) — payload alternativo ao bundle
- `alias` (String) — alias alternativo ao `cryptoMaterial.identifier`
- `pin` (String, raiz) — PIN alternativo ao `cryptoMaterial.pin`

### Response — Sucesso (HTTP 200)

```json
{
  "success": true,
  "signature": "base64EncodedSignature...",
  "algorithm": "SHA256withRSA"
}
```

### Response — Erro (HTTP 400)

```json
{
  "success": false,
  "error": "mensagem de erro",
  "timestamp": "2026-05-26T10:30:00"
}
```

---

## POST `/api/validate`

### Request

```json
{
  "data": "conteúdo do documento assinado (string)",
  "signature": "base64EncodedSignature..."
}
```

### Response — Sucesso (HTTP 200)

```json
{
  "valid": true,
  "message": "Validação simulada"
}
```

### Response — Erro (HTTP 400)

```json
{
  "success": false,
  "error": "mensagem de erro",
  "timestamp": "2026-05-26T10:30:00"
}
```

---

## Notas

| Aspecto | Detalhe |
|---------|---------|
| Content-Type | `application/json` (request e response) |
| Accept | `application/json` |
| Timeout padrão | 30 segundos (configurável via `--timeout`) |
| Autenticação | Nenhuma (serviço local) |
| Porta padrão | 8080 (`--port` ou env `ASSINATURA_SERVICE_URL`) |
| certChain | PEM bundle ou JSON array de strings base64 — CLI normaliza para array |
| operationalConfig | Chaves permitidas: `verification`, `trustStore`, `temporalPolicy`, `security`, `middlewareCrypto` |
