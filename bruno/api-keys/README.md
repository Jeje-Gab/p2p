# API Keys Management Collection

Esta collection permite que você **gerencie API keys** para integração com sistemas externos.

## 🎯 Visão Geral

O sistema de API Keys permite que usuários autenticados criem chaves de acesso para sistemas third-party consumirem a API externa (`/services/*` endpoints).

## 🔄 Fluxo de Uso

### 1️⃣ Como Administrador/Usuário (Você)

1. **Login**: Faça login para obter um JWT token
2. **Criar API Key**: Use `1. Create API Key` para gerar uma nova key
3. **Copiar Key**: A key completa é mostrada **apenas uma vez** - copie imediatamente!
4. **Distribuir**: Envie a key para o sistema externo que vai consumir sua API

### 2️⃣ Como Cliente Externo

1. **Receber Key**: Receba a API key do administrador
2. **Usar Key**: Adicione no header `X-API-Key` em todas as requisições para `/services/*`
3. **Consumir API**: Acesse os endpoints públicos (skins, offers, trades, stats)

## 📋 Endpoints Disponíveis

### 1. Create API Key
`POST /api/api-keys`

Gera uma nova API key.

**Autenticação**: JWT (Bearer token)

**Body**:
```json
{
  "name": "Mobile App Production",
  "description": "Optional description",
  "expires_at": "2026-12-31T23:59:59Z"  // Optional
}
```

**⚠️ IMPORTANTE**: A key completa é retornada apenas neste momento!

---

### 2. List My API Keys
`GET /api/api-keys`

Lista todas as suas API keys.

**Autenticação**: JWT (Bearer token)

Mostra apenas o `key_prefix` para segurança (ex: `sk_live_abc123...`)

---

### 3. Revoke API Key
`POST /api/api-keys/:id/revoke`

Desativa uma API key sem deletá-la.

**Autenticação**: JWT (Bearer token)

Útil para suspender acesso temporariamente.

---

### 4. Delete API Key
`DELETE /api/api-keys/:id`

Remove permanentemente uma API key.

**Autenticação**: JWT (Bearer token)

⚠️ **Ação irreversível!**

## 🔐 Segurança

### Armazenamento Seguro
- Keys são armazenadas como **SHA256 hash** no banco
- Apenas o **hash** é salvo - a key original nunca é armazenada
- Impossível recuperar a key original - por isso mostramos apenas uma vez!

### Formato da Key
```
sk_live_<64 caracteres hexadecimais>
```

Exemplo:
```
sk_live_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6a7b8c9d0e1f2
```

### Validação
Quando um sistema externo usa a key:
1. Key é extraída do header `X-API-Key`
2. Hash da key é calculado
3. Hash é buscado no banco de dados
4. Verificações:
   - ✅ Key existe?
   - ✅ Está ativa (`is_active = true`)?
   - ✅ Não expirou?
5. Se tudo OK, acesso é liberado
6. `last_used_at` é atualizado automaticamente

## 📊 Monitoramento

Cada API key rastreia:
- **Criação**: `created_at`
- **Último uso**: `last_used_at` (atualizado a cada requisição)
- **Expiração**: `expires_at` (opcional)
- **Status**: `is_active` (true/false)

## 💡 Casos de Uso

### Cenário 1: Mobile App
```json
{
  "name": "Mobile App - Production",
  "description": "Main API key for iOS and Android apps",
  "expires_at": null  // Sem expiração
}
```

### Cenário 2: Dashboard Analytics
```json
{
  "name": "Analytics Dashboard",
  "description": "Read-only access for business intelligence",
  "expires_at": "2026-01-01T00:00:00Z"  // Expira em 1 ano
}
```

### Cenário 3: Parceiro Temporário
```json
{
  "name": "Partner XYZ - Trial",
  "description": "30-day trial access",
  "expires_at": "2025-12-16T00:00:00Z"  // 30 dias
}
```

## 🔄 Rotação de Keys

Boas práticas para rotacionar API keys:

1. **Criar nova key**
   ```
   POST /api/api-keys
   ```

2. **Distribuir nova key** para o sistema externo

3. **Sistema externo atualiza** para usar nova key

4. **Validar** que a nova key está funcionando

5. **Revogar key antiga**
   ```
   POST /api/api-keys/{old_id}/revoke
   ```

6. **(Opcional) Deletar** após período de segurança
   ```
   DELETE /api/api-keys/{old_id}
   ```

## 🚨 Em Caso de Comprometimento

Se uma API key for comprometida:

1. **Revogue imediatamente**:
   ```
   POST /api/api-keys/{id}/revoke
   ```

2. **Crie nova key** para o sistema legítimo

3. **Distribua nova key** com segurança

4. **Investigue** uso em `last_used_at`

## 🧪 Testando

### Passo a Passo

1. **Faça login** para obter JWT token:
   ```
   POST /api/auth/login
   ```

2. **Crie uma API key**:
   ```
   POST /api/api-keys
   Body: { "name": "Test Key" }
   ```

3. **Copie a key** do response (campo `key`)

4. **Teste nos services**:
   ```
   GET /services/skins/available
   Header: X-API-Key: sk_live_...
   ```

5. **Verifique último uso**:
   ```
   GET /api/api-keys
   ```
   O campo `last_used_at` deve estar atualizado!

## 📝 Script Automático

O request "1. Create API Key" inclui um **script pós-resposta** que:
- ✅ Extrai a key automaticamente
- ✅ Salva na variável de ambiente `api_key`
- ✅ Mostra no console para você copiar

Após criar uma key, ela já estará disponível em `{{api_key}}` para usar nos requests de `/services/*`!

## ⚙️ Configuração Antiga (.env) - DEPRECIADA

Antes usávamos:
```env
EXTERNAL_API_KEYS=sk_prod_abc123,sk_prod_def456
```

**Agora isso NÃO É MAIS NECESSÁRIO!** 🎉

Todas as keys são gerenciadas dinamicamente no banco de dados.
