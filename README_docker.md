# Quick start (Docker Compose)

This project ships a dev-friendly Compose stack: Postgres + Redis + API + Web (Vite dev server).

## Prerequisites

- Docker and Docker Compose

## Run (dev stack: live code, hot reload)

```bash
# From repo root
docker compose up --build
```

Services/ports (dev):

- API: http://localhost:8080 (Gin)
- Web: http://localhost:5173 (Vite dev, proxies directly to API via base URL in axios)
- Postgres: localhost:5432 (user `postgres`, password `password`, db `pp_db`)
- Redis: localhost:6379

Data volumes:

- Postgres data: `pgdata`
- Redis data: `redisdata`

## Run (prod-ish stack: built binaries/static + nginx proxy)

```bash
docker compose -f docker-compose.prod.yml up --build
```

Services/ports (prod-ish):
- API: http://localhost:8080 (built binary)
- Web: http://localhost:5173 (nginx static + `/api/` proxy to API + `/healthz`)
- Postgres: localhost:5432
- Redis: localhost:6379

## Configuration

Config resolution:
- Default file: `config.yaml`
- Override via env `CONFIG_FILE` (Compose already sets `config.docker` for service hostnames)
- Override any field with env, e.g. `DATABASE_HOST=127.0.0.1`, `SERVER_PORT=8080`, `TRADING_MAXOPENORDERS=5`

Axios/Vite base URL can be overridden with `VITE_API_BASE`; default is `http://localhost:8080/api/v1`.
For prod build container, set build arg `API_BASE` (Compose prod already sets `http://localhost:8080/api/v1`).

## TLS / certificates

- 快速启用 TLS（自备证书）：使用 `docker-compose.prod.tls.yml`，在 repo 根放置 `./certs/fullchain.pem` 和 `./certs/privkey.pem`（PEM），会自动挂载到 web 容器并使用 `web/nginx.tls.conf`，监听 443（同时保留 80/5173）。
- 如需自定义 API 基址为 https，构建 web 时调整 `API_BASE`（TLS compose 已设为 `https://localhost/api/v1`）。
- 自动签发示例：将本栈置于 Traefik/Caddy 后面，由外层代理处理 ACME。参考 Traefik 片段：
  ```yaml
  web:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.pp.rule=Host(`pp.example.com`)"
      - "traefik.http.routers.pp.entrypoints=websecure"
      - "traefik.http.routers.pp.tls.certresolver=letsencrypt"
      - "traefik.http.services.pp.loadbalancer.server.port=80"
  ```
  在这种场景下，可使用普通 `docker-compose.prod.yml`，由 Traefik 终结 TLS。

## Useful commands

- Stop: `docker compose down`
- Stop and clean volumes: `docker compose down -v`
- Rebuild a single service: `docker compose build api` (or `web`)
