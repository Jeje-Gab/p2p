# MinIO Storage Setup - Imagens das Skins

Este documento explica como o sistema de armazenamento de imagens funciona com MinIO.

## 📦 Como Funciona

O sistema foi configurado para armazenar automaticamente todas as imagens das skins no MinIO, um serviço de armazenamento de objetos compatível com S3.

### Bucket Criado
- **Nome**: `net.public.p2p`
- **Permissões**: Público (leitura para todos)
- **Estrutura**: `skins/{arma}/{nome}.svg`

### Inicialização Automática

Toda vez que o backend iniciar, ele **automaticamente**:

1. ✅ Conecta ao MinIO
2. ✅ Verifica se o bucket `net.public.p2p` existe
3. ✅ Cria o bucket se não existir
4. ✅ Define permissões públicas no bucket
5. ✅ Verifica quais skins não têm imagens no MinIO
6. ✅ Cria imagens placeholder SVG para skins sem imagem
7. ✅ Atualiza o banco de dados com os URLs corretos

**Você não precisa fazer NADA manualmente!** 🎉

## 🚀 Rodando do Zero

### Pré-requisitos
1. MinIO rodando em `localhost:9000`
2. Credenciais configuradas no `.env`:
```env
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
MINIO_BUCKET=net.public.p2p
MINIO_USE_SSL=false
MINIO_PUBLIC_URL=http://localhost:9000
```

### Inicialização
```bash
# 1. Certifique-se que o MinIO está rodando
# (verifique se http://localhost:9000 está acessível)

# 2. Rode as migrations do banco de dados
# (isso inclui a migration 05_update_image_urls.sql)

# 3. Inicie o backend
cd backend
go run cmd/server/main.go

# O backend vai:
# - Criar o bucket automaticamente
# - Gerar todas as imagens das skins
# - Atualizar o banco de dados
```

Você verá logs como:
```
[MINIO] Connected to localhost:9000 | Bucket: net.public.p2p
[MINIO] Checking for missing skin images...
[MINIO] Created placeholder for: Redline (AK-47)
[MINIO] Created placeholder for: Asiimov (AWP)
...
[MINIO] Placeholder initialization complete: 10 created, 0 skipped
```

Na próxima vez que rodar:
```
[MINIO] All skin images already exist (10 total)
```

## 🎨 Imagens Geradas

As imagens são SVGs bonitos e coloridos com:
- **Gradientes baseados na raridade** (dourado para Exceedingly Rare, vermelho para Covert, etc.)
- **Nome da skin** em destaque
- **Tipo de arma** no topo
- **Indicador de raridade** com a cor correspondente
- **Design moderno** com padrões e sombras

### Estrutura das URLs
Padrão: `http://localhost:9000/net.public.p2p/skins/{arma}/{nome}.svg`

Exemplos:
- `http://localhost:9000/net.public.p2p/skins/AK-47/Redline.svg`
- `http://localhost:9000/net.public.p2p/skins/AWP/Dragon-Lore.svg`
- `http://localhost:9000/net.public.p2p/skins/Desert-Eagle/Hypnotic.svg`

## 🔧 Scripts Manuais (Opcional)

Se precisar recriar todas as imagens manualmente:

```bash
# Criar placeholders para todas as skins
cd backend
go run cmd/create-placeholders/main.go
```

## 🗄️ Estrutura do Bucket

```
net.public.p2p/
├── skins/
│   ├── AK-47/
│   │   └── Redline.svg
│   ├── AWP/
│   │   ├── Asiimov.svg
│   │   └── Dragon-Lore.svg
│   ├── M4A4/
│   │   ├── Howl.svg
│   │   └── The-Emperor.svg
│   ├── Desert-Eagle/
│   │   ├── Hypnotic.svg
│   │   ├── Neo-Noir.svg
│   │   └── Printstream.svg
│   ├── USP-S/
│   │   └── Kill-Confirmed.svg
│   └── Glock-18/
│       └── Fade.svg
```

## ✅ Vantagens

1. **Automático**: Não precisa rodar scripts manualmente
2. **Idempotente**: Pode rodar múltiplas vezes sem problemas
3. **Resiliente**: Se o MinIO estiver offline, o backend apenas avisa mas continua funcionando
4. **Escalável**: Fácil adicionar novas skins - só inserir no banco e reiniciar
5. **Performance**: Imagens servidas diretamente do MinIO (cache, CDN-ready)

## 🔍 Troubleshooting

### MinIO não está conectando
```
[MINIO] Warning: Failed to initialize MinIO: ...
```
**Solução**: Verifique se o MinIO está rodando e as credenciais no `.env` estão corretas.

### Imagens não aparecem no frontend
1. Verifique se as URLs estão corretas no banco de dados
2. Certifique-se que `next.config.mjs` permite o domínio do MinIO
3. Abra `http://localhost:9000/net.public.p2p/skins/AK-47/Redline.svg` no navegador

### Bucket com permissões incorretas
Execute:
```bash
cd backend
go run cmd/migrate-images/main.go
```
Isso recria o bucket com as permissões corretas.

## 📝 Notas

- As imagens são SVG (vetoriais), então sempre ficam nítidas em qualquer tamanho
- O sistema detecta automaticamente se uma imagem já existe antes de criar
- Cada skin tem cores únicas baseadas na raridade para fácil identificação visual
