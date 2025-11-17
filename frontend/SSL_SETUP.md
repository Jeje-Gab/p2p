# Configuração SSL para Desenvolvimento

Este documento explica como aceitar o certificado SSL auto-assinado do backend durante o desenvolvimento.

## Por que isso é necessário?

O backend agora roda em HTTPS (porta 8443) com certificados auto-assinados. Navegadores modernos bloqueiam requisições para servidores com certificados não confiáveis por segurança.

## Passo a Passo: Aceitar o Certificado no Navegador

### Opção 1: Aceitar uma vez (Mais Rápido)

1. **Abra o backend diretamente no navegador:**
   ```
   https://localhost:8443/healthz
   ```

2. **Você verá um aviso de segurança:**
   - Chrome: "Your connection is not private"
   - Firefox: "Warning: Potential Security Risk Ahead"
   - Edge: "Your connection isn't private"

3. **Aceite o risco:**
   - **Chrome/Edge:** Clique em "Advanced" → "Proceed to localhost (unsafe)"
   - **Firefox:** Clique em "Advanced" → "Accept the Risk and Continue"

4. **Você deve ver:**
   ```json
   {"status":"healthy","timestamp":1234567890}
   ```

5. **Agora seu frontend funcionará!** O navegador lembrou que você confia neste certificado.

### Opção 2: Confiar no Certificado (Permanente)

Se você quer que o sistema confie permanentemente no certificado:

#### Windows

1. Abra PowerShell **como Administrador**
2. Execute:
   ```powershell
   Import-Certificate -FilePath ..\backend\certs\server.crt -CertStoreLocation Cert:\LocalMachine\Root
   ```

#### Linux

```bash
sudo cp ../backend/certs/server.crt /usr/local/share/ca-certificates/cs2-p2p.crt
sudo update-ca-certificates
```

#### macOS

```bash
sudo security add-trusted-cert -d -r trustRoot \
  -k /Library/Keychains/System.keychain ../backend/certs/server.crt
```

## Configuração do Next.js

O arquivo `src/lib/api.ts` já está configurado para aceitar certificados auto-assinados **apenas em desenvolvimento**:

```typescript
const httpsAgent = process.env.NODE_ENV !== 'production'
  ? new https.Agent({ rejectUnauthorized: false })
  : undefined;
```

⚠️ **IMPORTANTE:** Esta configuração é desabilitada automaticamente em produção.

## Testando

Depois de aceitar o certificado:

1. Reinicie o servidor de desenvolvimento:
   ```bash
   npm run dev
   ```

2. Acesse: http://localhost:3000

3. Tente fazer login - deve funcionar sem erros de SSL!

## Troubleshooting

### Erro: "TLS handshake error"

**Causa:** Certificado não aceito pelo navegador.

**Solução:** Siga a "Opção 1" acima.

### Erro: "NET::ERR_CERT_AUTHORITY_INVALID"

**Causa:** Normal com certificados auto-assinados.

**Solução:** Aceite o certificado conforme instruções acima.

### Erro persiste após aceitar

1. Limpe o cache do navegador (Ctrl+Shift+Delete)
2. Feche e reabra o navegador
3. Tente novamente

---

## Produção

Em produção, você **DEVE** usar certificados de uma CA confiável:

- **Let's Encrypt** (grátis)
- **DigiCert** (pago)
- **Cloud Provider** (AWS ACM, Google Cloud, Azure)

Nunca use certificados auto-assinados em produção!

---

**Nota:** Este setup é apenas para desenvolvimento local. Em produção, sempre use HTTPS com certificados válidos.
