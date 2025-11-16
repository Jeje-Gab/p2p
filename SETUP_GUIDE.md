# 🎮 CS2 P2P Skins Trading Platform - Guia de Configuração

## 📋 Índice
- [Correções Implementadas](#correções-implementadas)
- [Aplicando as Mudanças](#aplicando-as-mudanças)
- [Configuração da Steam API](#configuração-da-steam-api)
- [Tornando o P2P Mais Real](#tornando-o-p2p-mais-real)
- [Troubleshooting](#troubleshooting)

---

## ✅ Correções Implementadas

### 1. **Rate Limiting Ajustado**
**Problema**: Erro 429 Too Many Requests

**Solução**: Aumentado os limites de requisições em `backend/pkg/middleware/ratelimit.go`
- **StrictRateLimiter**: 5 req/min → **500 req/min** (endpoints de autenticação)
- **ModerateRateLimiter**: 30 req/min → **1000 req/min** (endpoints de API protegidos)

### 2. **Imagens das Skins Corrigidas**
**Problema**: URLs das imagens truncadas/quebradas

**Solução**: Atualizadas todas as URLs em `backend/migrations/04_seed_data.sql`
- Antes: `steamcommunity-a.akamaihd.net/economy/image/...` (truncadas)
- Agora: `community.cloudflare.steamstatic.com/economy/image/...` (URLs completas)

**Skins Incluídas:**
- AK-47 Redline
- AWP Asiimov
- M4A4 Howl
- AWP Dragon Lore
- USP-S Kill Confirmed
- Desert Eagle Printstream
- Desert Eagle Hypnotic
- Glock-18 Fade
- M4A4 The Emperor
- Desert Eagle Neo-Noir

### 3. **Espaçamento dos Botões Melhorado**
**Problema**: Botões sem espaçamento adequado

**Solução**: Ajustados gaps nos containers:
- `frontend/src/app/(dashboard)/offers/page.tsx` - gaps aumentados de 2 para 3
- `frontend/src/app/(dashboard)/dashboard/page.tsx` - spacing aumentado de 3 para 4

### 4. **CSS de Raridades**
Sistema de cores para raridades CS2 já implementado em `frontend/src/app/globals.css`:
- Consumer Grade: Cinza
- Industrial Grade: Azul Claro
- Mil-Spec Grade: Azul Escuro
- Restricted: Roxo Claro
- Classified: Roxo Escuro
- Covert: Vermelho
- Exceedingly Rare: Amarelo

---

## 🚀 Aplicando as Mudanças

### Passo 1: Parar os Servidores
```bash
# Pressione Ctrl+C em ambos os terminais (backend e frontend)
```

### Passo 2: Recriar o Banco de Dados
```bash
# Conectar ao PostgreSQL e recriar o schema
psql -U postgres -d p2p_db -c "DROP SCHEMA IF EXISTS p2p CASCADE;"

# Executar as migrations na ordem
psql -U postgres -d p2p_db -f backend/migrations/01_init_users.sql
psql -U postgres -d p2p_db -f backend/migrations/02_init_skins.sql
psql -U postgres -d p2p_db -f backend/migrations/03_init_offers_trades.sql
psql -U postgres -d p2p_db -f backend/migrations/04_seed_data.sql
```

**OU via script SQL único:**
```bash
psql -U postgres -d p2p_db << 'EOF'
DROP SCHEMA IF EXISTS p2p CASCADE;
\i backend/migrations/01_init_users.sql
\i backend/migrations/02_init_skins.sql
\i backend/migrations/03_init_offers_trades.sql
\i backend/migrations/04_seed_data.sql
EOF
```

### Passo 3: Reiniciar o Backend
```bash
cd backend
go run cmd/server/main.go
```

Você deve ver:
```
[BOOT] Environment: dev | Port: 8080
[HTTP] Server starting on port 8080
```

### Passo 4: Reiniciar o Frontend
Em outro terminal:
```bash
cd frontend
npm run dev
```

Você deve ver:
```
▲ Next.js 14.x.x
- Local:   http://localhost:3000
```

### Passo 5: Verificar
Acesse `http://localhost:3000` e as imagens devem carregar corretamente! ✨

---

## 🔑 Configuração da Steam API

### Por que usar a Steam API?
A integração com Steam permite:
- ✅ Login autêntico via Steam
- ✅ Buscar inventário real dos jogadores
- ✅ Criar trade offers reais
- ✅ Sincronizar skins automaticamente
- ✅ Validar propriedade das skins

### Passo 1: Obter Steam API Key

1. **Acesse**: https://steamcommunity.com/dev/apikey
2. **Faça login** com sua conta Steam
3. **Preencha o formulário:**
   - Domain Name: `localhost` (para desenvolvimento)
   - Agree to Steam Web API Terms: ✅
4. **Copie** sua API Key gerada

⚠️ **IMPORTANTE**: Nunca compartilhe ou commite sua API Key!

### Passo 2: Configurar no Backend

Edite o arquivo `backend/.env`:

```env
# Steam OAuth Configuration
STEAM_API_KEY=SUA-API-KEY-AQUI-123456789ABCDEF
STEAM_CALLBACK_URL=http://localhost:8080/api/auth/steam/callback
```

Se o arquivo `.env` não existir, copie do exemplo:
```bash
cd backend
cp .env.example .env
# Depois edite o .env com sua API Key
```

### Passo 3: Reiniciar o Backend
```bash
cd backend
go run cmd/server/main.go
```

### Passo 4: Testar Login Steam
1. Acesse o frontend: `http://localhost:3000`
2. Clique em "Login with Steam"
3. Você será redirecionado para Steam
4. Após autorizar, retornará ao site autenticado

---

## 🎯 Tornando o P2P Mais Real

### Funcionalidades Já Implementadas

#### 1. **Autenticação Steam** ✅
Localização: `backend/pkg/auth/steam.go`

**Recursos:**
- OpenID 2.0 Authentication
- OAuth flow completo
- Proteção CSRF com state
- Fetch de dados do usuário (Avatar, Nome, SteamID)

**Endpoints:**
- `GET /api/auth/steam/login` - Inicia login
- `GET /api/auth/steam/callback` - Callback OAuth

#### 2. **Sistema de Rate Limiting** ✅
Localização: `backend/pkg/middleware/ratelimit.go`

**Proteção contra:**
- DoS attacks
- Spam de requisições
- Abuse de API

### Próximas Funcionalidades para P2P Real

#### 1. **Buscar Inventário Real do CS2**

**API Endpoint:**
```
http://api.steampowered.com/IEconItems_730/GetPlayerItems/v0001/
```

**Parâmetros:**
- `key`: Sua Steam API Key
- `steamid`: SteamID64 do jogador

**Exemplo de Implementação:**
```go
// backend/internal/skins/usecase/steam_inventory.go
func (u *useCase) FetchSteamInventory(ctx context.Context, steamID string) ([]*entity.Skin, error) {
    apiURL := fmt.Sprintf(
        "http://api.steampowered.com/IEconItems_730/GetPlayerItems/v0001/?key=%s&steamid=%s",
        u.steamAPIKey, steamID,
    )

    // Fazer requisição HTTP
    // Parsear resposta JSON
    // Mapear para entity.Skin
    // Salvar no banco de dados
}
```

#### 2. **Trade Offers Reais**

**API Endpoints Steam:**
- `IEconService/GetTradeOffers` - Listar offers
- `IEconService/GetTradeOffer` - Detalhes de uma offer
- `IEconService/DeclineTradeOffer` - Recusar offer

**Limitações:**
- Criar trade offers requer **Steam Trade API Key** (diferente da Web API Key)
- Necessita autenticação via Steam Guard Mobile
- Precisa de Trade URL do usuário

**Documentação Oficial:**
https://developer.valvesoftware.com/wiki/Steam_Web_API#GetTradeOffers_.28v1.29

#### 3. **Preços de Mercado em Tempo Real**

**API Endpoint:**
```
http://steamcommunity.com/market/priceoverview/
```

**Parâmetros:**
- `appid=730` (CS2)
- `currency=7` (USD)
- `market_hash_name=AK-47%20%7C%20Redline%20%28Field-Tested%29`

**Exemplo:**
```typescript
// frontend/src/services/market.service.ts
async function getSkinPrice(marketHashName: string) {
    const url = `http://steamcommunity.com/market/priceoverview/?appid=730&currency=7&market_hash_name=${encodeURIComponent(marketHashName)}`;
    const response = await fetch(url);
    const data = await response.json();
    return {
        lowest_price: data.lowest_price,
        median_price: data.median_price,
        volume: data.volume
    };
}
```

#### 4. **Notificações em Tempo Real**

**Tecnologias:**
- WebSockets (para notificações instantâneas)
- Server-Sent Events (SSE)
- Redis Pub/Sub

**Casos de uso:**
- Notificar quando alguém aceitar sua offer
- Alertas de novas offers para suas skins
- Mudanças de preço de mercado

#### 5. **Sistema de Reputação**

**Funcionalidades:**
- Rating de usuários (1-5 estrelas)
- Comentários em perfis
- Histórico de trades completos
- Badges por número de trades

**Schema adicional:**
```sql
CREATE TABLE p2p.user_ratings (
    id SERIAL PRIMARY KEY,
    from_user_id INT REFERENCES p2p.users(id),
    to_user_id INT REFERENCES p2p.users(id),
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    trade_id INT REFERENCES p2p.trades(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔧 Troubleshooting

### Erro: "429 Too Many Requests"
**Solução**: Reinicie o backend após as alterações no rate limiter
```bash
cd backend
go run cmd/server/main.go
```

### Imagens ainda não carregam
**Verificar:**
1. Database foi recriado com as novas URLs?
2. Frontend está apontando para `http://localhost:8080`?
3. CORS configurado corretamente no backend?

**Verificar URLs no banco:**
```sql
SELECT id, name, image_url FROM p2p.skins LIMIT 3;
```

### Erro: "Cannot connect to database"
**Verificar:**
```bash
# PostgreSQL está rodando?
pg_isready -U postgres

# Verificar variável de ambiente
echo $DB_DSN  # ou no .env: DB_DSN
```

### Botões ainda sem espaçamento
**Solução**: Limpar cache do Next.js
```bash
cd frontend
rm -rf .next
npm run dev
```

### Steam Login não funciona
**Verificar:**
1. `STEAM_API_KEY` configurado no `.env`?
2. `STEAM_CALLBACK_URL` está correto?
3. Backend está rodando na porta 8080?
4. Verificar logs do backend para erros

**Testar callback:**
```bash
# No navegador, verificar se este endpoint existe:
http://localhost:8080/api/auth/steam/callback
```

---

## 📚 Referências Úteis

### Documentação Oficial
- [Steam Web API Documentation](https://developer.valvesoftware.com/wiki/Steam_Web_API)
- [Steam API Key Registration](https://steamcommunity.com/dev/apikey)
- [CS2 Item Schema](https://github.com/SteamDatabase/GameTracking-CS2)

### APIs Úteis
- **Steam Market API**: Preços em tempo real
- **Steam Inventory API**: Inventário dos jogadores
- **Steam Trade API**: Gerenciar trade offers
- **Steam User API**: Informações de perfil

### Comunidade
- [r/csgomarketforum](https://reddit.com/r/csgomarketforum) - Trading discussion
- [SteamDB](https://steamdb.info/) - Database de jogos/items Steam
- [CS2 Backpack](https://csbackpack.net/) - Preços e estatísticas

---

## 🎉 Próximos Passos

1. ✅ **Aplicar todas as correções** (migrations + restart)
2. ✅ **Configurar Steam API Key**
3. ⬜ **Implementar fetch de inventário real**
4. ⬜ **Adicionar preços de mercado**
5. ⬜ **Implementar sistema de notificações**
6. ⬜ **Adicionar sistema de reputação**

---

**Desenvolvido com ❤️ para a comunidade CS2**

Para questões ou sugestões, abra uma issue no repositório!
