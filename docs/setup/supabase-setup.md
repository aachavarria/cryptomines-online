# Supabase Local Development Setup

## Prerequisites

- **Docker Desktop** (running)
- **Supabase CLI** v2.54+ (`brew install supabase/tap/supabase`)

## Quick Start

### Start Supabase

```bash
cd /Users/yurei/cryptomines-online
supabase start
```

First run will pull Docker images (may take a few minutes).

### Stop Supabase

```bash
supabase stop
```

To stop and reset all data:

```bash
supabase stop --no-backup
```

### Check Status

```bash
supabase status
```

For environment variable format:

```bash
supabase status -o env
```

## Service URLs and Ports

| Service       | URL / Connection String                                      | Port  |
|---------------|--------------------------------------------------------------|-------|
| API (REST)    | http://127.0.0.1:54321                                       | 54321 |
| GraphQL       | http://127.0.0.1:54321/graphql/v1                            | 54321 |
| Database      | postgresql://postgres:postgres@127.0.0.1:54322/postgres      | 54322 |
| Studio (UI)   | http://127.0.0.1:54323                                       | 54323 |
| Mailpit       | http://127.0.0.1:54324                                       | 54324 |
| Analytics     | http://127.0.0.1:54327                                       | 54327 |
| S3 Storage    | http://127.0.0.1:54321/storage/v1/s3                         | 54321 |

## Keys (Local Development Only)

These keys are deterministic for local development and are safe to commit.

### ANON_KEY (public, client-side)

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0
```

### SERVICE_ROLE_KEY (secret, server-side only)

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImV4cCI6MTk4MzgxMjk5Nn0.EGIM96RAZx35lJzdJsyH-qQwv8Hdp7fsn3W0YpN81IU
```

### JWT_SECRET

```
super-secret-jwt-token-with-at-least-32-characters-long
```

### Publishable Key

```
sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH
```

### Secret Key

```
sb_secret_N7UND0UgjKTVK-Uodkm0Hg_xSvEMPvz
```

## Configuration

The Supabase config lives at `supabase/config.toml`. Key settings:

- **Anonymous sign-ins**: Enabled (`enable_anonymous_sign_ins = true`)
- **API port**: 54321
- **DB port**: 54322
- **Studio port**: 54323
- **Project ID**: cryptomines-online

## Database Migrations

Migrations are stored in `supabase/migrations/`.

### Create a new migration

```bash
supabase migration new <migration_name>
```

### Apply migrations (reset local DB)

```bash
supabase db reset
```

### Diff schema changes

```bash
supabase db diff --use-migra -f <migration_name>
```

## Connecting from Go Backend

Use these environment variables in your backend:

```env
SUPABASE_URL=http://127.0.0.1:54321
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0
SUPABASE_SERVICE_ROLE_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImV4cCI6MTk4MzgxMjk5Nn0.EGIM96RAZx35lJzdJsyH-qQwv8Hdp7fsn3W0YpN81IU
DATABASE_URL=postgresql://postgres:postgres@127.0.0.1:54322/postgres
```

## Connecting from React Frontend

```typescript
import { createClient } from '@supabase/supabase-js'

const supabase = createClient(
  'http://127.0.0.1:54321',
  'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0'
)
```
