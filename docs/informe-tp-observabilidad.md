# Informe Técnico — TP Observabilidad y Confiabilidad

> **Documento unificado de entrega y guion de presentación.**
> Cada sección indica su responsable. Las secciones marcadas como `PENDIENTE` deben completarse con capturas y detalles de la implementación.

---

## Equipo

| Integrante | Rol en el TP |
|---|---|
| Josias | Infraestructura (Harbor, ArgoCD, Argo Rollouts), OTel, Jaeger, Loki |
| Maxi | Logs estructurados y correlación del frontend |
| Jero | Trazas distribuidas (OpenTelemetry) |
| Lautaro | Saturación (cAdvisor) y laboratorio de simulación |
| Gabriel | Grafana como código, dashboard Golden Signals, documentación y postmortem |

---

## 0. Resumen Ejecutivo

FinFlow es una fintech de billeteras y pagos desplegada en un clúster local Kubernetes (K3D) gestionado por GitOps (ArgoCD). Sobre ese sistema implementamos las tres pilares de la observabilidad —**métricas, logs y trazas**— con Prometheus/Grafana, Loki y Jaeger, más un dashboard con las 4 Golden Signals de Google SRE. Como ejercicio adicional, diseñamos un **laboratorio de inyección de fallos** (`/lab`) y documentamos un **postmortem blameless (RCA)** de un incidente simulado.

---

## 1. El Sistema: FinFlow

- **Qué es:** fintech con billeteras (`/api/wallets`) y transferencias (`/api/payments`).
- **Servicios:** frontend (React) + backend (Go/Gin) + PostgreSQL + Redis + Unleash (feature flags).
- **Comunicación entre servicios:** el frontend consume la API REST del backend; el backend persiste en Postgres y cachea en Redis.

> **PENDIENTE — Diagrama de arquitectura** (generar con archify o adjuntar captura).

---

## 2. Despliegue y GitOps

> **PENDIENTE — Responsable: Josias**

Puntos a cubrir:
- Clúster local K3D y cómo se levanta (`make up`).
- ArgoCD app-of-apps (`root-local` → aplicaciones por servicio).
- **Argo Rollouts** con estrategia canary (`setWeight: 33` + `pause`) y promoción manual.
- Harbor privado vía Tailscale y el runner local de CI/CD (build multi-arquitectura).
- *Evidencia:* captura de la UI de ArgoCD / Argo Rollouts con el canary pausado.

---

## 3. Métricas y las 4 Golden Signals

> **Responsable: Gabriel** (dashboard) + **Lautaro** (saturación con cAdvisor)

### 3.1 Recolección de métricas
- **Prometheus** scrapea el backend Go (métricas de Gin: `gin_requests_total`, `gin_request_duration_seconds_bucket`).
- **cAdvisor** (kubelet) expone métricas reales de contenedor: `container_cpu_usage_seconds_total`, `container_memory_working_set_bytes`, `container_spec_cpu_quota`, etc.

### 3.2 Grafana como código (IaC)
- Los **DataSources** y **dashboards** se provisionan por ConfigMaps montados en Grafana (no se tocan a mano).
- Los dashboards se versionan en Git y se despliegan con ArgoCD.

### 3.3 Las 4 Golden Signals

| Señal | Métrica | Query PromQL |
|---|---|---|
| **Tráfico** | `gin_requests_total` | `sum(rate(gin_requests_total{namespace=~"$namespace"}[1m])) by (code)` |
| **Latencia** | `gin_request_duration_seconds_bucket` | `histogram_quantile(0.95, sum(rate(...{code=~"[23].."}[5m])) by (le))` |
| **Errores** | `gin_requests_total` | `((sum(rate(...{code=~"5.."}[1m])) or vector(0)) / ...) * 100` |
| **Saturación** | cAdvisor | `rate(container_cpu_usage_seconds_total[2m]) / (container_spec_cpu_quota / container_spec_cpu_period)` |

> **Nota SRE:** la latencia se mide solo sobre peticiones **exitosas** (`code=~"[23].."`) para no falsear el percentil con errores rápidos de validación.

*Evidencia:* captura del dashboard `FinFlow - 4 Golden Signals`.

---

## 4. Logs estructurados y correlación

> **PENDIENTE — Responsable: Maxi** (frontend) + **Josias** (backend/Loki)

Puntos a cubrir:
- Logs JSON estructurados en el backend (`log/slog`) con `request_id` y `latency_ms`.
- Logs JSON en el frontend (misma clave `msg: "http_request"`).
- **Correlación:** el frontend genera el `request_id` → lo envía por header `X-Request-ID` → el backend lo **reusa**.
- Loki centraliza los logs de todos los pods (vía OTel collector agent).
- *Evidencia:* captura de Loki mostrando logs de front y back con el mismo `request_id`.

---

## 5. Trazas distribuidas con OpenTelemetry

> **PENDIENTE — Responsable: Jero**

Puntos a cubrir:
- SDK OpenTelemetry en el backend Go (`internal/telemetry/tracer.go`).
- Web SDK en el frontend (`WebTracerProvider` + `FetchInstrumentation` con propagación de `traceparent`).
- Collector OTel (gateway) → Jaeger.
- **Traza distribuida de punta a punta** (frontend → backend).
- Correlación métricas/logs/trazas por `trace_id` (Data Links de Loki → Jaeger).
- *Evidencia:* captura de Jaeger con un trace completo.

---

## 6. Incidente simulado y postmortem (RCA)

> **Responsable: Gabriel** (documentación) + **Lautaro** (laboratorio `/lab`)

### 6.1 Laboratorio de inyección de fallos
- SPA en `/lab` con una página por señal (`/lab/traffic`, `/lab/latency`, `/lab/errors`, `/lab/saturation`).
- Endpoints de inyección controlada (con topes de seguridad): `/api/lab/cpu`, `/api/lab/slow`, `/api/lab/error`, `/api/lab/burst`, `/api/lab/memory`.

### 6.2 Simulacro de incidente
- Escenario A: inyección por el laboratorio (reproducible).
- Escenario B: caída de una dependencia real (PostgreSQL) para el RCA sistémico.
- Procedimiento completo en `docs/instructivo-simulacro-incidente.md`.

### 6.3 Postmortem blameless
- Estructura: resumen, impacto, timeline, **5 Whys**, Error Budget, action items.
- Documento completo en `docs/postmortem-incidente-simulado.md`.

*Evidencia:* capturas del dashboard durante el incidente (pico de errores/latencia/saturación) + fragmento de log con `request_id`.

---

## 7. Demo en vivo (guion de defensa)

> **PENDIENTE — Todos**

Flujo sugerido (≈15 min):

```text
1. Contexto (2 min) — qué es FinFlow y cómo está desplegado.
2. Despliegue GitOps (2 min) — ArgoCD + canary.
3. Métricas → Golden Signals (3 min) — dashboard.
4. Logs → correlación request_id (2 min) — Loki.
5. Trazas → traza distribuida (2 min) — Jaeger.
6. Incidente → postmortem (3 min) — lab + RCA.
7. Cierre (1 min) — lecciones aprendidas.
```

---

## Checklist de entrega

- [ ] Diagrama de arquitectura (Parte 1)
- [ ] Sección de despliegue (Parte 2 — Josias)
- [ ] Sección de logs (Parte 4 — Maxi)
- [ ] Sección de trazas (Parte 5 — Jero)
- [ ] Guion de demo (Parte 7 — Todos)
- [ ] Revisión final y unificación de docs en carpeta `docs/`
