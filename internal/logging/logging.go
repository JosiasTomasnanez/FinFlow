// Package logging centraliza el logging estructurado (JSON) de FinFlow y
// la correlación de logs entre capas/servicios vía un request_id.
//
// Idea general: cada request HTTP que entra al backend recibe (o hereda,
// si ya venía en un header) un request_id. Ese id se guarda en el
// context.Context de Go y se propaga a todas las capas (handler -> service
// -> storage) para que cualquier log emitido durante ese request incluya
// el mismo request_id. También se devuelve en la respuesta HTTP para que
// otros servicios (por ejemplo el frontend) puedan loguear el mismo id y
// así correlacionar logs entre servicios distintos.
package logging

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const (
	requestIDKey ctxKey = "request_id"
	loggerKey    ctxKey = "logger"
)

// HeaderRequestID es el header HTTP usado para propagar el id de
// correlación entre servicios (frontend <-> backend, y eventualmente
// backend <-> lo que instrumente Jero con OpenTelemetry).
const HeaderRequestID = "X-Request-ID"

// base es el logger de proceso. Todos los loggers "por request" se derivan
// de este con .With(...), así todos comparten el mismo formato JSON y el
// mismo campo "service".
var base = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
})).With(slog.String("service", "finflow-backend"))

// Base devuelve el logger de proceso, para usar fuera de un request
// (arranque de la app, jobs de background, etc.).
func Base() *slog.Logger {
	return base
}

// WithRequestID devuelve una copia de ctx que lleva el requestID y un
// logger que ya tiene el campo request_id seteado, para que todos los
// logs posteriores lo incluyan automáticamente.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	ctx = context.WithValue(ctx, requestIDKey, requestID)
	logger := base.With(slog.String("request_id", requestID))
	ctx = context.WithValue(ctx, loggerKey, logger)
	return ctx
}

// FromContext devuelve el logger scopeado al request guardado por
// WithRequestID, o el logger base si no hay ninguno (por ejemplo en tests
// o en código que corre fuera de un request HTTP).
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return base
}

// RequestIDFromContext devuelve el request_id guardado en ctx, o "" si no
// hay ninguno.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
