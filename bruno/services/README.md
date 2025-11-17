# External Services API Collection

Esta collection contém exemplos de uso da **API externa** do CS2 P2P Skins para integração com sistemas externos.

## 🔐 Autenticação

Todos os endpoints requerem autenticação via **API Key** no header:

```
X-API-Key: your-api-key-here
```

## 📋 Configuração

### 1. Criar API Key (Usuário Autenticado)

Use a collection **API Keys Management** para criar uma nova key:

1. Faça login para obter JWT token
2. Execute `POST /api/api-keys` com body:
   ```json
   {
     "name": "Mobile App Production",
     "description": "API key for external access"
   }
   ```
3. **IMPORTANTE**: Copie a key retornada - ela será mostrada apenas uma vez!

### 2. Distribuir para Cliente Externo

Envie a API key gerada para o sistema/cliente que vai consumir sua API.

### 3. Cliente Usa a Key

O cliente externo adiciona a key no header de todas as requisições:

```bash
curl -H "X-API-Key: sk_live_abc123..." \
  https://localhost:8443/services/skins/available
```

**Nota**: Não é mais necessário configurar keys no `.env` - tudo é gerenciado dinamicamente no banco de dados!

## 📡 Endpoints Disponíveis

### 1. Get Available Skins
`GET /services/skins/available`

Lista todas as skins disponíveis com paginação.

**Query Params:**
- `limit` (default: 100, max: 1000)
- `offset` (default: 0)

---

### 2. Get Skin by ID
`GET /services/skins/:id`

Retorna detalhes de uma skin específica.

---

### 3. Get Open Offers
`GET /services/offers/open`

Lista todas as ofertas de trade abertas.

**Query Params:**
- `limit` (default: 100, max: 1000)
- `offset` (default: 0)

---

### 4. Get Recent Trades
`GET /services/trades/recent`

Retorna histórico de trades completadas.

**Query Params:**
- `limit` (default: 50, max: 500)
- `offset` (default: 0)

---

### 5. Get Market Stats
`GET /services/market/stats`

Estatísticas agregadas do marketplace:
- Total de skins
- Total de trades
- Skins disponíveis para trade

## 🛡️ Segurança

- ✅ HTTPS obrigatório
- ✅ Rate limiting: 30 requisições/minuto
- ✅ Validação de API Key
- ✅ Security headers habilitados
- ✅ CORS configurável

## 💡 Casos de Uso

Esta API externa permite que sistemas third-party:

- **Dashboards**: Criar dashboards customizados com dados do marketplace
- **Analytics**: Análise de tendências de mercado e preços
- **Bots**: Automação de monitoramento de ofertas
- **Integrações**: Conectar com outros sistemas e plataformas
- **Relatórios**: Gerar relatórios de negócio
- **Mobile Apps**: Aplicativos mobile podem consumir os dados

## 🔒 Segurança em Produção

**IMPORTANTE**: Em produção, sempre:

1. Use HTTPS (porta 443)
2. Gere API keys fortes (min 32 caracteres)
3. Rotacione as keys periodicamente
4. Monitore uso via logs
5. Configure CORS adequadamente
6. Nunca commite API keys no git

## 📊 Exemplo de Resposta

Todas as respostas seguem o padrão:

```json
{
  "success": true,
  "message": "Optional message",
  "data": {
    // Response data
  }
}
```

Em caso de erro:

```json
{
  "success": false,
  "message": "Error description"
}
```

## 🚀 Como Testar

1. Certifique-se que o backend está rodando (porta 8443)
2. Configure sua API key no ambiente Bruno
3. Execute qualquer request da collection
4. Verifique a resposta JSON

## 📝 Notas

- A API só é habilitada se houver API keys configuradas no `.env`
- Sem keys configuradas, os endpoints retornam 404
- Rate limiting é aplicado por IP
- Paginação é obrigatória para endpoints de lista
