# Logs estructurados + request_id (Backend Go) — FinFlow

## Qué hace este patch

1. **Logs en JSON** en vez de texto plano (`log/slog` de la librería estándar
   de Go, sin agregar dependencias nuevas).
2. Cada request HTTP recibe un **request_id**: si ya venía en el header
   `X-Request-ID` lo reusa, si no lo genera. Se guarda en el `context.Context`
   de Go y se propaga a la capa de handlers y de services, así que todos los
   logs emitidos mientras se procesa ese request llevan el mismo `request_id`.
3. El `request_id` se devuelve en la respuesta HTTP (header `X-Request-ID`),
   para que otros servicios (el frontend) puedan loguear el mismo id.
4. Se agregó una línea de "access log" por request (método, path, status,
   latencia) y algunas líneas de negocio (`wallet_created`,
   `payment_completed`, `login_failed`, hit/miss de Redis, etc.), todas con
   el `request_id` adentro.

## Archivos

- `internal/logging/logging.go` **(nuevo)** — el logger central y los
  helpers `WithRequestID` / `FromContext`.
- `internal/api/middleware.go` **(nuevo)** — el middleware de Gin que genera
  o propaga el `request_id` y loguea cada request.
- `internal/api/server.go` — se cambió `gin.Default()` por `gin.New()` +
  `gin.Recovery()` + el middleware nuevo (el logger de texto de Gin quedaba
  duplicado con el nuestro).
- `internal/api/handlers.go`, `internal/api/auth.go` — ahora pasan
  `c.Request.Context()` a los services y loguean los eventos importantes.
- `internal/service/wallet.go`, `internal/service/auth.go` — los métodos
  ahora reciben `ctx context.Context` como primer parámetro, para poder
  loguear con el `request_id` del request que los originó. Los
  `fmt.Printf` con emojis se reemplazaron por logs estructurados.
- `internal/service/wallet_test.go` — **ojo:** este archivo ya no compilaba
  contra el código actual (llamaba a `NewWalletService(store)` con un solo
  argumento). Lo dejé arreglado y actualizado a las nuevas firmas con `ctx`.

No toqué `go.mod` ni `cmd/finflow/main.go` (los `log.Printf` de arranque de
Postgres/Redis/Unleash quedaron como estaban — se pueden migrar después,
no son parte de la correlación por request).

## Ejemplo de log generado

```json
{"time":"2026-09-15T02:10:03Z","level":"INFO","msg":"http_request","service":"finflow-backend","request_id":"a1b2c3...","method":"POST","path":"/api/payments","status":201,"latency_ms":4}
{"time":"2026-09-15T02:10:03Z","level":"INFO","msg":"payment_completed","service":"finflow-backend","request_id":"a1b2c3...","from_wallet_id":"w1","to_wallet_id":"w2","amount":300}
```

## Coordinación con el resto del equipo

- **Con Maxi (logs en el frontend):** el contrato es el header
  `X-Request-ID`. Si el frontend manda ese header en cada request al
  backend, el backend lo reusa (no genera uno nuevo) y así un mismo
  `request_id` aparece tanto en los logs del frontend como en los del
  backend para la misma acción del usuario. Si el frontend no lo manda,
  igual funciona: el backend genera uno y lo devuelve en la respuesta, por
  si el frontend lo quiere tomar de ahí para loguear.

- **Con Jero (OpenTelemetry / trazas):** cuando él instrumente el backend
  con OTel, lo más prolijo es que el `trace_id` que genera OTel reemplace
  (o conviva con) este `request_id` — la consigna pide "request_id/trace_id"
  como si fueran intercambiables. Como referencia: `logging.FromContext(ctx)`
  es el único lugar que decide qué id usar, así que cuando esté el
  middleware de OTel alcanza con leer el `trace_id` del span activo ahí
  adentro en vez del `X-Request-ID`, sin tener que tocar handlers ni
  services de nuevo. Vale la pena que se sincronicen para no terminar con
  dos ids distintos por request.

## Cómo aplicarlo

Copiá estos archivos a las mismas rutas dentro de `FinFlow/`, respetando la
carpeta `internal/`. Son ediciones completas de archivo (no diffs), así que
simplemente reemplazan a los actuales. Después:

```bash
go build ./...
go test ./...
```

(Validé la lógica de estos archivos de forma aislada porque el `go.mod` del
proyecto pide Go 1.25 y acá solo tenía 1.22 disponible para probar — no
alcancé a correr el build completo del módulo, así que antes de dar por
cerrado esto convendría correr `go build ./... && go test ./...` con la
toolchain real.)
