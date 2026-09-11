# Uso del Makefile — FinFlow

Este documento explica **solo** el `Makefile` de la raíz del repo: qué hace, en qué orden, qué requiere de antemano y cómo debuguear cuando algo no arranca.

## 1. Requisitos previos (antes de tocar `make`)

| Requisito | Cómo se obtiene |
|---|---|
| `docker`, `k3d`, `kubectl`, `tailscale`, `envsubst` instalados | `make doctor` te dice cuál falta |
| Cuenta en la tailnet de Tailscale de Josias | Invitación de Josias (no lo puede automatizar el Makefile) |
| `secrets/sealed-secrets-master-key.yaml` | Josias te la pasa por un canal seguro (nunca va en el repo) |
| `registries.yaml` en la raíz del repo | Ya viene versionado en git |

Corré `make doctor` antes de `make up` la primera vez — te va a listar exactamente qué falta en vez de que `make up` explote a mitad de camino.

## 2. Comando único: `make up`

```bash
make up
```

Es lo único que necesitás correr. Internamente encadena 11 pasos, en este orden exacto (podés correr cada uno suelto si necesitás repetir solo uno):

| # | Target | Qué hace |
|---|---|---|
| 1 | `doctor` | Chequea binarios instalados, Docker corriendo, `registries.yaml` y la clave de sealed-secrets |
| 2 | `resolve-ip` | Resuelve la IP de Tailscale del host que corre Harbor y la guarda en `.harbor-ip` |
| 3 | `render-config` | Genera `k3d-config.yaml` desde el template, inyectando esa IP |
| 4 | `cluster-up` | Crea el cluster k3d (si ya existe, no hace nada) |
| 5 | `wait-cluster` | Espera a que el nodo esté `Ready` |
| 6 | `restore-sealed-secrets-key` | Restaura la clave privada de sealed-secrets (para que las `SealedSecret` del repo se puedan desencriptar) |
| 7 | `bootstrap-argocd` | Instala Argo CD (manifiesto oficial) |
| 8 | `wait-argocd` | Espera a que `argocd-server` esté disponible |
| 9 | `apply-argocd-config` | Aplica `argocd-infrastructure/argocd-config/` por `kubectl` (ver sección 3) |
| 10 | `apply-root` | Aplica `argocd-infrastructure/root-local.yaml` (ver sección 3) |
| 11 | `wait-apps` | Espera a que **todas** las Applications de Argo CD queden `Synced` + `Healthy` |

## 3. Qué levanta cada capa (importante, no es intuitivo)

Hay **dos árboles de Applications separados**, con distinto mecanismo de arranque:

### `argocd-config/` — bootstrap manual, corre en el paso 9
Se aplica directo con `kubectl apply` (Argo CD todavía no puede auto-gestionarlo porque recién se está instalando). Contiene:

- **`sealed-secrets`** → el operador de sealed-secrets, en `kube-system`
- **`harbor-infra`** → Harbor completo (registry privado), en el namespace `harbor`
- **`argocd-self-config`** → parchea `argocd-cm` (los `ignoreDifferences` para Rollouts/Deployments)

El Makefile espera explícitamente a que el CRD `SealedSecret` quede registrado y el controller esté disponible **antes** de seguir, porque `finflow-prod`/`finflow-staging` (paso siguiente) dependen de eso para desencriptar sus secretos.

### `apps-local/` — gestionado por Argo CD, corre en el paso 10
`root-local.yaml` es una Application "app-of-apps": apunta a `argocd-infrastructure/apps-local` con `recurse: true`, así que Argo CD crea automáticamente estas 4 Applications:

- **`finflow-infra`** → namespace `finflow-infra`, chart de infra propio del proyecto
- **`finflow-prod`** → namespace `finflow-prod`
- **`finflow-staging`** → namespace `finflow-staging`
- **`keda-operator`** → KEDA (chart oficial), también en `finflow-infra`

> **Ojo:** "infra" aparece en los dos árboles pero significa cosas distintas. `argocd-config`'s Harbor/sealed-secrets son infraestructura de *bootstrap del cluster*. `apps-local/finflow-infra` es infraestructura *de la app* (lo que antes eran `app-infra.yaml`, `app-prod.yaml`, etc.). `root-local` solo trae la segunda.

### Consecuencia práctica
Si `finflow-prod`/`finflow-staging` quedan un rato en `ImagePullBackOff` recién arrancado el cluster, es esperable: Harbor tarda en levantar todos sus componentes (core, portal, jobservice, DB, redis, trivy). Con `selfHeal: true` se recupera solo apenas Harbor esté listo. `make wait-apps` ya espera esto por vos.

## 4. Comandos de diagnóstico

| Comando | Para qué |
|---|---|
| `make status` | Lista rápida de todas las Applications y su estado |
| `make debug` | Nodos + pods que no están `Running`/`Completed` + Applications, todo junto |
| `make logs-failing` | Logs (o `describe` si no hay logs) de cada pod que no está `Ready` |
| `make check-harbor` | Prueba conectividad real a Harbor desde dentro del cluster vía la IP de Tailscale |

Si `make up` se cuelga en el paso 11 (`wait-apps`) con timeout, el cluster **queda arriba igual** — no aborta nada previo. Corré `make debug` o `make logs-failing` para ver qué está trabado.

## 5. Otros comandos

```bash
make down                          # borra el cluster completo
make clean                         # borra cluster + k3d-config.yaml + .harbor-ip
TAILSCALE_HOST=otro-host make up   # si tu host de Tailscale no es el default
```

## 6. Lo que el Makefile NO automatiza (pasos manuales, una sola vez por persona)

1. Que Josias te invite a su tailnet de Tailscale.
2. Que Josias te pase `secrets/sealed-secrets-master-key.yaml` por un canal seguro y la coloques en esa ruta antes de correr `make up`.

Sin esos dos pasos, `make doctor` te lo va a marcar como error/warning antes de que pierdas tiempo con el resto del pipeline.
