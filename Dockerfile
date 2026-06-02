# --- Etapa 1: Construcción del Frontend ---
FROM node:20-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package.json frontend/package-lock.json ./
# Usamos npm ci para velocidad y consistencia, e ignore-scripts por seguridad (SonarQube)
RUN npm ci --ignore-scripts
COPY frontend ./
RUN npm run build

# --- Etapa 2: Construcción del Backend en Go ---
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copiamos todo el código fuente
COPY . .
# Traemos el build del frontend antes de compilar Go (crucial si usas go:embed)
COPY --from=frontend-builder /frontend/dist ./frontend/dist
# Compilamos de forma optimizada para entornos Alpine (estático)
RUN CGO_ENABLED=0 GOOS=linux go build -o finflow ./cmd/finflow

# --- Etapa 3: Imagen Final de Producción ---
FROM alpine:latest
# Por seguridad, es buena práctica instalar certificados CA por si tu app de Go hace peticiones HTTPS
RUN apk --no-cache add ca-certificates
WORKDIR /app
# Copiamos solo el binario ejecutable
COPY --from=builder /app/finflow ./finflow

# NOTA: Descomenta la siguiente línea SOLO si tu backend de Go NO usa "go:embed" 
# y necesita leer la carpeta dist directamente desde el disco.
COPY --from=builder /app/frontend/dist ./frontend/dist

EXPOSE 8080
ENTRYPOINT ["./finflow"]