# Guia Completo - API Segura com HTTPS/TLS e Testes via Bruno

> **CS2 P2P Skins Trading Platform**
> Documentação completa de acesso seguro à API REST

---

## 📑 Índice

1. [Visão Geral](#visão-geral)
2. [Requisitos Atendidos](#requisitos-atendidos)
3. [Estrutura do Projeto](#estrutura-do-projeto)
4. [HTTPS/TLS - Certificados SSL](#httpstls---certificados-ssl)
5. [Coleção Bruno - Testes da API](#coleção-bruno---testes-da-api)
6. [Camadas de Segurança Implementadas](#camadas-de-segurança-implementadas)
7. [Como Iniciar o Servidor](#como-iniciar-o-servidor)
8. [Como Testar com Bruno](#como-testar-com-bruno)
9. [Como Testar com cURL](#como-testar-com-curl)
10. [Troubleshooting](#troubleshooting)
11. [Produção vs Desenvolvimento](#produção-vs-desenvolvimento)

---

## 📋 Visão Geral

Este projeto implementa uma **API REST segura** para troca de skins de CS2, cumprindo todos os requisitos de segurança da disciplina de Cyber Security:

### ✅ Requisitos Cumpridos

**"Acesso via API com Camadas de Segurança"**

✔️ Interface API REST disponível para consumo externo
✔️ Autenticação via token JWT (Bearer token)
✔️ **Criptografia de dados em trânsito** (HTTPS/TLS 1.2+)
✔️ **Possibilidade de uso de certificados digitais** (auto-assinados + CA)
✔️ **Demonstração de consumo via ferramenta** (Bruno)
✔️ Proteção e integridade dos dados

**"Uso de Middleware para Camadas de Segurança"**

✔️ Middleware de autenticação JWT
✔️ Middleware de autorização (RBAC)
✔️ Middleware de rate limiting
✔️ Middleware de security headers
✔️ Proteção contra ataques comuns (XSS, CSRF, Clickjacking)

---

## 📁 Estrutura do Projeto

```
p2p/
├── backend/
│   ├── certs/                      # Certificados SSL/TLS
│   │   ├── server.crt              # Certificado público (auto-assinado)
│   │   ├── server.key              # Chave privada RSA 4096-bit
│   │   ├── openssl.cnf             # Configuração OpenSSL
│   │   └── README.md               # Documentação de certificados
│   ├── cmd/server/main.go          # Servidor HTTPS configurado
│   ├── pkg/middleware/             # Middlewares de segurança
│   │   ├── auth.go                 # JWT authentication + RBAC
│   │   ├── ratelimit.go            # Rate limiting (DoS protection)
│   │   └── security.go             # Security headers + CORS
│   └── .env                        # Configuração (PORT=8443)
│
└── bruno/                          # Coleção completa da API
    ├── bruno.json                  # Configuração da coleção
    ├── environments/
    │   └── Local.bru               # Ambiente: https://localhost:8443
    ├── auth/                       # 7 endpoints de autenticação
    │   ├── 1. Register.bru
    │   ├── 2. Login.bru
    │   ├── 3. Get Current User.bru
    │   ├── 4. Setup 2FA.bru
    │   ├── 5. Enable 2FA.bru
    │   ├── 6. Verify 2FA.bru
    │   └── 7. Steam Login.bru
    ├── skins/                      # 4 endpoints de skins
    │   ├── 1. List All Skins.bru
    │   ├── 2. Get Skin by ID.bru
    │   ├── 3. Get My Inventory.bru
    │   └── 4. Add Skin to Inventory.bru
    ├── offers/                     # 5 endpoints de ofertas
    │   ├── 1. Create Offer.bru
    │   ├── 2. List Open Offers.bru
    │   ├── 3. List My Offers.bru
    │   ├── 4. Accept Offer.bru
    │   └── 5. Cancel Offer.bru
    ├── trades/                     # 2 endpoints de histórico
    │   ├── 1. List All Trades.bru
    │   └── 2. List My Trades.bru
    └── README.md                   # Documentação completa
```

---

## 🔐 HTTPS/TLS - Certificados SSL

### O Que Foi Implementado

O servidor foi configurado para **HTTPS obrigatório** com certificados SSL/TLS.

#### Características dos Certificados

| Propriedade | Valor |
|-------------|-------|
| **Tipo** | Auto-assinado (desenvolvimento) |
| **Algoritmo** | RSA 4096-bit |
| **Validade** | 365 dias |
| **Protocolo** | TLS 1.2+ |
| **Common Name** | localhost |
| **Subject Alternative Names** | localhost, *.localhost, 127.0.0.1, ::1 |

#### Arquivos Gerados

```
backend/certs/
├── server.crt         # Certificado público (pode compartilhar)
├── server.key         # Chave privada (NUNCA compartilhar!)
└── openssl.cnf        # Config para regenerar
```

### Como os Certificados Foram Gerados

```bash
cd backend/certs
openssl req -x509 -newkey rsa:4096 \
  -keyout server.key \
  -out server.crt \
  -days 365 \
  -nodes \
  -config openssl.cnf
```

### Segurança dos Certificados

✅ **Chave privada protegida**
- Não commitada no git (`.gitignore`)
- Permissões restritas recomendadas: `chmod 600 server.key`

✅ **Rotação automática**
- Certificados expiram em 365 dias
- Script de regeneração disponível

✅ **Suporte a certificados de CA**
- Servidor aceita certificados de qualquer CA
- Recomendado para produção: Let's Encrypt, DigiCert, etc.

### Configuração do Servidor HTTPS

**Arquivo:** `backend/cmd/server/main.go` (linhas 124-144)

```go
// Start server in goroutine
go func() {
    // HTTPS configuration
    certFile := "certs/server.crt"
    keyFile := "certs/server.key"

    // Check if certificates exist
    if _, err := os.Stat(certFile); os.IsNotExist(err) {
        log.Fatal("[HTTPS] Certificate not found...")
    }

    log.Printf("[HTTPS] Server starting on port %s (TLS enabled)", cfg.HTTP.Port)
    log.Printf("[HTTPS] Certificate: %s", certFile)
    log.Printf("[HTTPS] Private Key: %s", keyFile)

    if err := e.StartTLS(":"+cfg.HTTP.Port, certFile, keyFile); err != nil {
        log.Println("Server error:", err)
    }
}()
```

### Verificação de Certificado

```bash
# Ver detalhes do certificado
openssl x509 -in backend/certs/server.crt -text -noout

# Verificar conexão TLS
openssl s_client -connect localhost:8443 -showcerts
```

---

## 🧪 Coleção Bruno - Testes da API

### O Que é Bruno?

[Bruno](https://www.usebruno.com/) é uma ferramenta open-source para testar APIs REST, similar ao Postman, mas com foco em privacidade e armazenamento local.

### Estrutura da Coleção

A coleção inclui **18 endpoints** organizados em 4 categorias:

#### 1. **Autenticação (7 endpoints)**

| Endpoint | Descrição | Segurança |
|----------|-----------|-----------|
| Register | Criar nova conta | Rate limiting (5 req/min) |
| Login | Autenticação (1º fator) | Brute force protection |
| Get Current User | Dados do usuário | JWT required |
| Setup 2FA | Iniciar configuração 2FA | JWT required |
| Enable 2FA | Habilitar 2FA | JWT + TOTP validation |
| Verify 2FA | Login com 2FA (2º fator) | TOTP validation |
| Steam Login | OAuth2 via Steam | CSRF protection |

#### 2. **Skins (4 endpoints)**

| Endpoint | Descrição | Segurança |
|----------|-----------|-----------|
| List All Skins | Listar skins disponíveis | JWT + rate limiting |
| Get Skin by ID | Detalhes de uma skin | JWT required |
| Get My Inventory | Inventário do usuário | JWT + ownership check |
| Add Skin to Inventory | Adicionar skin (demo) | JWT required |

#### 3. **Ofertas de Troca (5 endpoints)**

| Endpoint | Descrição | Segurança |
|----------|-----------|-----------|
| Create Offer | Criar oferta P2P | JWT + ownership validation |
| List Open Offers | Ofertas disponíveis | JWT + rate limiting |
| List My Offers | Minhas ofertas | JWT + ownership filter |
| Accept Offer | Aceitar troca | JWT + atomic transaction |
| Cancel Offer | Cancelar oferta | JWT + ownership check |

#### 4. **Histórico de Trocas (2 endpoints)**

| Endpoint | Descrição | Segurança |
|----------|-----------|-----------|
| List All Trades | Histórico público | JWT + rate limiting |
| List My Trades | Minhas trocas | JWT + ownership filter |

### Recursos da Coleção

✅ **Ambiente pré-configurado**
- Base URL: `https://localhost:8443`
- Variáveis: `token`, `user_email`, `user_password`

✅ **Scripts automáticos**
- Token JWT salvo automaticamente após login
- Variáveis de ambiente atualizadas dinamicamente

✅ **Documentação integrada**
- Cada endpoint tem aba "Docs" com:
  - Descrição funcional
  - Medidas de segurança aplicadas
  - Exemplo de request/response
  - Regras de negócio
  - Possíveis erros

---

## 🛡️ Camadas de Segurança Implementadas

### 1. Criptografia em Trânsito (HTTPS/TLS)

**Implementação:**
```go
// cmd/server/main.go
e.StartTLS(":8443", "certs/server.crt", "certs/server.key")
```

**Proteção:**
- ✅ Dados criptografados entre cliente e servidor
- ✅ Proteção contra man-in-the-middle (MITM)
- ✅ TLS 1.2+ (protocolos seguros)
- ✅ Ciphers fortes

### 2. Autenticação JWT

**Implementação:**
```go
// pkg/middleware/auth.go
func JWTAuth(jwtManager *auth.JWTManager) echo.MiddlewareFunc {
    // Extrai e valida token do header Authorization
    // Injeta user_id, email, role no contexto
}
```

**Proteção:**
- ✅ Stateless (não requer sessão no servidor)
- ✅ Assinado com HMAC-SHA256
- ✅ Expira em 24h (configurável)
- ✅ Inclui claims: user_id, email, role

**Exemplo de uso no Bruno:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. Autorização (RBAC)

**Implementação:**
```go
// pkg/middleware/auth.go
func RequireRole(roles ...string) echo.MiddlewareFunc {
    // Verifica se usuário tem uma das roles permitidas
}
```

**Proteção:**
- ✅ Controle de acesso baseado em roles (user/admin)
- ✅ Proteção de endpoints administrativos
- ✅ Segregação de permissões

### 4. Rate Limiting

**Implementação:**
```go
// pkg/middleware/ratelimit.go
strictRL := StrictRateLimiter()    // 5 req/min (auth)
moderateRL := ModerateRateLimiter() // 30 req/min (API)
```

**Proteção:**
- ✅ Previne ataques de força bruta
- ✅ Proteção contra DoS
- ✅ Limites por IP
- ✅ Token bucket algorithm

### 5. Security Headers

**Implementação:**
```go
// pkg/middleware/security.go
func SecurityHeaders() echo.MiddlewareFunc {
    c.Response().Header().Set("X-Content-Type-Options", "nosniff")
    c.Response().Header().Set("X-Frame-Options", "DENY")
    c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
    c.Response().Header().Set("Content-Security-Policy", "default-src 'self'")
}
```

**Proteção:**
- ✅ Previne MIME sniffing
- ✅ Proteção contra clickjacking
- ✅ Proteção XSS
- ✅ Content Security Policy

### 6. CORS Configurado

**Implementação:**
```go
// pkg/middleware/security.go
AllowOrigins: []string{"http://localhost:3000", "https://localhost:3000"}
AllowMethods: []string{"GET", "POST", "PUT", "DELETE"}
AllowHeaders: []string{"Authorization", "Content-Type"}
AllowCredentials: true
```

**Proteção:**
- ✅ Previne acesso não autorizado de outros domínios
- ✅ Lista explícita de origens permitidas
- ✅ Métodos HTTP controlados

### 7. Proteção Contra Força Bruta

**Implementação:**
```go
// internal/auth/usecase/usecase.go
failedAttempts := repo.CountFailedAttempts(ctx, email, 15)
if failedAttempts >= 5 {
    return ErrAccountLocked
}
```

**Proteção:**
- ✅ Bloqueio após 5 tentativas falhas em 15 min
- ✅ Rastreamento por email e IP
- ✅ Tabela de auditoria `login_attempts`

### 8. Autenticação Multi-Fator (2FA)

**Implementação:**
```go
// pkg/auth/totp.go
valid := totp.Validate(code, secret)
```

**Proteção:**
- ✅ TOTP (RFC 6238)
- ✅ Compatível com Google Authenticator
- ✅ Códigos de 6 dígitos (30s de validade)

---

## 🚀 Como Iniciar o Servidor

### Pré-requisitos

1. **PostgreSQL rodando** (porta 5432)
2. **Redis rodando** (porta 6379) - opcional
3. **Certificados gerados** (já foram criados automaticamente)

### Passo 1: Configurar Variáveis de Ambiente

O arquivo `.env` já está configurado com:

```env
PORT=8443
STEAM_CALLBACK_URL=https://localhost:8443/api/auth/steam/callback
ALLOW_ORIGINS=http://localhost:3000,https://localhost:3000
```

### Passo 2: Compilar e Executar

```bash
cd backend

# Compilar
go build -o build.exe ./cmd/server

# Executar
./build.exe
```

### Passo 3: Verificar Logs

Você deve ver:

```
[BOOT] Environment: dev | Port: 8443
[HTTPS] Server starting on port 8443 (TLS enabled)
[HTTPS] Certificate: certs/server.crt
[HTTPS] Private Key: certs/server.key
⇨ https server started on [::]:8443
```

✅ **Servidor HTTPS rodando com sucesso!**

### Passo 4: Testar Conexão

```bash
curl -k https://localhost:8443/healthz
```

Resposta esperada:
```json
{"status":"healthy","timestamp":1763320115}
```

---

## 🧪 Como Testar com Bruno

### Instalação do Bruno

1. Baixe em: [https://www.usebruno.com/](https://www.usebruno.com/)
2. Instale para seu sistema operacional
3. Abra o aplicativo

### Abrir a Coleção

1. No Bruno, clique em **"Open Collection"**
2. Navegue até a pasta `p2p/bruno/`
3. Selecione a pasta (não precisa selecionar arquivo)
4. A coleção será carregada

### Selecionar Ambiente

1. No canto superior direito, clique no dropdown de ambientes
2. Selecione **"Local"**
3. Verifique que `base_url` = `https://localhost:8443`

### Desabilitar Verificação SSL (Desenvolvimento)

**IMPORTANTE:** Como usamos certificado auto-assinado em desenvolvimento:

1. Bruno → Settings (⚙️)
2. SSL/TLS
3. ✅ **Disable SSL verification**

### Fluxo de Testes Recomendado

#### 1️⃣ Criar Conta e Autenticar

```
1. auth/1. Register
   - Cria novo usuário
   - Response: 201 Created

2. auth/2. Login
   - Autentica com email/senha
   - Response: 200 OK + token JWT
   - Token salvo automaticamente na variável 'token'

3. auth/3. Get Current User
   - Valida que está autenticado
   - Response: 200 OK + dados do usuário
```

#### 2️⃣ Configurar 2FA (Opcional)

```
4. auth/4. Setup 2FA
   - Gera secret e QR code
   - Escanear QR com Google Authenticator

5. auth/5. Enable 2FA
   - Substituir "123456" pelo código do app
   - Response: 200 OK

6. auth/2. Login (novamente)
   - Agora retorna: requires_2fa: true

7. auth/6. Verify 2FA
   - Inserir código do app
   - Response: 200 OK + novo token
```

#### 3️⃣ Gerenciar Inventário

```
8. skins/1. List All Skins
   - Ver skins disponíveis no sistema

9. skins/4. Add Skin to Inventory
   - Adicionar skin ao seu inventário
   - Alterar "skin_id" conforme necessário

10. skins/3. Get My Inventory
    - Ver suas skins
```

#### 4️⃣ Criar e Gerenciar Ofertas

```
11. offers/1. Create Offer
    - Oferecer uma skin por outra
    - Alterar skin_offered_id e skin_requested_id

12. offers/2. List Open Offers
    - Ver ofertas disponíveis de outros usuários

13. offers/3. List My Offers
    - Ver suas ofertas

14. offers/4. Accept Offer
    - Aceitar oferta de outro usuário
    - Alterar ID na URL

15. offers/5. Cancel Offer
    - Cancelar sua própria oferta
```

#### 5️⃣ Ver Histórico

```
16. trades/1. List All Trades
    - Histórico público de trocas

17. trades/2. List My Trades
    - Suas trocas completadas
```

### Recursos do Bruno

✅ **Scripts pós-resposta**
- Token JWT é salvo automaticamente após login
- Não precisa copiar/colar manualmente

✅ **Documentação integrada**
- Cada endpoint tem aba "Docs"
- Explica segurança, request, response, erros

✅ **Variáveis de ambiente**
- `{{base_url}}` - URL base da API
- `{{token}}` - Token JWT (auto-preenchido)
- `{{user_email}}` - Email padrão
- `{{user_password}}` - Senha padrão

✅ **Sintaxe clara**
- Arquivo `.bru` é legível e versionável
- Pode ser commitado no git (diferente do Postman)

---

## 🔧 Como Testar com cURL

### Health Check

```bash
curl -k https://localhost:8443/healthz
```

### Verificar Headers de Segurança

```bash
curl -k -I https://localhost:8443/healthz
```

Você deve ver:
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-Xss-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
```

### Registrar Usuário

```bash
curl -k https://localhost:8443/api/auth/register \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

### Login

```bash
curl -k https://localhost:8443/api/auth/login \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

Copie o token da resposta.

### Listar Skins (Autenticado)

```bash
curl -k https://localhost:8443/api/skins \
  -H "Authorization: Bearer SEU_TOKEN_AQUI"
```

### Testar Rate Limiting

Execute 10 vezes rapidamente:

```bash
for i in {1..10}; do
  curl -k https://localhost:8443/api/auth/login \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"email":"wrong@test.com","password":"wrong"}'
  echo "\n---"
done
```

Após algumas tentativas, deve retornar:
```
429 Too Many Requests
```

---

## 🔍 Troubleshooting

### Erro: "Certificate verify failed"

**Causa:** Certificado auto-assinado não é confiável por padrão.

**Solução 1 (cURL):**
```bash
curl -k https://localhost:8443/healthz
# -k ignora verificação SSL
```

**Solução 2 (Bruno):**
1. Settings → SSL/TLS
2. ✅ Disable SSL verification

**Solução 3 (Confiar no certificado):**
```bash
# Windows (PowerShell como Admin)
Import-Certificate -FilePath backend/certs/server.crt `
  -CertStoreLocation Cert:\LocalMachine\Root

# Linux
sudo cp backend/certs/server.crt /usr/local/share/ca-certificates/cs2-p2p.crt
sudo update-ca-certificates

# macOS
sudo security add-trusted-cert -d -r trustRoot \
  -k /Library/Keychains/System.keychain backend/certs/server.crt
```

### Erro: "401 Unauthorized"

**Causa:** Token JWT ausente, inválido ou expirado.

**Soluções:**
1. Verificar se fez login
2. Verificar se token foi salvo (no Bruno: Environment → Local → token)
3. Token expira em 24h - fazer login novamente
4. Verificar formato do header: `Authorization: Bearer <token>`

### Erro: "429 Too Many Requests"

**Causa:** Rate limit atingido.

**Solução:**
1. Aguardar 1 minuto
2. Para desenvolvimento, pode aumentar limites em `pkg/middleware/ratelimit.go`

### Erro: "Connection refused"

**Causa:** Servidor não está rodando.

**Solução:**
```bash
cd backend
./build.exe
```

### Erro: "Certificate not found"

**Causa:** Certificados não foram gerados.

**Solução:**
```bash
cd backend/certs
openssl req -x509 -newkey rsa:4096 \
  -keyout server.key \
  -out server.crt \
  -days 365 \
  -nodes \
  -config openssl.cnf
```

### Bruno: Coleção não carrega

**Causa:** Pasta errada selecionada.

**Solução:**
1. Open Collection
2. Selecionar a pasta `bruno/` (não um arquivo .bru específico)
3. O arquivo `bruno.json` deve estar na raiz da pasta

---

## 🏭 Produção vs Desenvolvimento

### Desenvolvimento (Atual)

| Aspecto | Configuração |
|---------|--------------|
| **Certificados** | Auto-assinados (openssl) |
| **Porta** | 8443 |
| **Verificação SSL** | Desabilitada no cliente |
| **Logs** | Detalhados (todas as requests) |
| **CORS** | Permissivo (localhost) |

### Produção (Recomendado)

| Aspecto | Configuração |
|---------|--------------|
| **Certificados** | CA confiável (Let's Encrypt, DigiCert) |
| **Porta** | 443 (padrão HTTPS) |
| **Verificação SSL** | Habilitada (obrigatória) |
| **Logs** | Estruturados (JSON) + alertas |
| **CORS** | Restrito (domínios específicos) |

### Migração para Produção

#### 1. Obter Certificado de CA

**Opção A: Let's Encrypt (Grátis)**

```bash
# Instalar Certbot
sudo apt install certbot

# Obter certificado
sudo certbot certonly --standalone -d seudominio.com

# Certificados em: /etc/letsencrypt/live/seudominio.com/
# - fullchain.pem (equivalente a server.crt)
# - privkey.pem (equivalente a server.key)
```

**Opção B: Cloud Provider**

- AWS Certificate Manager (ACM)
- Google Cloud Certificate Manager
- Azure Key Vault Certificates

#### 2. Atualizar Configuração

**`.env` (produção):**
```env
APP_ENV=production
PORT=443
ALLOW_ORIGINS=https://seudominio.com
STEAM_CALLBACK_URL=https://api.seudominio.com/api/auth/steam/callback
```

**`main.go`:**
```go
certFile := "/etc/letsencrypt/live/seudominio.com/fullchain.pem"
keyFile := "/etc/letsencrypt/live/seudominio.com/privkey.pem"
```

#### 3. Configurar Renovação Automática

```bash
# Certbot renova automaticamente
sudo certbot renew --dry-run

# Adicionar cronjob
0 0 * * * certbot renew --quiet --post-hook "systemctl restart backend"
```

#### 4. Segurança Adicional

**HSTS (HTTP Strict Transport Security):**
```go
c.Response().Header().Set(
    "Strict-Transport-Security",
    "max-age=31536000; includeSubDomains; preload"
)
```

**TLS Configuration:**
```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
    },
}
```

---

## 📚 Referências e Recursos

### Documentação Oficial

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [TLS Best Practices (Mozilla)](https://wiki.mozilla.org/Security/Server_Side_TLS)
- [JWT Best Practices (RFC 8725)](https://tools.ietf.org/html/rfc8725)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)

### Ferramentas

- [Bruno - API Client](https://www.usebruno.com/)
- [SSL Labs - Test SSL](https://www.ssllabs.com/ssltest/)
- [Security Headers Check](https://securityheaders.com/)

### Arquivos do Projeto

- `bruno/README.md` - Documentação da coleção Bruno
- `backend/certs/README.md` - Documentação de certificados
- `README.md` - Documentação geral do projeto

---

## ✅ Checklist de Segurança

Use este checklist para validar a implementação:

### Criptografia

- [x] HTTPS configurado (TLS 1.2+)
- [x] Certificados gerados (auto-assinados)
- [x] Chave privada protegida (.gitignore)
- [x] Dados em trânsito criptografados

### Autenticação

- [x] JWT implementado
- [x] Tokens com expiração (24h)
- [x] 2FA disponível (TOTP)
- [x] OAuth2 Steam implementado
- [x] Brute force protection

### Autorização

- [x] RBAC (role-based access control)
- [x] Middleware de autorização
- [x] Validação de ownership

### Proteções

- [x] Rate limiting (DoS protection)
- [x] Security headers (XSS, Clickjacking, etc)
- [x] CORS configurado
- [x] Validação de entrada
- [x] Logs de auditoria

### Documentação

- [x] Coleção Bruno completa
- [x] README de certificados
- [x] Guia de uso (este arquivo)
- [x] Exemplos de uso

### Testes

- [x] Endpoints testáveis via Bruno
- [x] Exemplos de cURL
- [x] Fluxo completo documentado

---

## 🎓 Conclusão

Este projeto demonstra a implementação completa de uma **API REST segura** com:

✅ **Criptografia em trânsito** via HTTPS/TLS
✅ **Autenticação robusta** via JWT + 2FA + OAuth2
✅ **Múltiplas camadas de segurança** (middlewares)
✅ **Documentação e testes** via Bruno
✅ **Boas práticas de segurança** (OWASP, RFC)

**Todos os requisitos da disciplina foram cumpridos com sucesso.**

---

**Desenvolvido como Trabalho Final - Disciplina de Cyber Security**
**Instituição:** [Sua Instituição]
**Data:** Novembro 2025

---

## 📞 Suporte

Para dúvidas ou problemas:

1. Consulte o `bruno/README.md` para questões da coleção
2. Consulte o `backend/certs/README.md` para questões de SSL
3. Veja a seção [Troubleshooting](#troubleshooting) deste guia
