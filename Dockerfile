# --- Etapa 1: Construcción del Backend en Go ---
# Usamos $BUILDPLATFORM para que Go se ejecute a velocidad nativa x86 en tu servidor
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Buildx inyecta automáticamente el OS y la arquitectura de destino (amd64, arm64, etc.)
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Copiamos el código fuente del backend
COPY . .

# Configuraciones de Go y descarga de dependencias
RUN go env -w GOTOOLCHAIN=auto
RUN go mod tidy
RUN go mod download

# Compilamos definiendo GOARCH para que Go genere el binario correcto sin usar QEMU
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o finflow ./cmd/finflow

# --- Etapa 2: Imagen Final de Producción ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Solo nos traemos el binario ejecutable
COPY --from=builder /app/finflow ./finflow

EXPOSE 8080

ENTRYPOINT ["./finflow"]
