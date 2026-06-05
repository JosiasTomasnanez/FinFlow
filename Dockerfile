# --- Etapa 1: Construcción del Backend en Go ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copiamos el código fuente del backend
COPY . .

# Configuraciones de Go y descarga de dependencias
RUN go env -w GOTOOLCHAIN=auto
RUN go mod tidy
RUN go mod download

# Compilamos el binario de forma optimizada
RUN CGO_ENABLED=0 GOOS=linux go build -o finflow ./cmd/finflow

# --- Etapa 2: Imagen Final de Producción (Ultra ligera)
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Solo nos traemos el binario ejecutable
COPY --from=builder /app/finflow ./finflow

# EXPOSE es informativo, Render usará el puerto que definas (ej: 8080)
EXPOSE 8080

ENTRYPOINT ["./finflow"]
