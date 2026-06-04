# --- Etapa 1: Construcción del Frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --ignore-scripts
COPY frontend ./
RUN npm run build

# --- Etapa 2: Construcción del Backend en Go ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# 1. Copiamos ABSOLUTAMENTE TODO el código fuente primero
# Esto incluye tu go.mod, go.sum y la carpeta internal/storage
COPY . .

# 2. Como Docker ya tiene tu código Y tiene Go 1.25 nativo,
# acá sí el tidy va a escanear todo y va a corregir el go.sum interno con éxito.
RUN go env -w GOTOOLCHAIN=auto
RUN go mod tidy
RUN go mod download

# Traemos el build del frontend antes de compilar Go
COPY --from=frontend-builder /frontend/dist ./frontend/dist

# Compilamos de forma optimizada para entornos Alpine (estático)
RUN CGO_ENABLED=0 GOOS=linux go build -o finflow ./cmd/finflow

# --- Etapa 3: Imagen Final de Producción ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/finflow ./finflow
COPY --from=builder /app/frontend/dist ./frontend/dist

EXPOSE 8080
ENTRYPOINT ["./finflow"]
