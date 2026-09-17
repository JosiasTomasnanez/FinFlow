# Lab Golden Signals — FinFlow

El laboratorio vive en **`/lab`**, separado de la caja (`/`). Cada señal tiene su ruta: `/lab/traffic`, `/lab/latency`, `/lab/errors`, `/lab/saturation`.

Los pagos reales no sirven: esperan Redis/Postgres (CPU/throttle/mem planos) y no fabrican 5xx ni latencia de handler.

Código:

- Trabajo sintético: [`internal/lab/`](../internal/lab/) (`Spin`, `Hold`, topes).
- HTTP: [`internal/api/lab_*.go`](../internal/api/lab_routes.go) (un archivo por señal).
- UI: [`frontend/src/lab/`](../frontend/src/lab/) (`LabLayout` + una página por señal).

`GET /api/flags` incluye `feature_lab: true`. Si viene `false`, `/lab` redirige a `/`.

## Qué botón pega a qué

| Señal | UI | Endpoint | Grafana |
|---|---|---|---|
| Tráfico | `/lab/traffic` | `GET /api/lab/ok` | RPS (`gin_requests_total`, `[1m]`) |
| Latencia | `/lab/latency` | `POST /api/lab/slow` `{ "delay_ms" }` → **200** | p95 exitosas (`code=~"[23].."`) |
| Errores | `/lab/errors` | `POST /api/lab/error` `{ "code": 500 }` | Tasa 5xx y barras 4xx/5xx |
| CPU | `/lab/saturation` | `POST /api/lab/cpu` `{ "duration_ms" }` loop (no sleep) | CPU uso/límite `[2m]` |
| Throttle | `/lab/saturation` | N × `POST /api/lab/burst` `{ "spin_ms": 80 }` | Throttle CFS |
| Memoria | `/lab/saturation` | `POST /api/lab/memory` y `/memory/release` | Working set/límite (hold **global**) |

Topes en server: delay 5 s, CPU 10 s, burst 200 ms, RAM 128 MiB, errores solo 500/502/503.

## Efectos cruzados

Una CPU de 5 s también suma 1 request lenta al RPS y al p95. El resto de botones está pensado para no mezclar: `ok` es barato, `slow` no quema CPU, `error` no duerme, `memory` no spinea.

Esperá ~1–2 minutos en saturación (`rate[2m]`). Tráfico y errores se ven antes (`[1m]`).

Staging con `limits.cpu: 250m` / `256Mi` hace visibles CPU y throttle. Sin limits, saturación no tiene denominador.

Más contexto de las tres métricas de techo: [saturacion-golden-signals.md](./saturacion-golden-signals.md).
