#!/usr/bin/env bash

#############################################
# FinFlow Dev Console
#############################################

set -u

PID_DIR="/tmp/finflow-dev"
mkdir -p "$PID_DIR"

#############################################
# Colores
#############################################

GREEN="\033[0;32m"
RED="\033[0;31m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
BOLD="\033[1m"
NC="\033[0m"

#############################################
# Configuración
#############################################

GRAFANA_URL="http://localhost:8082"
UNLEASH_URL="http://localhost:8081"
ARGOCD_URL="https://localhost:8080"
ROLLOUTS_URL="http://localhost:3100"

# Cambiá estas URLs por las de tu ALB
STAGING_URL="http://k8s-finflowstaginggro-2d1f40137d-1302572390.us-east-1.elb.amazonaws.com / http://staging.finflow.local/"
PRODUCTION_URL="http://k8s-finflowprodgroup-be3700653a-396980022.us-east-1.elb.amazonaws.com / http://prod.finflow.local/"

#############################################

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

#############################################

service_running() {

    local pidfile="$PID_DIR/$1.pid"

    if [[ ! -f "$pidfile" ]]; then
        return 1
    fi

    local pid
    pid=$(cat "$pidfile")

    kill -0 "$pid" 2>/dev/null
}

#############################################

start_service() {

    local name="$1"
    local command="$2"

    if service_running "$name"; then
        echo "✓ $name ya estaba iniciado."
        return
    fi

    bash -c "$command" >/dev/null 2>&1 &
    echo $! > "$PID_DIR/$name.pid"

    echo "✓ $name iniciado."
}

#############################################

stop_service() {

    local name="$1"

    if service_running "$name"; then
        kill "$(cat "$PID_DIR/$name.pid")"
        rm -f "$PID_DIR/$name.pid"
        echo "✓ $name detenido."
    fi
}

#############################################

status_icon() {

    if service_running "$1"; then
        echo -e "${GREEN}🟢${NC}"
    else
        echo -e "${RED}🔴${NC}"
    fi
}

#############################################

open_url() {

    if command_exists xdg-open; then
        xdg-open "$1" >/dev/null 2>&1 &
    fi
}

#############################################

start_all() {

    echo
    echo -e "${BLUE}Iniciando servicios...${NC}"
    echo

    start_service grafana \
        "kubectl port-forward svc/finflow-infra-grafana-service 8082:3000 -n finflow-infra"

    start_service unleash \
        "kubectl port-forward svc/unleash-web 8081:4242 -n finflow-infra"

    start_service argocd \
        "kubectl port-forward svc/argocd-server 8080:443 -n argocd"

    start_service rollouts \
        "kubectl argo rollouts dashboard -n argo-rollouts"

    sleep 2
    echo
    read -p "Presione Enter para continuar..."
}

#############################################

stop_all() {

    echo
    echo "Deteniendo servicios..."
    echo

    stop_service grafana
    stop_service unleash
    stop_service argocd
    stop_service rollouts

    echo
    read -p "Presione Enter para continuar..."
}

#############################################

cluster_info() {

    CONTEXT=$(kubectl config current-context 2>/dev/null)

    if [[ -z "$CONTEXT" ]]; then
        CONTEXT="No conectado"
    fi

    echo "$CONTEXT"
}

#############################################

draw() {

clear

echo -e "${BOLD}==============================================================${NC}"
echo -e "${CYAN}                     FinFlow Dev Console${NC}"
echo -e "${BOLD}==============================================================${NC}"

echo
echo "Cluster:"
echo "  $(cluster_info)"
echo

printf "%-3s %-24s %s\n" "$(status_icon grafana)" "Grafana" "$GRAFANA_URL"
printf "%-3s %-24s %s\n" "$(status_icon unleash)" "Unleash" "$UNLEASH_URL"
printf "%-3s %-24s %s\n" "$(status_icon argocd)" "Argo CD" "$ARGOCD_URL"
printf "%-3s %-24s %s\n" "$(status_icon rollouts)" "Rollouts Dashboard" "$ROLLOUTS_URL"

echo
echo "Staging : $STAGING_URL"
echo "Production : $PRODUCTION_URL"

echo
echo "--------------------------------------------------------------"

echo "1) Start Local Services"
echo "2) Stop Local Services"
echo "3) Open Grafana"
echo "4) Open Unleash"
echo "5) Open Argo CD"
echo "6) Open Rollouts"
echo "7) Open Staging"
echo "8) Open Production"
echo "9) k9s"
echo "10) Refresh"
echo "11) Exit"
echo
}

#############################################

while true
do

draw

read -rp "Seleccione una opción: " option

case $option in

1)
    start_all
    ;;

2)
    stop_all
    ;;

3)
    open_url "$GRAFANA_URL"
    ;;

4)
    open_url "$UNLEASH_URL"
    ;;

5)
    open_url "$ARGOCD_URL"
    ;;

6)
    open_url "$ROLLOUTS_URL"
    ;;

7)
    open_url "$STAGING_URL"
    ;;

8)
    open_url "$PRODUCTION_URL"
    ;;

9)
    clear
    if command_exists k9s; then
        k9s
    else
        echo "k9s no está instalado."
        read -p "Presione Enter para continuar..."
    fi
    ;;

10)
    ;;

11)
    clear
    exit 0
    ;;

*)
    echo "Opción inválida."
    sleep 1
    ;;

esac

done
