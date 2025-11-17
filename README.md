# CS2 P2P Skins Trading Platform

> **Trabalho Final - Disciplina de Cyber Security**
> Plataforma P2P para troca de skins de Counter-Strike 2 com foco em segurança

## 📋 Visão Geral

Este projeto demonstra a implementação prática de conceitos avançados de segurança cibernética em uma aplicação web completa. A plataforma permite que usuários troquem skins do CS2 de forma peer-to-peer, com múltiplas camadas de proteção e autenticação.

## 🎯 Objetivos de Cyber Security Demonstrados

### 1. **Autenticação Multi-Fator (2FA)**
- Implementação de TOTP (Time-based One-Time Password)
- Integração com Google Authenticator
- Segunda camada de segurança além de email/senha

### 2. **Autenticação OAuth2**
- Login via Steam usando OpenID/OAuth2
- Proteção contra CSRF com state parameter
- Armazenamento seguro de tokens

### 3. **Criptografia e Hashing**
- Senhas armazenadas com **bcrypt** (cost factor 12)
- NUNCA armazena senhas em texto plano
- JWT para sessões stateless

### 4. **Proteção Contra Força Bruta**
- Rastreamento de tentativas de login falhas
- Bloqueio temporário após N tentativas (lockout)
- Monitoramento por IP e por usuário
- Tabela de auditoria (`login_attempts`)

### 5. **Middlewares de Segurança**
- **JWT Authentication**: Validação de tokens
- **Role-Based Access Control (RBAC)**: Controle de permissões (user/admin)
- **Rate Limiting**: Proteção contra DoS
- **Security Headers**: X-Content-Type-Options, X-Frame-Options, CSP
- **CORS**: Configuração adequada de origens permitidas

### 6. **HTTPS/TLS (Criptografia em Trânsito)**
- Certificados SSL/TLS auto-assinados para desenvolvimento
- Suporte a TLS 1.2+
- Criptografia de dados em trânsito
- Possibilidade de uso de certificados de CA (Let's Encrypt, etc)

### 7. **Arquitetura Segura**
- Clean Architecture (separação de camadas)
- Validação de entrada em todos os endpoints
- Tratamento adequado de erros sem expor detalhes internos
- Logs estruturados para auditoria

## 🛠️ Tecnologias Utilizadas

- **Backend**: Go 1.24 + Echo Framework
- **Banco de Dados**: PostgreSQL 16
- **Cache**: Redis 7
- **Autenticação**:
  - JWT (golang-jwt/jwt/v5)
  - TOTP (pquerna/otp)
  - OAuth2 (golang.org/x/oauth2)
- **Segurança**:
  - bcrypt (golang.org/x/crypto)
  - Rate limiting (golang.org/x/time/rate)
- **Containerização**: Docker + Docker Compose

## 📁 Estrutura do Projeto

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point da aplicação
├── internal/
│   ├── entity/                     # Entidades do domínio
│   │   ├── user.go
│   │   ├── skin.go
│   │   ├── offer.go
│   │   └── trade.go
│   ├── auth/                       # Módulo de autenticação
│   │   ├── repository/
│   │   ├── usecase/
│   │   └── delivery/http/
│   ├── skins/                      # Módulo de skins
│   ├── offers/                     # Módulo de ofertas P2P
│   └── trades/                     # Módulo de histórico de trocas
├── pkg/
│   ├── auth/                       # Utilitários de autenticação
│   │   ├── jwt.go
│   │   ├── password.go
│   │   ├── totp.go
│   │   └── steam.go
│   ├── middleware/                 # Middlewares de segurança
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   └── security.go
│   ├── config/                     # Configuração
│   └── db/                         # Database clients
├── migrations/                     # SQL migrations
│   ├── 01_init_users.sql
│   ├── 02_init_skins.sql
│   ├── 03_init_offers_trades.sql
│   └── 04_seed_data.sql
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── go.mod
```

## 🚀 Como Rodar o Projeto

### Pré-requisitos
- Docker e Docker Compose instalados
- Go 1.24+ (opcional, para desenvolvimento local)

### Passo 1: Configurar Variáveis de Ambiente

```bash
cd backend
cp .env.example .env
```

Edite o arquivo `.env` e configure:
```env
JWT_SECRET=seu-secret-forte-aqui-min-32-caracteres
STEAM_API_KEY=sua-chave-steam-api  # Obtenha em https://steamcommunity.com/dev/apikey
```

### Passo 2: Subir com Docker Compose

```bash
docker-compose up --build
```

Isso irá:
1. ✅ Subir PostgreSQL na porta 5432
2. ✅ Subir Redis na porta 6379
3. ✅ Executar migrations automaticamente
4. ✅ Subir a API na porta 8080

### Passo 3: Verificar Health (HTTPS)

```bash
curl -k https://localhost:8443/healthz
```

Resposta esperada:
```json
{
  "status": "healthy",
  "timestamp": 1234567890
}
```

**Nota**: A flag `-k` ignora a validação do certificado SSL (somente para desenvolvimento)

## 📚 API Endpoints

### 🔐 Autenticação

#### Registro
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Login (1º Fator)
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

Resposta (sem 2FA):
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

Resposta (com 2FA habilitado):
```json
{
  "success": true,
  "message": "2FA code required",
  "data": {
    "requires_2fa": true
  }
}
```

#### Verificar 2FA (2º Fator)
```http
POST /api/auth/2fa/verify
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "123456"
}
```

#### Configurar 2FA
```http
POST /api/auth/2fa/setup
Authorization: Bearer <token>
```

Resposta:
```json
{
  "success": true,
  "message": "2FA setup initiated",
  "data": {
    "secret": "JBSWY3DPEHPK3PXP",
    "qr_url": "otpauth://totp/CS2%20P2P%20Skins:user@example.com?secret=..."
  }
}
```

#### Habilitar 2FA
```http
POST /api/auth/2fa/enable
Authorization: Bearer <token>
Content-Type: application/json

{
  "code": "123456"
}
```

#### Login via Steam OAuth
```http
GET /api/auth/steam/login
```
Retorna a URL de autenticação do Steam. Redirecione o usuário para essa URL.

#### Usuário Atual
```http
GET /api/auth/me
Authorization: Bearer <token>
```

### 🎮 Skins

#### Listar Skins Disponíveis
```http
GET /api/skins?limit=50&offset=0
Authorization: Bearer <token>
```

#### Obter Skin por ID
```http
GET /api/skins/1
Authorization: Bearer <token>
```

#### Listar Minhas Skins
```http
GET /api/skins/me
Authorization: Bearer <token>
```

#### Adicionar Skin ao Inventário (Demo)
```http
POST /api/skins/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "skin_id": 1,
  "quantity": 1
}
```

### 💱 Ofertas de Troca

#### Criar Oferta
```http
POST /api/offers
Authorization: Bearer <token>
Content-Type: application/json

{
  "skin_offered_id": 1,
  "skin_requested_id": 2
}
```

#### Listar Ofertas Abertas
```http
GET /api/offers?status=open&limit=50&offset=0
Authorization: Bearer <token>
```

#### Listar Minhas Ofertas
```http
GET /api/offers/me
Authorization: Bearer <token>
```

#### Aceitar Oferta
```http
POST /api/offers/1/accept
Authorization: Bearer <token>
```

#### Cancelar Oferta
```http
POST /api/offers/1/cancel
Authorization: Bearer <token>
```

### 📜 Histórico de Trocas

#### Listar Todas as Trocas
```http
GET /api/trades?limit=50&offset=0
Authorization: Bearer <token>
```

#### Listar Minhas Trocas
```http
GET /api/trades/me
Authorization: Bearer <token>
```

## 🔒 Explicação das Medidas de Segurança

### 1. Hash de Senha (Bcrypt)

**Por que não guardar senhas em texto puro?**
- Se o banco de dados for comprometido, as senhas ficam expostas
- Bcrypt usa salt automático e é resistente a ataques de força bruta
- Cost factor 12 torna cada tentativa mais lenta (dificulta rainbow tables)

```go
// pkg/auth/password.go
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
```

### 2. Autenticação de Dois Fatores (2FA)

**Como aumenta a segurança?**
- Mesmo se a senha for comprometida, o atacante precisa do código TOTP
- Códigos são temporários (30 segundos de validade)
- Baseado em TOTP (RFC 6238)

```go
// pkg/auth/totp.go
valid := totp.Validate(code, secret)
```

### 3. Proteção Contra Força Bruta

**Como mitiga ataques?**
- Registra todas as tentativas de login (sucesso/falha)
- Bloqueia conta após 5 tentativas falhas em 15 minutos
- Bloqueia IP após 10 tentativas falhas em 15 minutos
- Tabela de auditoria para análise posterior

```go
// internal/auth/usecase/usecase.go
failedAttempts, _ := repo.CountFailedAttempts(ctx, email, 15)
if failedAttempts >= 5 {
    return ErrAccountLocked
}
```

### 4. JWT (JSON Web Tokens)

**Por que usar JWT?**
- Stateless: servidor não precisa armazenar sessões
- Auto-contido: inclui informações do usuário (user_id, role)
- Expira após 24h (configurável)
- Assinado com HMAC-SHA256

```go
// pkg/auth/jwt.go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
```

### 5. Rate Limiting

**Como previne abuso?**
- Limita requisições por IP
- Endpoints de login: 5 req/minuto (strict)
- Endpoints de API: 30 req/minuto (moderate)
- Usa Token Bucket algorithm

```go
// pkg/middleware/ratelimit.go
limiter := rate.NewLimiter(rate.Limit(5.0/60.0), 5)
```

### 6. Middlewares de Segurança

**Headers adicionados:**
- `X-Content-Type-Options: nosniff` - Previne MIME sniffing
- `X-Frame-Options: DENY` - Previne clickjacking
- `X-XSS-Protection: 1; mode=block` - Habilita proteção XSS
- `Content-Security-Policy` - Controla recursos carregados
- `Referrer-Policy` - Controla informações do referrer

### 7. CORS Configurado Adequadamente

```go
// pkg/middleware/security.go
AllowOrigins: []string{"http://localhost:3000"},
AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
AllowHeaders: []string{"Authorization", "Content-Type"},
AllowCredentials: true,
```

## 🧪 Testando com Bruno

A API inclui uma coleção completa do Bruno para testar todos os endpoints de forma segura via HTTPS.

### Como Usar a Coleção Bruno

1. **Instale o [Bruno](https://www.usebruno.com/)**
2. **Abra a coleção**: `bruno/` na raiz do projeto
3. **Selecione o ambiente**: "Local" (já configurado para https://localhost:8443)
4. **Execute os testes** seguindo a ordem:

#### Fluxo de Testes:
1. **Registre um usuário** - `auth/1. Register`
2. **Faça login e pegue o token JWT** - `auth/2. Login` (token salvo automaticamente)
3. **Valide autenticação** - `auth/3. Get Current User`
4. **Configure 2FA** (opcional):
   - Setup 2FA (pega QR code)
   - Escaneie com Google Authenticator
   - Habilite 2FA com código
5. **Teste login com 2FA**:
   - Login retorna `requires_2fa: true`
   - Verify 2FA com código do app
6. **Adicione skins ao inventário** - `skins/4. Add Skin to Inventory`
7. **Crie ofertas de troca** - `offers/1. Create Offer`
8. **Aceite ofertas de outros usuários** - `offers/4. Accept Offer`
9. **Veja histórico de trocas** - `trades/2. List My Trades`

📖 **Documentação completa**: Veja `bruno/README.md` para detalhes sobre segurança, TLS e troubleshooting.

## 🔐 Steam OAuth Flow

1. Cliente chama `GET /api/auth/steam/login`
2. API retorna URL do Steam
3. Cliente redireciona usuário para Steam
4. Usuário autentica no Steam
5. Steam redireciona para `/api/auth/steam/callback`
6. API valida resposta do Steam (verifica state CSRF)
7. API cria/busca usuário pelo Steam ID
8. API retorna JWT token

## 🎓 Conceitos de Cyber Security Aplicados

| Conceito | Implementação |
|----------|---------------|
| **Confidencialidade** | Bcrypt para senhas, JWT para sessões, HTTPS/TLS |
| **Integridade** | HMAC-SHA256 em JWTs, validação de entrada, TLS |
| **Disponibilidade** | Rate limiting, proteção contra DoS |
| **Autenticação** | Email/senha + 2FA + OAuth2 Steam |
| **Autorização** | RBAC (user/admin roles) |
| **Auditoria** | Tabela login_attempts, logs estruturados |
| **Defense in Depth** | Múltiplas camadas de segurança |
| **Criptografia em Trânsito** | TLS 1.2+ com certificados SSL |

## 🚨 Limitações e Melhorias Futuras

### Limitações Atuais:
1. **Steam Integration**: Simplificado para demo (não consome API real de inventário)
2. **Email 2FA**: Não implementado (apenas TOTP/Google Authenticator)
3. **Session Management**: Implementar refresh tokens
4. **Password Policy**: Adicionar requisitos de complexidade
5. **Production Certificates**: Certificados auto-assinados (usar Let's Encrypt em produção)

### Melhorias Recomendadas:
1. ✅ Implementar refresh tokens (JWT de longa duração)
2. ✅ Adicionar logs estruturados (Zap/Zerolog)
3. ✅ Implementar WebSockets para notificações em tempo real
4. ✅ Cache com Redis para performance
5. ✅ Testes unitários e de integração
6. ✅ CI/CD pipeline
7. ✅ Monitoring e alertas (Prometheus + Grafana)
8. ✅ Web Application Firewall (WAF)
9. ✅ Implementar CAPTCHA em login
10. ✅ Backup automático do banco de dados

## 📖 Referências

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [RFC 6238 - TOTP](https://tools.ietf.org/html/rfc6238)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OAuth 2.0 RFC 6749](https://tools.ietf.org/html/rfc6749)
- [Bcrypt Algorithm](https://en.wikipedia.org/wiki/Bcrypt)
- [Steam OpenID](https://steamcommunity.com/dev)

## 👨‍💻 Autor

Desenvolvido como trabalho final da disciplina de Cyber Security.

## 📄 Licença

Este projeto é apenas para fins educacionais.

---

**⚠️ IMPORTANTE**: Este projeto foi desenvolvido para demonstrar conceitos de segurança. Em um ambiente de produção real, seria necessário:
- Auditoria de segurança profissional
- Testes de penetração
- Compliance com LGPD/GDPR
- Monitoramento contínuo de vulnerabilidades
- Backup e disaster recovery
- Documentação completa de processos de segurança
