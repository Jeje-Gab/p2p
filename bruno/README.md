# CS2 P2P Skins API - Bruno Collection

Esta é a coleção completa da API do sistema de troca P2P de skins do CS2, documentada para uso com o [Bruno](https://www.usebruno.com/).

## 📋 Requisitos

- [Bruno](https://www.usebruno.com/) instalado
- Backend rodando em HTTPS (porta 8443)
- PostgreSQL e Redis configurados

## 🚀 Como Usar

### 1. Abrir a Coleção no Bruno

1. Abra o Bruno
2. Clique em "Open Collection"
3. Selecione a pasta `bruno/` deste projeto
4. A coleção será carregada com todos os endpoints

### 2. Configurar Variáveis de Ambiente

A coleção já vem com o ambiente "Local" pré-configurado:

```
base_url: https://localhost:8443
token: (será preenchido automaticamente após login)
user_email: user@example.com
user_password: SecurePass123!
```

Para usar suas próprias credenciais, edite o arquivo `environments/Local.bru`.

### 3. Fluxo de Teste Recomendado

#### Passo 1: Autenticação Básica

1. **Register** - Crie um novo usuário
2. **Login** - Faça login (token será salvo automaticamente)
3. **Get Current User** - Valide que está autenticado

#### Passo 2: Configurar 2FA (Opcional)

1. **Setup 2FA** - Obtenha o QR code
2. Escaneie o QR code com Google Authenticator
3. **Enable 2FA** - Habilite 2FA com o código do app
4. **Logout e Login novamente** - Login retornará `requires_2fa: true`
5. **Verify 2FA** - Complete o login com código 2FA

#### Passo 3: Gerenciar Inventário

1. **List All Skins** - Veja skins disponíveis
2. **Get Skin by ID** - Detalhes de uma skin específica
3. **Add Skin to Inventory** - Adicione skins ao seu inventário
4. **Get My Inventory** - Veja suas skins

#### Passo 4: Criar e Gerenciar Ofertas

1. **Create Offer** - Crie uma oferta de troca
2. **List Open Offers** - Veja ofertas disponíveis
3. **List My Offers** - Veja suas ofertas
4. **Accept Offer** - Aceite uma oferta (use outro usuário)
5. **Cancel Offer** - Cancele uma oferta sua

#### Passo 5: Histórico de Trocas

1. **List All Trades** - Veja todas as trocas completadas
2. **List My Trades** - Veja suas trocas

## 🔐 Segurança Implementada

### 1. HTTPS/TLS
- Todos os endpoints usam HTTPS
- Certificado auto-assinado para desenvolvimento
- Criptografia em trânsito (TLS 1.2+)

### 2. Autenticação JWT
- Token Bearer em todas as requisições autenticadas
- Token expira em 24h
- Assinado com HMAC-SHA256

### 3. Autenticação Multi-Fator (2FA)
- TOTP (RFC 6238)
- Compatível com Google Authenticator
- Códigos de 6 dígitos com 30s de validade

### 4. Rate Limiting
- Endpoints de auth: 5 req/min (strict)
- Endpoints de API: 30 req/min (moderate)
- Baseado em IP do cliente

### 5. Headers de Segurança
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy: default-src 'self'`
- `Referrer-Policy: strict-origin-when-cross-origin`

### 6. Proteção Contra Força Bruta
- Bloqueio após 5 tentativas falhas em 15 min
- Rastreamento por IP e por usuário
- Tabela de auditoria de tentativas

## 📁 Estrutura da Coleção

```
bruno/
├── bruno.json                  # Configuração da coleção
├── environments/
│   └── Local.bru              # Variáveis de ambiente
├── auth/
│   ├── 1. Register.bru
│   ├── 2. Login.bru
│   ├── 3. Get Current User.bru
│   ├── 4. Setup 2FA.bru
│   ├── 5. Enable 2FA.bru
│   ├── 6. Verify 2FA.bru
│   └── 7. Steam Login.bru
├── skins/
│   ├── 1. List All Skins.bru
│   ├── 2. Get Skin by ID.bru
│   ├── 3. Get My Inventory.bru
│   └── 4. Add Skin to Inventory.bru
├── offers/
│   ├── 1. Create Offer.bru
│   ├── 2. List Open Offers.bru
│   ├── 3. List My Offers.bru
│   ├── 4. Accept Offer.bru
│   └── 5. Cancel Offer.bru
└── trades/
    ├── 1. List All Trades.bru
    └── 2. List My Trades.bru
```

## 🧪 Testando a Segurança

### Teste 1: Verificar HTTPS
```bash
# Deve falhar (HTTP não permitido)
curl http://localhost:8080/healthz

# Deve funcionar (HTTPS)
curl -k https://localhost:8443/healthz
```

### Teste 2: Verificar Autenticação
```bash
# Deve retornar 401 Unauthorized
curl -k https://localhost:8443/api/skins

# Deve funcionar com token
curl -k https://localhost:8443/api/skins \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Teste 3: Verificar Rate Limiting
Execute o mesmo endpoint 10 vezes rapidamente:
```bash
for i in {1..10}; do
  curl -k https://localhost:8443/api/auth/login \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}'
done
```
Deve retornar `429 Too Many Requests` após algumas tentativas.

### Teste 4: Verificar Headers de Segurança
```bash
curl -k -I https://localhost:8443/healthz
```
Deve incluir:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-Xss-Protection: 1; mode=block`

## 🔧 Troubleshooting

### Erro: "certificate verify failed"
O certificado é auto-assinado. No Bruno, você pode desabilitar a verificação SSL:
- Settings → SSL/TLS → Disable SSL verification (apenas para desenvolvimento)

### Erro: "401 Unauthorized"
- Verifique se fez login
- Verifique se o token foi salvo no ambiente (variável `token`)
- Token expira em 24h - faça login novamente

### Erro: "429 Too Many Requests"
- Rate limit atingido
- Aguarde 1 minuto e tente novamente
- Para desenvolvimento, você pode aumentar os limites no código

## 📚 Documentação da API

Cada requisição no Bruno inclui documentação detalhada na aba "Docs":
- Descrição do endpoint
- Parâmetros de entrada
- Exemplo de resposta
- Regras de negócio
- Possíveis erros

## 🎓 Demonstração de Segurança

Esta coleção demonstra:

✅ **API REST completa e segura**
✅ **Autenticação via JWT**
✅ **Criptografia em trânsito (HTTPS/TLS)**
✅ **Possibilidade de certificados digitais**
✅ **Testes práticos com ferramenta (Bruno)**
✅ **Proteção e integridade dos dados**

---

**Desenvolvido para disciplina de Cyber Security**
