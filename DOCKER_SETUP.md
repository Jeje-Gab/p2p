# 🐳 Docker Setup - CS2 P2P Skins Trading Platform

Este guia mostra como rodar toda a aplicação usando Docker em **qualquer máquina** com um único comando!

## 📋 Pré-requisitos

Apenas Docker e Docker Compose instalados:

- **Docker**: https://www.docker.com/get-started
- **Docker Compose**: Incluído no Docker Desktop (Windows/Mac) ou instale separadamente no Linux

Verifique a instalação:
```bash
docker --version
docker-compose --version
```

## 🚀 Início Rápido (One-Command Setup)

### 1. Clone o repositório
```bash
git clone <repo-url>
cd p2p
```

### 2. Configure as variáveis de ambiente
```bash
# Copie o arquivo de exemplo
cp .env.docker.example .env

# Edite o arquivo .env e configure:
# - STEAM_API_KEY (obrigatório - pegue em https://steamcommunity.com/dev/apikey)
# - Outras variáveis conforme necessário
```

### 3. Inicie tudo com um único comando
```bash
docker-compose up -d --build
```

**Pronto!** 🎉

A aplicação estará rodando em:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **MinIO Console**: http://localhost:9001 (usuário: `minio`, senha: `minio123`)
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 📦 O que é criado automaticamente

Quando você roda `docker-compose up -d --build`, o sistema:

### ✅ Banco de Dados (PostgreSQL)
1. Cria o container PostgreSQL
2. Inicializa o schema do banco
3. Espera estar pronto para aceitar conexões

### ✅ MinIO (Object Storage)
1. Cria o container MinIO
2. Inicia o serviço de armazenamento
3. Console web disponível na porta 9001

### ✅ Redis (Cache)
1. Cria o container Redis
2. Configura senha de acesso
3. Pronto para cache de sessões

### ✅ Backend (Go API)
1. **Espera** PostgreSQL, MinIO e Redis estarem prontos
2. **Roda as migrations** automaticamente
3. **Cria o bucket MinIO** `net.public.p2p`
4. **Gera imagens placeholder** das skins
5. Inicia a API na porta 8080

### ✅ Frontend (Next.js)
1. **Espera** o backend estar pronto
2. Build otimizado da aplicação
3. Inicia o servidor web na porta 3000

## 🎯 Comandos Úteis

### Iniciar todos os serviços
```bash
docker-compose up -d
```

### Ver logs em tempo real
```bash
# Todos os serviços
docker-compose logs -f

# Apenas backend
docker-compose logs -f backend

# Apenas frontend
docker-compose logs -f frontend
```

### Parar todos os serviços
```bash
docker-compose down
```

### Parar e remover volumes (limpar tudo)
```bash
docker-compose down -v
```

### Rebuild completo (útil após mudanças no código)
```bash
docker-compose up -d --build --force-recreate
```

### Ver status dos containers
```bash
docker-compose ps
```

### Acessar o shell de um container
```bash
# Backend
docker-compose exec backend sh

# Frontend
docker-compose exec frontend sh

# PostgreSQL
docker-compose exec postgres psql -U postgres
```

### Rodar migrations manualmente (se necessário)
```bash
docker-compose exec backend /usr/local/bin/migrate -path=/app/migrations -database="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up
```

## 🔍 Troubleshooting

### Erro: "Port already in use"
```bash
# Verifique quais processos estão usando as portas
netstat -ano | findstr :3000   # Windows
lsof -i :3000                  # Linux/Mac

# Mude as portas no arquivo .env
FRONTEND_PORT=3001
BACKEND_PORT=8081
```

### Backend não conecta ao banco
```bash
# Verifique se o PostgreSQL está healthy
docker-compose ps

# Veja os logs do backend
docker-compose logs backend

# Reinicie o backend
docker-compose restart backend
```

### Imagens das skins não aparecem
```bash
# Verifique se o MinIO está rodando
docker-compose ps minio

# Acesse o console do MinIO
# http://localhost:9001
# Usuário: minio
# Senha: minio123

# Verifique se o bucket existe
docker-compose exec backend sh
# curl http://minio:9000/net.public.p2p/
```

### Frontend não conecta ao backend
```bash
# Verifique a variável NEXT_PUBLIC_API_URL no .env
# Deve ser: http://localhost:8080/api

# Rebuild o frontend
docker-compose up -d --build frontend
```

### Reset completo do ambiente
```bash
# Para tudo e remove volumes
docker-compose down -v

# Remove imagens antigas
docker-compose down --rmi all

# Rebuild do zero
docker-compose up -d --build
```

## 📊 Health Checks

Todos os serviços têm health checks configurados:

- **PostgreSQL**: `pg_isready`
- **Redis**: `redis-cli ping`
- **MinIO**: `/minio/health/live`
- **Backend**: `GET /healthz`
- **Frontend**: Verifica resposta HTTP

Você pode verificar o status:
```bash
docker-compose ps

# Ou inspecionar um container específico
docker inspect --format='{{.State.Health.Status}}' p2p-backend
```

## 🌐 Acessos

### Serviços Web
| Serviço | URL | Credenciais |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | - |
| Backend API | http://localhost:8080 | - |
| MinIO Console | http://localhost:9001 | minio / minio123 |

### APIs
| Endpoint | Descrição |
|----------|-----------|
| `GET /healthz` | Health check do backend |
| `GET /api/skins` | Lista de skins |
| `GET /api/offers` | Ofertas disponíveis |
| `POST /api/auth/login` | Login de usuário |

### Bancos de Dados
| Serviço | Host | Porta | Usuário | Senha |
|---------|------|-------|---------|-------|
| PostgreSQL | localhost | 5432 | postgres | postgres |
| Redis | localhost | 6379 | - | redis123 |

## 📁 Estrutura de Volumes

Os dados são persistidos nos seguintes volumes Docker:

```
p2p_postgres_data    → Banco de dados PostgreSQL
p2p_minio_data       → Arquivos do MinIO (imagens das skins)
p2p_redis_data       → Cache do Redis
```

Para fazer backup:
```bash
# Backup do PostgreSQL
docker-compose exec postgres pg_dump -U postgres postgres > backup.sql

# Backup do MinIO
docker-compose exec minio mc mirror /data ./minio-backup
```

## 🔐 Segurança

### Produção

Antes de fazer deploy em produção:

1. **Mude todas as senhas** no `.env`
2. **Configure TLS** (`TLS_ENABLED=true`)
3. **Use secrets** em vez de variáveis de ambiente
4. **Configure CORS** adequadamente
5. **Use um JWT_SECRET** forte e único
6. **Configure rate limiting** apropriado

### Secrets do Docker (Produção)

Para produção, use Docker secrets:

```yaml
secrets:
  db_password:
    external: true
  jwt_secret:
    external: true
```

## 🔄 CI/CD

Exemplo de pipeline GitLab CI:

```yaml
build:
  stage: build
  script:
    - docker-compose build
    - docker-compose push

deploy:
  stage: deploy
  script:
    - docker-compose pull
    - docker-compose up -d
```

## 📝 Customização

### Mudando Portas

Edite `.env`:
```env
FRONTEND_PORT=3001
BACKEND_PORT=8081
MINIO_PORT=9002
```

### Adicionando Serviços

Edite `docker-compose.yml`:
```yaml
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    depends_on:
      - frontend
      - backend
```

## 🎓 Recursos Adicionais

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Next.js Docker Deployment](https://nextjs.org/docs/deployment#docker-image)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## 💡 Dicas

1. **Use `.env` file**: Nunca commite o `.env` no Git
2. **Multi-stage builds**: Reduz o tamanho das imagens
3. **Health checks**: Garante que serviços estão realmente prontos
4. **Volumes**: Use para persistência de dados
5. **Networks**: Isola serviços por ambiente

---

**Precisa de ajuda?** Abra uma issue no repositório! 🚀
