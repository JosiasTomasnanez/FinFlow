FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o finflow ./cmd/finflow

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/finflow ./finflow
EXPOSE 8080
ENTRYPOINT ["./finflow"]
