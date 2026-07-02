#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# seal-secrets.sh
# Genera SealedSecrets para todos los ambientes de FinFlow a partir de
# archivos .env locales (NUNCA commiteados). El output SÍ se commitea con
# condicionales nativos de Helm inyectados dinámicamente.
# ==============================================================================

CONTROLLER_NAME="sealed-secrets"
CONTROLLER_NAMESPACE="kube-system"
OUTPUT_DIR="finflow-chart/templates/secrets"
SECRETS_DIR="secrets"

# Mantenemos tus nombres originales consistentes con el backend (-app-secret)
declare -A CONFIGS=(
  ["prod"]="finflow-prod finflow-prod-app-secret"
  ["staging"]="finflow-staging finflow-staging-app-secret"
  ["infra"]="finflow-infra finflow-infra-unleash-secret"
)

mkdir -p "$OUTPUT_DIR"

echo "=================================================="
echo "  Verificando conexión al cluster..."
echo "=================================================="
kubectl cluster-info > /dev/null || { echo "❌ No se pudo conectar al cluster"; exit 1; }

echo "Contexto actual: $(kubectl config current-context)"
read -p "¿Es el cluster correcto? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Abortado. Cambiá de contexto con: kubectl config use-context <nombre>"
  exit 1
fi

echo ""
echo "=================================================="
echo "  Generando SealedSecrets con bloques Helm {{ if }}"
echo "=================================================="

for key in "${!CONFIGS[@]}"; do
  read -r NAMESPACE SECRET_NAME <<< "${CONFIGS[$key]}"
  ENV_FILE="${SECRETS_DIR}/${key}.env"

  if [[ ! -f "$ENV_FILE" ]]; then
    echo "⚠️  $ENV_FILE no existe, salteando '$key'"
    continue
  fi

  echo ""
  echo "🔐 Sellando '$key' -> namespace=$NAMESPACE secret=$SECRET_NAME"

  # Armar los --from-literal a partir del .env, ignorando comentarios y líneas vacías
  LITERAL_ARGS=()
  while IFS='=' read -r k v; do
    [[ -z "$k" || "$k" == \#* ]] && continue
    LITERAL_ARGS+=(--from-literal="${k}=${v}")
  done < "$ENV_FILE"

  if [[ ${#LITERAL_ARGS[@]} -eq 0 ]]; then
    echo "⚠️  $ENV_FILE está vacío, salteando"
    continue
  fi

  OUTPUT_FILE="${OUTPUT_DIR}/${key}-sealed-secret.yaml"
  TEMP_FILE="${OUTPUT_FILE}.tmp"

  # 1. Generar el SealedSecret plano temporal mediante kubeseal
  kubectl create secret generic "$SECRET_NAME" \
    --namespace "$NAMESPACE" \
    "${LITERAL_ARGS[@]}" \
    --dry-run=client -o yaml | \
  kubeseal --format=yaml \
    --controller-name="$CONTROLLER_NAME" \
    --controller-namespace="$CONTROLLER_NAMESPACE" \
    > "$TEMP_FILE"

  # 2. Inyectar de forma dinámica el condicional {{- if }} de Helm al inicio del archivo
  echo "{{- if or (contains \"${key}\" .Release.Name) (eq .Values.environment \"${key}\") }}" > "$OUTPUT_FILE"
  
  # 3. Modificar dinámicamente los nombres fijos por tags de Helm e inyectar el contenido
  # Esto cambia "finflow-prod-app-secret" por "{{ .Release.Name }}-app-secret" dinámicamente
  if [[ "$key" != "infra" ]]; then
    sed -e "s/name: ${SECRET_NAME}/name: {{ .Release.Name }}-app-secret/g" "$TEMP_FILE" >> "$OUTPUT_FILE"
  else
    cat "$TEMP_FILE" >> "$OUTPUT_FILE"
  fi
  
  # 4. Inyectar el {{- end }} de Helm al final del archivo
  echo "{{- end }}" >> "$OUTPUT_FILE"

  # Limpiar el archivo temporal
  rm "$TEMP_FILE"

  echo "✅ Generado e inyectado con éxito: $OUTPUT_FILE"
done

echo ""
echo "=================================================="
echo "  Listo. Próximos pasos:"
echo "  1. Revisá los archivos generados en $OUTPUT_DIR"
echo "  2. git add $OUTPUT_DIR"
echo "  3. Commit + push"
echo "  4. ¡Dale Sync en Argo CD!"
echo "=================================================="
