SHELL := /bin/bash

# ── Config ──────────────────────────────────────────────────────────────
CLUSTER_NAME       := finflow
TAILSCALE_HOST     ?= josias-latitude-e5450
SEALED_SECRETS_KEY := secrets/sealed-secrets-master-key.yaml
KUBECONFIG_FILE    := $(HOME)/.k3d/kubeconfig-$(CLUSTER_NAME).yaml

.PHONY: all up down clean doctor render-config resolve-ip cluster-up wait-cluster \
        restore-sealed-secrets-key bootstrap-argocd wait-argocd apply-argocd-config apply-root \
        wait-apps status debug logs-failing check-harbor

all: up

## Pipeline completo: chequeo de dependencias + cluster + clave de sealed-secrets + Argo CD + argocd-config + app-of-apps
up: doctor resolve-ip render-config cluster-up wait-cluster restore-sealed-secrets-key bootstrap-argocd wait-argocd apply-argocd-config apply-root wait-apps
	@echo ""
	@echo "✅ Cluster arriba y Argo CD sincronizado."
	@echo "   Ver estado en cualquier momento con: make status"

## 0. Verifica que las herramientas necesarias estén instaladas ANTES de tocar nada.
##    No reemplaza instalarlas: solo evita que 'make up' falle a mitad de camino
##    con un error críptico.
doctor:
	@ok=1; \
	for bin in docker k3d kubectl tailscale envsubst; do \
		if ! command -v $$bin >/dev/null 2>&1; then \
			echo "❌ Falta '$$bin' en el PATH. Instalalo antes de seguir."; ok=0; \
		fi; \
	done; \
	if command -v docker >/dev/null 2>&1 && ! docker info >/dev/null 2>&1; then \
		echo "❌ Docker está instalado pero no responde (¿está corriendo el daemon? ¿tenés permisos?)."; ok=0; \
	fi; \
	if [ ! -f registries.yaml ]; then \
		echo "❌ Falta registries.yaml en la raíz del repo (viene versionado en git, revisá el clone)."; ok=0; \
	fi; \
	if [ ! -f $(SEALED_SECRETS_KEY) ]; then \
		echo "⚠️  Falta $(SEALED_SECRETS_KEY) (no viaja en el repo por seguridad)."; \
		echo "    Pedísela a Josias por un canal seguro y colocala en esa ruta antes de correr 'make up'."; \
	fi; \
	if [ $$ok -eq 1 ]; then \
		echo "✅ Dependencias básicas OK."; \
	else \
		echo ""; echo "🛑 Resolvé lo anterior antes de correr 'make up'."; exit 1; \
	fi

## 1. Resuelve la IP de Tailscale del host que corre Harbor
resolve-ip:
	@command -v tailscale >/dev/null 2>&1 || { echo "❌ No encontré el cliente 'tailscale' instalado."; exit 1; }
	@IP=$$(tailscale ip -4 $(TAILSCALE_HOST) 2>/dev/null); \
	if [ -z "$$IP" ]; then \
		echo "❌ No pude resolver la IP de Tailscale para '$(TAILSCALE_HOST)'."; \
		echo "   Corré 'tailscale status' y ajustá TAILSCALE_HOST=<nombre-correcto> al llamar make."; \
		exit 1; \
	fi; \
	echo "$$IP" > .harbor-ip; \
	echo "🔎 IP de Tailscale de Harbor: $$IP"

## 2. Genera k3d-config.yaml a partir del template (este archivo NO se commitea)
render-config:
	@if [ ! -f registries.yaml ]; then \
		echo "❌ No encuentro registries.yaml en la raíz del repo (necesario para que el nodo confíe en Harbor)."; \
		exit 1; \
	fi
	@HARBOR_TAILSCALE_IP=$$(cat .harbor-ip) envsubst < k3d-config.yaml.tmpl > k3d-config.yaml
	@echo "📝 k3d-config.yaml generado."

## 3. Levanta el cluster k3d (si ya existe, no falla)
cluster-up:
	@if k3d cluster list | grep -q "^$(CLUSTER_NAME) "; then \
		echo "ℹ️  El cluster '$(CLUSTER_NAME)' ya existe, lo dejo como está."; \
	else \
		k3d cluster create --config k3d-config.yaml; \
	fi

## 4. Espera a que el nodo esté Ready
wait-cluster:
	@echo "⏳ Esperando a que el nodo esté Ready..."
	@kubectl wait --for=condition=Ready node --all --timeout=120s

## 5. Restaura la clave privada de sealed-secrets ANTES de que el controller
##    genere una propia. Si el controller ya está corriendo, lo reinicia
##    para que relea esta clave en vez de la que haya generado.
restore-sealed-secrets-key:
	@if [ ! -f $(SEALED_SECRETS_KEY) ]; then \
		echo "❌ No encuentro $(SEALED_SECRETS_KEY). Pedile a Josias que te la pase por un canal seguro."; \
		exit 1; \
	fi
	@kubectl create namespace kube-system --dry-run=client -o yaml | kubectl apply -f -
	@kubectl delete secret -n kube-system -l sealedsecrets.bitnami.com/sealed-secrets-key=active --ignore-not-found=true
	@kubectl apply -n kube-system -f $(SEALED_SECRETS_KEY)
	@echo "🔑 Clave de sealed-secrets restaurada en kube-system."
	@kubectl delete pod -n kube-system -l app.kubernetes.io/name=sealed-secrets --ignore-not-found=true

## 6. Bootstrap de Argo CD (instalación base, todavía sin apps).
##    OJO: esto asume que usás el manifiesto oficial. Si vos ya instalás
##    Argo CD de otra forma (helm, script propio, etc.) reemplazá esta receta.
bootstrap-argocd:
	@kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -n argocd --server-side --force-conflicts -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

## 7. Espera a que el server de Argo CD esté listo
wait-argocd:
	@echo "⏳ Esperando a que Argo CD esté listo..."
	@kubectl wait --for=condition=available deployment/argocd-server -n argocd --timeout=300s

## 8. Aplica manualmente lo que antes hacías a mano: sealed-secrets, harbor y
##    el self-management de Argo CD (argocd-config/). Deben ir ANTES que
##    apps-local, porque finflow-staging/prod dependen del CRD SealedSecret.
apply-argocd-config:
	@kubectl apply -n argocd -f argocd-infrastructure/argocd-config/
	@echo "⏳ Esperando a que el CRD de sealed-secrets quede registrado..."
	@until kubectl get crd sealedsecrets.bitnami.com >/dev/null 2>&1; do sleep 3; done
	@echo "⏳ Esperando a que el controller de sealed-secrets esté disponible..."
	@kubectl wait --for=condition=available deployment -n kube-system -l app.kubernetes.io/name=sealed-secrets --timeout=180s || true
	@echo "✅ argocd-config aplicado y sealed-secrets listo."

## 9. Aplica la app-of-apps raíz: de acá en más, Argo CD toma el control
##    de apps-local (finflow-infra, finflow-prod, finflow-staging, keda)
apply-root:
	@kubectl apply -f argocd-infrastructure/root-local.yaml -n argocd

## 10. Espera a que TODAS las Applications de Argo CD queden Synced + Healthy
##     (hasta acá 'up' solo aplicaba manifiestos, no confirmaba que hayan prendido).
wait-apps:
	@echo "⏳ Esperando a que Argo CD sincronice todas las Applications (puede tardar varios minutos)..."
	@for i in $$(seq 1 60); do \
		total=$$(kubectl get applications -n argocd --no-headers 2>/dev/null | wc -l | tr -d ' '); \
		ready=$$(kubectl get applications -n argocd -o jsonpath='{range .items[*]}{.status.sync.status}{" "}{.status.health.status}{"\n"}{end}' 2>/dev/null | grep -c "^Synced Healthy$$" || true); \
		if [ "$$total" -gt 0 ] && [ "$$ready" = "$$total" ]; then \
			echo "✅ Las $$total Applications están Synced + Healthy."; \
			exit 0; \
		fi; \
		echo "   ($$ready/$$total listas, intento $$i/60) reintentando en 10s..."; \
		sleep 10; \
	done; \
	echo "⚠️  Timeout esperando que todas las apps queden Healthy."; \
	echo "    Revisá con 'make debug' o 'make logs-failing'. El cluster sigue arriba igual."

## Ver el estado de las Applications de Argo CD
status:
	@kubectl get applications -n argocd

## Diagnóstico rápido: nodos, pods que no están Running/Completed, y Applications
debug:
	@echo "── Nodos ──"; kubectl get nodes -o wide
	@echo ""; echo "── Pods que NO están Running/Completed (todos los namespaces) ──"; \
	kubectl get pods -A --field-selector=status.phase!=Running,status.phase!=Succeeded 2>/dev/null || true
	@echo ""; echo "── Applications de Argo CD ──"; kubectl get applications -n argocd

## Logs (o describe si no hay logs) de cada pod que no está Ready. Útil para debug rápido en grupo.
logs-failing:
	@kubectl get pods -A --no-headers 2>/dev/null | awk '$$4 != "Running" && $$4 != "Completed" {print $$1, $$2}' | while read -r ns pod; do \
		echo "── $$ns/$$pod ──"; \
		kubectl logs -n "$$ns" "$$pod" --tail=30 --all-containers 2>/dev/null || kubectl describe pod -n "$$ns" "$$pod" | tail -20; \
		echo ""; \
	done

## Prueba que el cluster pueda resolver y hablar con Harbor vía la IP de Tailscale
## (corré 'make resolve-ip' antes si no lo hiciste en esta sesión)
check-harbor:
	@if [ ! -f .harbor-ip ]; then echo "❌ Corré 'make resolve-ip' primero."; exit 1; fi
	@kubectl run harbor-check --rm -i --restart=Never --image=curlimages/curl:latest --quiet -- \
		curl -sk -o /dev/null -w "HTTP %{http_code}\n" \
		--resolve harbor.finflow.local:443:$$(cat .harbor-ip) \
		https://harbor.finflow.local/v2/ \
	&& echo "✅ Harbor responde." || echo "❌ No se pudo contactar a Harbor (revisá Tailscale y el /etc/hosts del nodo k3d)."

## Borra el cluster completo
down:
	@k3d cluster delete $(CLUSTER_NAME)

## Borra el cluster + archivos generados localmente
clean: down
	@rm -f k3d-config.yaml .harbor-ip
