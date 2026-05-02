# Cryptomines Online

Juego online multiplayer compuesto por tres servicios:

- **`backend/`** — API HTTP en Go 1.25 (puerto `8080`).
- **`frontend/`** — React 19 + Vite 7 + TypeScript (puerto `5173`).
- **`supabase/`** — Postgres + Auth local vía Supabase CLI (API `54321`, DB `54322`).

## Requisitos

- [Go](https://go.dev/dl/) 1.25+
- [Node.js](https://nodejs.org/) 20+ y [Yarn](https://classic.yarnpkg.com/)
- [Supabase CLI](https://supabase.com/docs/guides/local-development/cli/getting-started)
- [Docker](https://docs.docker.com/get-docker/) (lo usa Supabase CLI)
- (Opcional) [Air](https://github.com/air-verse/air) para hot-reload del backend: `go install github.com/air-verse/air@latest`

## Puesta en marcha

Abre **tres terminales**, una por servicio.

### 1. Supabase (DB + Auth)

```bash
cd supabase
supabase start          # arranca Postgres, Auth, Studio, etc.
supabase db reset       # aplica todas las migraciones de supabase/migrations/
```

Tras `supabase start` la CLI imprime las claves locales (`anon`, `service_role`). Las del `.env.example` ya coinciden con la demo, normalmente no hace falta cambiarlas.

Studio queda disponible en http://127.0.0.1:54323.

### 2. Backend (Go)

```bash
cd backend
cp .env.example .env    # solo la primera vez
go mod download
go run ./cmd/server     # o `air` si quieres hot-reload
```

API disponible en http://localhost:8080. Healthcheck: `GET /api/health`.

### 3. Frontend (React + Vite)

```bash
cd frontend
yarn install            # solo la primera vez
yarn dev
```

UI disponible en http://localhost:5173.

## Variables de entorno (backend)

`backend/.env` (ver `backend/.env.example`):

| Variable | Default | Descripción |
| --- | --- | --- |
| `SUPABASE_URL` | `http://127.0.0.1:54321` | URL de la API local de Supabase |
| `SUPABASE_ANON_KEY` | demo key | Clave pública anon |
| `SUPABASE_SERVICE_ROLE_KEY` | demo key | Clave service-role (admin) |
| `SUPABASE_DB_URL` | `postgresql://postgres:postgres@127.0.0.1:54322/postgres?sslmode=disable` | Conexión Postgres |
| `JWT_SECRET` | demo secret | Debe coincidir con el JWT secret de Supabase |
| `PORT` | `8080` | Puerto del backend |
| `FRONTEND_URL` | `http://localhost:5173` | Origin permitido por CORS |

## Comandos útiles

**Backend**

```bash
go build -o ./tmp/main ./cmd/server     # build
go test ./...                           # tests
```

**Frontend**

```bash
yarn dev        # dev server
yarn build      # build de producción
yarn lint       # eslint
```

**Supabase**

```bash
supabase status                         # ver puertos y claves
supabase db reset                       # rehace la DB y reaplica migraciones
supabase migration new <nombre>         # nueva migración
supabase stop                           # detiene los contenedores
```

## Orden recomendado al arrancar de cero

1. `supabase start` y `supabase db reset` en `supabase/`.
2. `go run ./cmd/server` en `backend/`.
3. `yarn dev` en `frontend/`.
4. Abrir http://localhost:5173.
