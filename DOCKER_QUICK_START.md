# 🚀 Docker Quick Start

Rodar toda a aplicação em **3 comandos**:

## 1️⃣ Configure as variáveis de ambiente

```bash
cp .env.docker.example .env
```

Edite o `.env` e adicione sua **Steam API Key**:
```env
STEAM_API_KEY=SUA_KEY_AQUI
```

Obtenha sua key em: https://steamcommunity.com/dev/apikey

## 2️⃣ Suba todos os serviços

```bash
docker-compose up -d --build
```

## 3️⃣ Acesse a aplicação

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (minio/minio123)

---

## 📋 O que acontece automaticamente:

✅ PostgreSQL database é criado
✅ MinIO object storage é configurado
✅ Redis cache é inicializado
✅ **Migrations do banco são executadas**
✅ **Bucket MinIO é criado** (net.public.p2p)
✅ **Imagens das skins são geradas**
✅ Backend API inicia na porta 8080
✅ Frontend inicia na porta 3000

---

## 🛠️ Comandos Úteis

### Ver logs em tempo real
```bash
docker-compose logs -f
```

### Parar todos os serviços
```bash
docker-compose down
```

### Rebuild completo (após mudanças no código)
```bash
docker-compose up -d --build --force-recreate
```

### Ver status dos containers
```bash
docker-compose ps
```

### Limpar tudo (⚠️ deleta dados)
```bash
docker-compose down -v
```

---

## 🎯 Usando o Makefile (Linux/Mac)

Se você tem `make` instalado:

```bash
make setup     # Configura .env
make up        # Inicia tudo
make logs      # Ver logs
make down      # Para tudo
make rebuild   # Rebuild completo
make clean     # Limpa tudo
make help      # Lista todos os comandos
```

---

## 🔍 Troubleshooting

### Porta já em uso?
Edite `.env` e mude as portas:
```env
FRONTEND_PORT=3001
BACKEND_PORT=8081
```

### Imagens não aparecem?
Verifique MinIO: http://localhost:9001

### Backend não conecta?
```bash
docker-compose logs backend
docker-compose restart backend
```

---

## 📚 Documentação Completa

Veja [DOCKER_SETUP.md](./DOCKER_SETUP.md) para documentação detalhada.

---

**É isso! Simples assim! 🎉**
