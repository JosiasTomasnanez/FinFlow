# FinFlow: plataforma GitOps sobre Kubernetes (K3s y Amazon EKS)

Proyecto integrador de la **Diplomatura en DevOps (FCEFyN, UNC)**. FinFlow es una billetera digital ficticia (transferencias P2P y pagos) que usamos como carga de trabajo para construir una plataforma completa: empaquetado con Helm, entrega por GitOps con Argo CD, despliegues canary, autoescalado por métricas, observabilidad y gestión segura de secretos. Todo corre en un clúster **K3s local** y en **Amazon EKS**.

![Arquitectura de producción](Diagramas/Diagrama%20arquitectura%20produccion.jpg)

## Qué hay implementado

| Área | Implementación | Dónde mirar |
|---|---|---|
| Aplicación | API en **Go (Gin)** + SPA en **React/Vite**, **PostgreSQL** y **Redis**. Monolito modular. | `cmd/`, `internal/`, `frontend/` |
| Empaquetado | Chart de Helm propio para la app y otro para la infraestructura de la plataforma. Un archivo `values` por entorno (staging/prod) y por destino (local/AWS). | `finflow-chart/`, `finflow-infra/`, `values-*.yaml` |
| GitOps | Argo CD con patrón **app-of-apps**: `root-local.yaml` y `root-aws.yaml` sincronizan app, infraestructura, KEDA, Loki, Jaeger, OpenTelemetry Collector, Harbor y Sealed Secrets. | `argocd-infrastructure/` |
| Despliegue progresivo | **Argo Rollouts** canary: 33% del tráfico a la nueva versión y pausa para promoción manual. | `finflow-chart/templates/backend.yaml`, `values-prod*.yaml` |
| Feature flags | **Unleash** self-hosted con SDK de Go. Permite liberar código apagado y activarlo por segmentos. | `finflow-infra/templates/unleash.yaml` |
| Autoescalado | **KEDA** `ScaledObject` con trigger de Prometheus sobre `rate(gin_requests_total[2m])`. | `finflow-chart/templates/scaledobject.yaml`, `values-keda.yaml` |
| Observabilidad | **Prometheus + Grafana** con dashboard de Golden Signals versionado como código, **Loki** para logs JSON con `request_id`, **Jaeger** + **OpenTelemetry** para trazas del frontend al backend. | `finflow-infra/dashboards/`, `docs/informe-tp-observabilidad.md` |
| Secretos | **Sealed Secrets**: solo se commitean secretos cifrados con la clave pública del controlador. `scripts/seal-secrets.sh` los genera a partir de archivos `.env` locales que nunca se versionan. | `finflow-chart/templates/secrets/`, `scripts/seal-secrets.sh` |
| Red | `NetworkPolicy` por entorno para aislar base de datos y Redis. | `finflow-chart/templates/networkpolicy.yaml` |
| Registry | **Harbor** privado (local, accesible por Tailscale) y **Amazon ECR** (AWS). | `argocd-infrastructure/argocd-config/app-harbor.yaml` |
| AWS | EKS con **AWS Load Balancer Controller (ALB)**, **EBS CSI** (`gp3`) para volúmenes persistentes y ECR. Política IAM del controller incluida. | `values-*-aws.yaml`, `storageclass.yaml`, `iam_policy.json` |
| Calidad y seguridad | Análisis estático con **SonarCloud** en cada push. En la etapa de staging (backend en Render y frontend en Vercel) se ejecutó **OWASP ZAP baseline (DAST)**, que publicaba los hallazgos como issues ([#10](https://github.com/JosiasTomasnanez/FinFlow/issues/10) a [#14](https://github.com/JosiasTomasnanez/FinFlow/issues/14)). | `.github/workflows/sonar.yml` |

## Flujo de entrega (CI/CD + GitOps)

![Flujo GitOps](Diagramas/Diagrama%20gitops.jpg)

1. Push a `main`: GitHub Actions construye las imágenes de backend y frontend y las publica en **Harbor** (runner self-hosted) y en **ECR**.
2. El mismo pipeline actualiza el tag de imagen en `values-staging*.yaml` con `yq` y commitea el cambio. Git es la única fuente de verdad.
3. Argo CD detecta el commit y sincroniza staging.
4. Tag `vX.Y.Z`: se actualizan los `values-prod*.yaml` y Argo Rollouts despliega la nueva versión en canary al 33% hasta la promoción manual.

> **Estado actual:** el clúster de EKS está apagado para no generar costos. Por eso el job `build-and-push-ecr` está desactivado (`if: false`) en `.github/workflows/deploy.yml`. Los manifiestos y values de AWS siguen en el repo y alcanza con reactivar ese job y aplicar `root-aws.yaml` sobre un clúster EKS nuevo.

## Levantar el entorno local (K3s con k3d)

Requisitos: `docker`, `k3d`, `kubectl`, `tailscale` y `envsubst`. `make doctor` verifica qué falta.

```bash
make up      # crea el clúster k3d, instala Argo CD y sincroniza todas las apps
make status  # estado de pods y aplicaciones de Argo CD
make down    # elimina el clúster
```

El paso a paso y el debugging están en [`docs/MAKE.md`](docs/MAKE.md).

### Solo la aplicación con Docker Compose

```bash
cp .env.example .env   # completar con valores locales
docker compose up -d --build
```

Nginx expone todo en `http://localhost`: `/` (frontend), `/api/*` (backend), `/grafana` y `/unleash`.

## Documentación

- [`docs/Construccion_y_orquestacion_de_contenedores.md`](docs/Construccion_y_orquestacion_de_contenedores.md): decisiones de contenedores, Helm y K3s/EKS.
- [`docs/informe-tp-observabilidad.md`](docs/informe-tp-observabilidad.md): métricas, logs y trazas, con capturas en `docs/capturas/`.
- [`docs/lab-golden-signals.md`](docs/lab-golden-signals.md): laboratorio de saturación y Golden Signals.
- [`Diagramas/`](Diagramas/): arquitectura de staging, producción, entorno de infraestructura y GitOps (`.drawio` + `.jpg`).

## Equipo

Proyecto grupal de la Diplomatura en DevOps (UNC): **Josias Tomas Ñañez**, Jerónimo Massaro, Maximiliano Cravero, Gabriel Oliva, Lautaro Castro y Marcos Reynoso, quien hizo el primer workflow de OWASP ZAP.

Josias diseñó e implementó la mayor parte de la plataforma: charts de Helm, estructura de Argo CD, pipelines de GitHub Actions, configuración de AWS y Sealed Secrets.
